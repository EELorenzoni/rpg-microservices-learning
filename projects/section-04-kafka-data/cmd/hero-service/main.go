package main

import (
	"fmt"
	"log"

	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/core/services/herosrv"
	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/repositories/herorepo"
)

func main() {
	fmt.Println("🛡️  Hero Service Starting...")

	// 1. INFRASTRUCTURE: Crear adaptadores concretos (Repositories)
	// broker: localhost:9094 (definido en docker-compose)
	// topic: hero-created-04
	repo := herorepo.NewKafka("localhost:9094", "hero-created-04")
	// Importante: Cerrar conexión al terminar
	defer repo.Close()

	// 2. CORE: Inyectar dependencias (Services)
	// createHeroService es ahora un *herosrv.Service
	createHeroService := herosrv.New(repo)

	// 3. HANDLER/EXECUTION: Simular una petición
	// command ahora pertenece a herosrv
	cmd := herosrv.CreateHeroCommand{
		ID:    "h-1",
		Name:  "Aragorn",
		Power: 90,
	}

	// Ejecutar el caso de uso
	err := createHeroService.Run(cmd)
	if err != nil {
		log.Fatalf("❌ Error ejecutando caso de uso: %v", err)
	}

	fmt.Println("🎉 DEMO FINALIZADA: Capas Integradas Correctamente (Refactorizado).")
}
