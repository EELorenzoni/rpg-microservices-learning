.PHONY: init up down logs

# Inicializa el módulo de Go (idempotente)
init:
	@if [ ! -f go.mod ]; then \
		go mod init github.com/EELorenzoni/rpg-microservices-learning; \
	fi
	go mod tidy

# Levanta los servicios de Docker (Kafka + Zookeeper)
up:
	docker compose up -d

# Baja los servicios
down:
	docker compose down

# Muestra los logs de Kafka
logs:
	docker compose logs -f kafka
