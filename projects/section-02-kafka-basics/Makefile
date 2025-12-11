.PHONY: init up down logs

# Default: inicializa y levanta todo
all: init up run-produce run-consume

# Inicializa el módulo de Go (idempotente)
init:
	@if [ ! -f go.mod ]; then \
		go mod init github.com/EELorenzoni/rpg-microservices-learning; \
	fi
	go mod tidy

# Ejecuta el productor
run-produce:
	go run cmd/tutorial-basics/main.go produce

# Ejecuta el consumidor
run-consume:
	go run cmd/tutorial-basics/main.go consume

# Levanta los servicios de Docker (Kafka + Zookeeper)
up:
	docker compose up -d

# Baja los servicios
down:
	docker compose down

# Muestra los logs de Kafka
logs:
	docker compose logs -f kafka
