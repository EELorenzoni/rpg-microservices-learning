package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// Constantes de configuración (hardcodeadas para simplicidad del tutorial)
const (
	topic         = "rpg-battles"
	brokerAddress = "localhost:9094"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go [produce|consume]")
		os.Exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "produce":
		produceMessages()
	case "consume":
		consumeMessages()
	default:
		fmt.Printf("Modo desconocido: %s. Usa 'produce' o 'consume'\n", mode)
		os.Exit(1)
	}
}

func produceMessages() {
	// 0. Asegurar que el topic existe
	ensureTopic()

	// 1. Crear el escritor (writer)
	// Kafka-go maneja el pooling de conexiones automáticamente.
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()

	fmt.Println("⚔️  Iniciando Productor de Batallas...")

	for i := 1; i <= 5; i++ {
		msgValue := fmt.Sprintf("Heroe ataca a Orco #%d con 50 de daño", i)

		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(fmt.Sprintf("Key-%d", i)),
				Value: []byte(msgValue),
			},
		)
		if err != nil {
			log.Fatal("Error enviando mensaje:", err)
		}

		fmt.Printf("Enviado: %s\n", msgValue)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("✅ Todos los ataques enviados.")
}

func consumeMessages() {
	ensureTopic()

	// 1. Crear el lector (reader)
	// Configurado con GroupID para simular un microservicio real
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		GroupID:     "battle-stats-service", // Identificador del grupo de consumidores
		MinBytes:    10e3,                   // 10KB
		MaxBytes:    10e6,                   // 10MB
		StartOffset: kafka.FirstOffset,      // Para tutorial: leer desde el inicio siempre
	})
	defer r.Close()

	fmt.Println("🛡️  Iniciando Consumidor de Batallas (Esperando eventos)...")

	// Loop infinito de lectura
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("⚠️  Error leyendo mensaje: %v\n    --> Reintentando en 1s...\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("Mensaje recibido: %s (offset %d)\n", string(m.Value), m.Offset)
	}
}

func ensureTopic() {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatal("Error conectando a Kafka para verificar topic:", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("Error obteniendo controlador:", err)
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		log.Fatal("Error conectando al controlador:", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		// Ignorar error si ya existe
		// Es una forma simplificada, en prod verificaríamos el tipo de error
		fmt.Printf("Nota: Topic '%s' ya existe o acaba de ser creado.\n", topic)
	}
}
