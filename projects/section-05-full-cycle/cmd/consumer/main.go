package main

import (
	"context"
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/handlers/herohdl"
	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("🐢 Hero Consumer Starting...")

	// 1. INFRASTRUCTURE (Kafka Reader)
	// Configuración para leer de "hero-events-05"
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9094"},
		Topic:    "hero-events-05", // Mismo topic que el API publica
		GroupID:  "hero-group-1",   // Importante para Consumer Groups
		MinBytes: 10e3,             // 10KB
		MaxBytes: 10e6,             // 10MB
	})
	defer reader.Close()

	// 2. HANDLER
	consumer := herohdl.NewConsumerHandler(reader)

	// 3. EXECUTION
	// Bloquea por siempre (o hasta Ctrl+C)
	consumer.Start(context.Background())
}
