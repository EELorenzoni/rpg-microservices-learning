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
// 🎓 PATRÓN: Consumer Group + DLQ (Robustness)
func (h *ConsumerHandler) Start(ctx context.Context) {
	fmt.Println("🎧 HANDLER (Consumer): Esperando eventos en Kafka...")

	// 1. DLQ Writer (Producer para errores)
	// En prod, esto debería inyectarse como dependencia. Lo creamos aquí por simplicidad.
	dlqWriter := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9094"),
		Topic:    "hero-events-05-dlq",
		Balancer: &kafka.LeastBytes{},
	}
	defer dlqWriter.Close()

	for {
		// 2. Leer Mensaje (Bloqueante)
		m, err := h.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("❌ CRITICAL: Error de conexión con Kafka: %v\n", err)
			break // Romper el loop si Kafka se cae
		}

		// 3. Procesar (Simulación de fallo aleatorio o validación)
		err = h.processMessage(m)

		if err != nil {
			// 💀 DEAD LETTER QUEUE (DLQ)
			// Si fallamos, no reintentamos infinitamente. Lo movemos al DLQ.
			log.Printf("⚠️ ERROR PROCESANDO (Offset %d): %v. Enviando a DLQ...\n", m.Offset, err)

			errDLQ := dlqWriter.WriteMessages(ctx, kafka.Message{
				Key:   m.Key,   // Mantenemos la Key original
				Value: m.Value, // Mantenemos el Payload original
				Headers: []kafka.Header{
					{Key: "original-topic", Value: []byte(m.Topic)},
					{Key: "error-reason", Value: []byte(err.Error())},
				},
			})

			if errDLQ != nil {
				log.Printf("🔥 FATAL: No se pudo escribir en DLQ: %v\n", errDLQ)
				// Aquí sí podríamos reintentar o pausar.
			} else {
				log.Printf("🗑️ Enviado a DLQ: hero-events-05-dlq\n")
			}
		}

		// 4. COMMIT (Siempre avanzamos, ya sea éxito o DLQ)
		// Si no hiciéramos commit tras el DLQ, leeríamos el mensaje venenoso infinitamente.
		if err := h.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("❌ Error haciendo commit: %v\n", err)
		}
	}
}

// processMessage simula la lógica de negocio.
func (h *ConsumerHandler) processMessage(m kafka.Message) error {
	// Log de Auditoría Completo
	fmt.Printf("\n📨 CONSUMER (P:%d | O:%d) Key[%s]\n", m.Partition, m.Offset, string(m.Key))
	fmt.Printf("   Payload: %s\n", string(m.Value))

	// Simulación de "Poison Message": Si el payload contiene "error", fallamos.
	payload := string(m.Value)
	if payload == `{"fail":true}` {
		return fmt.Errorf("simulated business error")
	}

	// Éxito
	fmt.Println("   ✅ Procesado correctamente.")
	return nil
}
