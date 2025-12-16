.PHONY: init up down logs run

# Default: inicializa y ejecuta la demo
all: init run-api

# Inicializa dependencias
init:
	go mod tidy

run-api:
	go run cmd/api/main.go

run-consumer:
	go run cmd/consumer/main.go

up:
	docker-compose -f ../../docker-compose.yml up -d

down:
	docker-compose -f ../../docker-compose.yml down

