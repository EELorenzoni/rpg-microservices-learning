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

// Save serializa el héroe a JSON y lo envía a Kafka.
//
// 💡 POINTERS: (repo *Kafka) vs (hero *domain.Hero)
//  1. (repo *Kafka): NECESARIO. El 'writer' de Kafka mantiene un pool de conexiones TCP interno.
//     Si copiáramos el repo (por valor), podríamos duplicar/perder el estado de la conexión.
//     Queremos que TODOS usen LA MISMA conexión abierta.
//  2. (hero *domain.Hero): EFICIENCIA. No queremos copiar todos los datos del héroe, solo leerlos.
//
// 💡 POINTERS (Sintaxis):
// - `(repo *Kafka)`: Receiver de tipo Puntero.
// - `(hero *domain.Hero)`: Argumento de tipo Puntero.
// - Dentro de la función, usamos `repo.writer` directamente. Go hace "dereference" automático (*repo).writer.
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
