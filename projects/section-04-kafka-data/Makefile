.PHONY: init up down logs run

# Default: inicializa y ejecuta la demo
all: init run

# Inicializa dependencias
init:
	go mod tidy

# Levanta contenedores (aunque no los usamos todavía en el código, está bien tenerlos)
up:
	docker compose up -d

# Baja contenedores
down:
	docker compose down

# Ejecuta el microservicio (Producción a Kafka)
run:
	go run cmd/hero-service/main.go
