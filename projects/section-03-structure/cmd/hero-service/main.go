package main

import (
	"fmt"
	"log"

	"github.com/EELorenzoni/rpg-microservices-learning/section-03/internal/application"
	"github.com/EELorenzoni/rpg-microservices-learning/section-03/internal/infrastructure/adapters"
)

func main() {
	fmt.Println("🛡️  Hero Service Starting...")

	// 1. INFRASTRUCTURE: Crear adaptadores concretos
	// En el futuro, aquí inicializaremos KafkaProducer o PostgresDB
	repo := adapters.NewInMemoryHeroRepository()

	// 2. APPLICATION: Inyectar dependencias (Wiring / Composición)
	// Le pasamos el 'repo' concreto, pero el servicio solo ve la interfaz 'HeroRepository'
	createHeroService := application.NewCreateHeroService(repo)

	// 3. EXECUTION: Simular una petición (ej: HTTP request o Mensaje Kafka)
	// Creamos un comando (DTO) como si viniera de un JSON body
	cmd := application.CreateHeroCommand{
		ID:    "h-1",
		Name:  "Aragorn",
		Power: 90,
	}

	// Ejecutar el caso de uso
	err := createHeroService.Run(cmd)
	if err != nil {
		log.Fatalf("❌ Error ejecutando caso de uso: %v", err)
	}

	fmt.Println("🎉 DEMO FINALIZADA: Capas Integradas Correctamente.")
}
