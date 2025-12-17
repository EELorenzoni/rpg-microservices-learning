package main

import (
	"fmt"
	"log"
	"os"

	"github.com/EELorenzoni/rpg-microservices-learning/platform-kafka-admin/internal/core"
	"github.com/EELorenzoni/rpg-microservices-learning/platform-kafka-admin/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 0. Load .env (if exists)
	if err := godotenv.Load(); err != nil {
		fmt.Println("ℹ️  No .env file found (relying on system env)")
	}

	port := os.Getenv("ADMIN_PORT")
	if port == "" {
		log.Fatal("❌ FATAL: ADMIN_PORT is not set in .env or environment")
	}

	fmt.Printf("🚀 Kafka Admin API starting on %s\n", port)

	// 1. Core
	brokerAddress := os.Getenv("KAFKA_BROKER")
	if brokerAddress == "" {
		log.Fatal("❌ FATAL: KAFKA_BROKER is not set in .env or environment")
	}
	fmt.Printf("🔧 Config: Broker=%s\n", brokerAddress)

	service := core.NewAdminService(brokerAddress)

	// 2. Handlers
	handler := handlers.NewAdminHandler(service)

	// 3. Router
	r := gin.Default()

	r.POST("/topics", handler.CreateTopic)
	r.GET("/topics", handler.ListTopics)
	r.DELETE("/topics/:name", handler.DeleteTopic)

	r.Run(port)
}
