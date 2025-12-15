package main

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/core/services/herosrv"
	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/handlers/herohdl"
	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/repositories/herorepo"
)

func main() {
	fmt.Println("🛡️  Hero Service Starting...")

	// 1. INFRASTRUCTURE (Repositories/Ports)
	repo := herorepo.NewKafka("localhost:9094", "hero-created-04")
	defer repo.Close()

	// 2. CORE (Services/UseCases)
	createHeroService := herosrv.New(repo)

	// 3. HANDLER (Driving Adapter)
	// Aquí conectamos la "Entrada" (CLI) con el "Core" (Service)
	cliHandler := herohdl.NewCLIHandler(createHeroService)

	// 4. EXECUTION
	// Simulamos que el usuario ejecuta el comando
	cliHandler.CreateHeroSimulated("h-1", "Aragorn")

	fmt.Println("🎉 DEMO FINALIZADA.")
}
