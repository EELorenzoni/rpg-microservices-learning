package herohdl

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

// ConsumerHandler es un "Driving Adapter" asíncrono.
// Reacciona a eventos en vez de peticiones HTTP.
type ConsumerHandler struct {
	reader *kafka.Reader
}

// NewConsumerHandler crea el consumidor.
func NewConsumerHandler(reader *kafka.Reader) *ConsumerHandler {
	return &ConsumerHandler{
		reader: reader,
	}
}

// Start inicia el loop de consumo.
func (h *ConsumerHandler) Start(ctx context.Context) {
	fmt.Println("🎧 HANDLER (Consumer): Esperando eventos en Kafka...")

	for {
		// 1. Leer Mensaje (Bloqueante)
		m, err := h.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Error leyendo mensaje: %v\n", err)
			break
		}

		// 2. Procesar (Loggear)
		fmt.Printf("\n📨 CONSUMER: Recibido Evento!\n")
		fmt.Printf("   Topic: %s\n", m.Topic)
		fmt.Printf("   Key:   %s\n", string(m.Key))
		fmt.Printf("   Value: %s\n", string(m.Value))

		// Aquí podríamos llamar a un Service si tuviéramos lógica de negocio.
	}
}
