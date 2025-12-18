.PHONY: all init run-api run-consumer

# Default: inicializa y ejecuta la demo
all: init run-api

# Inicializa dependencias
init:
	go mod tidy

run-api:
	go run cmd/api/main.go

run-consumer:
	go run cmd/consumer/main.go


