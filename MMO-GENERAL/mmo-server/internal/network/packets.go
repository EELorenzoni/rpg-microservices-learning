package network

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Tipos de paquetes (Protocolo Fase 0)
const (
	PacketTypeHandshake uint8 = 0 // El primer saludo
	PacketTypeMove      uint8 = 1 // Actualización de posición
	PacketTypeHeartbeat uint8 = 2 // Latido de conexión (proximamente)
)

// Tamaño del encabezado:
// 1 byte (Tipo) + 4 bytes (Secuencia) + 8 bytes (PlayerID) = 13 bytes en total
const HeaderSize = 13

// PacketHeader es lo mínimo que tienen todos nuestros paquetes
type PacketHeader struct {
	Type     uint8  // Qué tipo de mensaje es
	Sequence uint32 // El número de mensaje (para ordenarlos si llegan desordenados)
	PlayerID uint64 // A qué jugador pertenece
}

// DeserializeHeader toma los bytes crudos y los convierte en una estructura entendible
func DeserializeHeader(data []byte) (PacketHeader, error) {
	if len(data) < HeaderSize {
		return PacketHeader{}, fmt.Errorf("paquete demasiado corto")
	}

	// Usamos BigEndian porque es el "idioma" estándar de las redes
	return PacketHeader{
		Type:     data[0],
		Sequence: binary.BigEndian.Uint32(data[1:5]),
		PlayerID: binary.BigEndian.Uint64(data[5:13]),
	}, nil
}

// SerializeHandshakeResponse construye el paquete de respuesta al handshake
func SerializeHandshakeResponse(playerID uint64) []byte {
	buf := make([]byte, HeaderSize)
	buf[0] = PacketTypeHandshake
	binary.BigEndian.PutUint32(buf[1:5], 0) // Secuencia 0 para respuestas simples
	binary.BigEndian.PutUint64(buf[5:13], playerID)
	return buf
}

// MoveData contiene las coordenadas de Unreal Engine
type MoveData struct {
	X, Y, Z, Yaw float32
}

// DeserializeMove extrae las coordenadas de un paquete de tipo Move
func DeserializeMove(data []byte) (MoveData, error) {
	// Un paquete de movimiento tiene: Header (13) + Payload (16) = 29 bytes
	if len(data) < HeaderSize+16 {
		return MoveData{}, fmt.Errorf("paquete de movimiento incompleto")
	}

	// El payload empieza después de los 13 bytes de la cabecera
	payload := data[HeaderSize : HeaderSize+16]

	return MoveData{
		X:   Float32frombytes(payload[0:4]),
		Y:   Float32frombytes(payload[4:8]),
		Z:   Float32frombytes(payload[8:12]),
		Yaw: Float32frombytes(payload[12:16]),
	}, nil
}

// Float32frombytes convierte 4 bytes en un número decimal de 32 bits (Estándar IEEE 754)
func Float32frombytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

// Float32bytes convierte un número decimal en 4 bytes para enviarlo por red
func Float32bytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)
	return bytes
}
