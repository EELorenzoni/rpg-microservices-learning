package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"mmo-server/internal/network"
)

// Configuración del servidor
const (
	TickRate = 30                     // Cuántas veces por segundo se actualiza el mundo (30Hz)
	TickTime = time.Second / TickRate // La duración exacta de cada tick (aprox 33.3ms)
)

// RawPacket representa un paquete tal cual llega del socket, antes de ser procesado
type RawPacket struct {
	Addr *net.UDPAddr // Quién lo envió (IP y Puerto)
	Data []byte       // El contenido binario (los bytes)
}

func main() {
	// 1. Inicializamos el Connection Manager (El que sabe quién está conectado)
	connMgr := network.NewConnectionManager()

	// 2. Definimos dónde va a escuchar el servidor (Cualquier IP, puerto 8080)
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.ParseIP("0.0.0.0"),
	}

	// 3. Abrimos el socket UDP
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("❌ Error inicializando el socket: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close() // Se asegura de cerrar el puerto al terminar el programa

	fmt.Printf("🚀 MMO Game Server iniciado - Fase 0 Completa\n")
	fmt.Printf("📡 Escuchando en UDP %s\n", conn.LocalAddr().String())
	fmt.Printf("💓 Tick Rate: %d Hz\n", TickRate)

	// Canal de Go: Es como una tubería para pasar datos entre diferentes partes del programa
	// Aquí lo usamos para pasar paquetes desde el socket al bucle principal
	packetChan := make(chan RawPacket, 100)

	// Goroutine: Es un "hilo" ligero. Aquí lanzamos un proceso en paralelo que solo lee del socket
	go func() {
		buffer := make([]byte, 1024) // Buffer temporal para cada lectura
		for {
			n, remoteAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				continue // Si hay error de lectura, seguimos esperando el siguiente
			}
			// Copiamos los bytes recibidos a una nueva rebanada (slice) para no sobrescribirlos
			data := make([]byte, n)
			copy(data, buffer[:n])

			// Metemos el paquete en la "tubería" (canal)
			packetChan <- RawPacket{Addr: remoteAddr, Data: data}
		}
	}()

	// Ticker: Un reloj que dispara un evento cada 33.3ms (nuestro Tick Rate)
	ticker := time.NewTicker(TickTime)
	defer ticker.Stop()

	tickCount := 0
	var nextPlayerID uint64 = 1001 // Empezamos a asignar IDs desde el 1001

	// BUCLE PRINCIPAL (El corazón del servidor)
	for range ticker.C {
		tickCount++

		// En cada tick, procesamos todos los paquetes que hayan llegado desde el último tick
		processPackets(packetChan, connMgr, &nextPlayerID, conn)

		// Imprimimos estadísticas cada 5 segundos (150 ticks)
		if tickCount%150 == 0 {
			fmt.Printf("📊 Tick #%d - Jugadores activos: %d\n", tickCount, connMgr.TotalPlayers())
		}
	}
}

// processPackets vacía el canal de paquetes y los procesa según su tipo
func processPackets(ch chan RawPacket, cm *network.ConnectionManager, nextID *uint64, conn *net.UDPConn) {
	for {
		select {
		case rp := <-ch:
			// 1. Deserializamos la cabecera (los primeros 13 bytes comunes a todo paquete)
			header, err := network.DeserializeHeader(rp.Data)
			if err != nil {
				continue // Si el paquete es basura, lo ignoramos
			}

			// 2. Dependiendo del tipo de paquete, hacemos una acción u otra
			switch header.Type {
			case network.PacketTypeHandshake:
				handleHandshake(rp.Addr, cm, nextID, conn)
			case network.PacketTypeMove:
				handleMove(header.PlayerID, rp.Data, cm, conn)
			}
		default:
			// Si no hay más paquetes en el canal, salimos del bucle de procesamiento
			return
		}
	}
}

// handleHandshake se encarga de registrar a un nuevo jugador y darle su ID
func handleHandshake(addr *net.UDPAddr, cm *network.ConnectionManager, nextID *uint64, conn *net.UDPConn) {
	player, exists := cm.GetPlayer(addr)
	id := *nextID

	if !exists {
		// Si es la primera vez que lo vemos, lo registramos
		cm.RegisterPlayer(addr, id)
		*nextID++
	} else {
		// Si ya lo conocíamos, usamos el ID que ya tenía
		id = player.ID
	}

	// Creamos la respuesta binaria de "Bienvenida" con el ID asignado
	response := network.SerializeHandshakeResponse(id)

	// Enviamos la respuesta de vuelta por el socket UDP
	conn.WriteToUDP(response, addr)
	fmt.Printf("📡 Handshake: ID %d asignado a %v\n", id, addr)
}

// handleMove se encarga de recibir una posición y avisar al resto de jugadores (Replicación)
func handleMove(playerID uint64, data []byte, cm *network.ConnectionManager, conn *net.UDPConn) {
	// Simplemente enviamos los mismos bytes que recibimos a todos los demás jugadores
	// Esto es el "Broadcast" o Replicación
	cm.Broadcast(data, playerID, conn)
}
