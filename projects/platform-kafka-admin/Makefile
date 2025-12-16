.PHONY: dockerUp dockerDown serverUp initAll

dockerUp:
	docker-compose up -d

dockerDown:
	docker-compose down

serverUp:
	go run cmd/admin-api/main.go

# 'initAll' ejecuta dockerUp y luego serverUp.
# Nota: serverUp bloqueará la terminal, así que 'initAll' se quedará ahí hasta que mates el proceso.
initAll: dockerUp serverUp