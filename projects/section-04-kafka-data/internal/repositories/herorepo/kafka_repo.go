package herorepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/core/domain"
	"github.com/segmentio/kafka-go"
)

// Kafka define la implementación REAL que habla con Kafka.
// Renombrado de KafkaHeroRepository para ser más idiomático: herorepo.Kafka
type Kafka struct {
	writer *kafka.Writer
}

// NewKafka inicializa la conexión.
func NewKafka(brokerAddress string, topic string) *Kafka {
	// 1. Intentar crear el topic explícitamente (Mejor práctica que auto-create)
	// Conectamos "crudo" al broker líder (o cualquiera)
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		fmt.Printf("⚠️ WARN: No se pudo conectar para crear topic: %v\n", err)
	} else {
		defer conn.Close()

		topics := []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		}

		err = conn.CreateTopics(topics...)
		if err != nil {
			// Si ya existe, dará error, pero no importa.
			// fmt.Printf("ℹ️ Info: Topic creation result: %v\n", err)
		} else {
			fmt.Printf("✨ INFRA (Kafka): Topic '%s' creado exitosamente!\n", topic)
		}
	}

	// 2. Configurar el Writer (Productor)
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokerAddress),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true, // Por si acaso
	}

	fmt.Printf("🔌 INFRA (Kafka): Conectado a %s -> Topic: %s\n", brokerAddress, topic)

	return &Kafka{
		writer: writer,
	}
}

// Save serializa el héroe a JSON y lo envía a Kafka.
func (repo *Kafka) Save(hero *domain.Hero) error {
	// 1. Serializar a JSON
	heroJSON, err := json.Marshal(hero)
	if err != nil {
		return fmt.Errorf("error serializando hero: %w", err)
	}

	// 2. Crear mensaje de Kafka
	// Usamos el ID del héroe como Key para garantizar orden.
	msg := kafka.Message{
		Key:   []byte(hero.ID),
		Value: heroJSON,
		Time:  time.Now(),
	}

	// 3. Enviar (con Contexto para timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Limpieza (como finally)

	err = repo.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("error escribiendo a kafka: %w", err)
	}

	fmt.Printf("🚀 INFRA (Kafka): Enviado mensaje! Key=%s Value=%s\n", hero.ID, string(heroJSON))
	return nil
}

// Close cierra la conexión (se debe llamar al apagar el servicio).
func (repo *Kafka) Close() error {
	return repo.writer.Close()
}
