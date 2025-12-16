package herorepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"
	"github.com/segmentio/kafka-go"
)

// Kafka define la implementación REAL que habla con Kafka.
// Renombrado de KafkaHeroRepository para ser más idiomático: herorepo.Kafka
type Kafka struct {
	writer *kafka.Writer
}

// NewKafka inicializa la conexión.
// Retorna: *Kafka (Dirección de memoria del struct creado).
func NewKafka(brokerAddress string, topic string) *Kafka {
	// 1. Configurar el Writer (Productor)
	// Ya no creamos el topic aquí. Asumimos que la "Plataforma" lo creó.
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokerAddress),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: false, // Forzamos a que exista
	}

	fmt.Printf("🔌 INFRA (Kafka): Conectado a %s -> Topic: %s\n", brokerAddress, topic)

	// 💡 POINTERS (Sintaxis):
	// Usamos '&' (address of) para devolver la dirección del struct literal.
	return &Kafka{
		writer: writer,
	}
}

// Publish implementa la interfaz ports.EventBus.
// 💡 SOLID (ISP): Ahora Kafka solo se usa para lo que es bueno: eventos.
func (repo *Kafka) Publish(hero *domain.Hero, eventType string) error {
	// 1. Enriquecer el evento (CloudEvents style - simplificado)
	event := struct {
		Type string       `json:"type"`
		Data *domain.Hero `json:"data"`
	}{
		Type: eventType,
		Data: hero,
	}

	// 2. Serializar
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error serializando evento: %w", err)
	}

	// 3. Enviar a Kafka
	msg := kafka.Message{
		Key:   []byte(hero.ID),
		Value: eventJSON,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = repo.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("error publicando en kafka: %w", err)
	}

	fmt.Printf("🚀 INFRA (Kafka): Evento '%s' publicado! Key=%s\n", eventType, hero.ID)
	return nil
}

// Close cierra la conexión (se debe llamar al apagar el servicio).
func (repo *Kafka) Close() error {
	return repo.writer.Close()
}
