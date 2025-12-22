package network

import (
	"fmt"
	"net"
	"sync"
)

// Player representa a un jugador conectado en la memoria del servidor
type Player struct {
	ID   uint64       // ID único (ej. 1001)
	Addr *net.UDPAddr // IP y Puerto (para saber a dónde mandarle paquetes)
}

// ConnectionManager es el "Libro de Registro" del servidor
type ConnectionManager struct {
	players map[string]*Player // El mapa donde guardamos a los jugadores (Clave: IP:Puerto)
	mu      sync.RWMutex       // Un "Candado" (Mutex) para que dos hilos no toquen el mapa al mismo tiempo
}

// NewConnectionManager crea una nueva instancia del gestor
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		players: make(map[string]*Player),
	}
}

// RegisterPlayer agrega un nuevo jugador al mapa de forma segura
func (cm *ConnectionManager) RegisterPlayer(addr *net.UDPAddr, playerID uint64) {
	cm.mu.Lock()         // Cerramos el candado (Nadie más puede escribir ahora)
	defer cm.mu.Unlock() // Al terminar la función, se abre el candado automáticamente

	addrStr := addr.String()
	if _, exists := cm.players[addrStr]; !exists {
		cm.players[addrStr] = &Player{
			ID:   playerID,
			Addr: addr,
		}
		fmt.Printf("✅ Jugador %d registrado desde %s\n", playerID, addrStr)
	}
}

// GetPlayer busca a un jugador por su dirección IP:Puerto
func (cm *ConnectionManager) GetPlayer(addr *net.UDPAddr) (*Player, bool) {
	cm.mu.RLock() // Candado de lectura (Varios pueden leer al mismo tiempo, pero nadie puede escribir)
	defer cm.mu.RUnlock()
	p, ok := cm.players[addr.String()]
	return p, ok
}

// RemovePlayer elimina a un jugador (ej. si se desconecta)
func (cm *ConnectionManager) RemovePlayer(addr *net.UDPAddr) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.players, addr.String())
}

// TotalPlayers nos dice cuánta gente hay conectada ahora mismo
func (cm *ConnectionManager) TotalPlayers() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.players)
}

// Broadcast envía un mismo mensaje a TODOS los jugadores conectados excepto a uno (normalmente el emisor)
func (cm *ConnectionManager) Broadcast(data []byte, exceptID uint64, conn *net.UDPConn) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, p := range cm.players {
		if p.ID != exceptID {
			// Enviamos los bytes directamente al socket UDP de cada jugador
			conn.WriteToUDP(data, p.Addr)
		}
	}
}
