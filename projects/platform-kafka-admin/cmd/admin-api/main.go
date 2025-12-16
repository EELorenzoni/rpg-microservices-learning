package main

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/platform-kafka-admin/internal/core"
	"github.com/EELorenzoni/rpg-microservices-learning/platform-kafka-admin/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Kafka Admin API starting on :3000")

	// 1. Core
	brokerAddress := "127.0.0.1:9094" // Use IPv4 specifically
	service := core.NewAdminService(brokerAddress)

	// 2. Handlers
	handler := handlers.NewAdminHandler(service)

	// 3. Router
	r := gin.Default()

	r.POST("/topics", handler.CreateTopic)
	r.GET("/topics", handler.ListTopics)
	r.DELETE("/topics/:name", handler.DeleteTopic)

	r.Run(":3000")
}
