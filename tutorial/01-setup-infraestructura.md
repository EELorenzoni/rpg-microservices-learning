# Unidad 1: Configuración de Infraestructura

Bienvenido al curso de Microservicios RPG con Go y Kafka. En esta primera unidad, prepararemos el "terreno de juego": nuestro entorno de desarrollo.

## ¿Qué vamos a construir?
Una arquitectura de microservicios para un juego RPG. Usaremos **Kafka** como el sistema nervioso central que comunica eventos (ej: "MonstruoAtacado", "ExperienciaGanada") entre servicios escritos en **Go**.

## Paso 1: Estructura del Proyecto en Go
En Go, existe un estándar no oficial pero muy aceptado llamado "Standard Go Project Layout". Hemos creado tres carpetas clave:

- **`cmd/`**: Aquí viven los "Main". Es el punto de entrada. Si tuviéramos un servicio de usuarios, tendríamos `cmd/users-service/main.go`.
- **`internal/`**: Código que es privado para tu proyecto. Nadie puede importar esto desde fuera. Aquí va la lógica de negocio.
- **`pkg/`**: Código que podría ser útil para otros proyectos (librerías compartidas). Úsalo con moderación.

## Paso 2: Docker y Kafka (KRaft)
Kafka tradicionalmente necesitaba otro software llamado **Zookeeper** para funcionar. Era complejo.
Modernamente, usamos **Kafka KRaft**, que elimina esa dependencia.

Hemos creado un `docker-compose.yml` que levanta Kafka en un contenedor.
**¿Por qué Docker?** Para no ensuciar tu máquina instalando servidores complejos. Con un comando (`make up`) tienes toda la infraestructura lista.

### Configuración Clave
- **Puerto 9094**: Hemos mapeado el puerto interno de Kafka (9092) al 9094 de tu máquina, para evitar conflictos con otros procesos.
- **Red `rpg-network`**: Una red privada para que nuestros futuros microservicios se hablen entre sí.

## Paso 3: Makefile
Para no memorizar comandos largos de Docker, creamos un `Makefile`. Es como un control remoto con botones simples:
- `make init`: Prepara Go.
- `make up`: Enciende Kafka.
- `make logs`: Muestra qué está pasando.

---
**✅ Hito Alcanzado**: Tienes un entorno profesional de Go + Kafka listo para recibir código.
