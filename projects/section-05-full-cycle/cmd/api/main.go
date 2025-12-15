package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/services/herosrv"
	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/handlers/herohdl"
	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/repositories/herorepo"
)

func main() {
	const port = ":8081"
	fmt.Println("🚀 Hero API (HTTP) Starting on port", port)

	// 1. INFRASTRUCTURE (Adapters)
	// a. Base de Datos (Memoria)
	dbRepo := herorepo.NewMemory()
	// b. Event Bus (Kafka)
	eventBus := herorepo.NewKafka("localhost:9094", "hero-events-05")
	defer eventBus.Close()

	// 2. CORE (Service)
	// Inyectamos AMBAS dependencias: DB y EventBus
	service := herosrv.New(dbRepo, eventBus)

	// 3. HANDLER (HTTP Adapter)
	handler := herohdl.NewHTTPHandler(service)

	// 4. ROUTER & SERVER
	http.HandleFunc("/heroes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateHero(w, r)
		case http.MethodGet:
			handler.GetHero(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}
