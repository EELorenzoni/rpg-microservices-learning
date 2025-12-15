# 🧙‍♂️ RPG Microservices: From Prompt to Production
> **A Prompt-Driven Engineering Course**

Este documento es un tutorial maestro diseñado para ser consumido tanto por humanos como por Agentes de IA. Su objetivo es guiar la construcción de un sistema de RPG distribuido usando Go, Kafka y PostgreSQL.

---

## 📂 Fase 0: El "Prime" (Configuración del Agente)

Esta fase no escribe código, sino comportamiento. Antes de empezar, copia y pega el siguiente prompt en tu chat con la IA. Esto configurará su "sys-call" mental y su rol como experto.

### 🤖 El Prompt Maestro del Sistema

Copiar el siguiente bloque y enviarlo a la IA para iniciar la sesión:

```markdown
ACTÚA COMO UN "ARQUITECTO SENIOR EN GO E INGENIERO DE PROMPTS".

## TU OBJETIVO
Guiar al usuario en la construcción de un Backend RPG Distribuido usando Go, Apache Kafka y PostgreSQL. Estamos haciendo "Desarrollo Guiado por Prompts": no solo escribirás código, sino que explicarás los *prompts* necesarios para generar ese código en el futuro o en otros agentes.

## RESTRICCIONES TÉCNICAS
1.  **Lenguaje:** Go 1.22+ (Tipado estricto, features modernas).
2.  **Arquitectura:** Hexagonal / Clean Architecture.
    - `cmd/`: Puntos de entrada.
    - `internal/core/domain`: Lógica y Modelos (Go Puro).
    - `internal/core/ports`: Interfaces (Driver/Driven).
    - `internal/adapters`: Implementaciones (HTTP, Postgres, Kafka).
3.  **Comunicación:**
    - Síncrona: REST API (Librería estándar + router minimalista como `chi` o `net/http`).
    - Asíncrona: Kafka (usando `github.com/segmentio/kafka-go`).
4.  **Persistencia:** PostgreSQL. Usar driver `pgx`. Enfoque SQL-First (escribir migraciones -> generar structs).
5.  **Observabilidad:** Todos los servicios deben implementar `slog` (Structured Logging) y propagación básica de contexto (trace context).

## REGLAS DE COMPORTAMIENTO
1.  **Piensa Primero:** Antes de codear, define la estructura de archivos o el flujo lógico.
2.  **Schema First (Esquema Primero):** Define especificaciones de API (OpenAPI) y Esquemas de Eventos (JSON) *antes* de escribir los handlers.
3.  **Educativo:** Explica *por qué* elegiste un patrón (ej. "¿Por qué usar el patrón Outbox aquí?").
4.  **Iterativo:** Comienza con el MVP (Producto Mínimo Viable), luego refactoriza.
5.  **Idioma:** Todas tus explicaciones y comentarios deben ser en Español Latinoamericano.

## CONTEXTO ACTUAL
Estamos empezando desde cero. Espera instrucciones para la Fase 1.
```

---

## 📂 Fase 1: Arquitectura y Patrones de Diseño (El Plano)

En esta fase, enseñamos al Agente a **"Pensar antes de Codear"**. No escribiremos Go todavía. Definiremos los contratos.

### 1.1 Definición de Dominio (DDD Lite)

Usa este prompt para que la IA entienda las entidades del juego y sus relaciones.

**Prompt para Diseño de Dominio:**

```markdown
TAREA: Análisis de Dominio (DDD)

Actúa como un Experto en Domain-Driven Design.
Analiza el concepto "RPG Sencillo por Turnos" y define los Contextos Acotados (Bounded Contexts) y Entidades principales.

Requisitos:
1.  **Contexto Jugador:** Manejo de perfil, estadísticas (HP, Fuerza).
2.  **Contexto Combate:** Lógica de atacar, defender, calcular daño.
3.  **Contexto Inventario:** (Opcional por ahora) Items y equipamiento.

Salida Esperada:
- Lista de Entidades (con atributos clave).
- Lista de Value Objects (ej. `Health`, `Damage`).
- Diagrama Mermaid (classDiagram) mostrando relaciones.
```

### 1.2 Event Storming (Diseño de Eventos)

Kafka necesita mensajes claros. Definiremos qué pasa en el sistema asíncronamente.

**Prompt para Event Storming:**

```markdown
TAREA: Diseño de Eventos (Event Storming)

Basado en el dominio anterior, define los Eventos de Dominio que viajarán por Kafka.
Formato de evento: `NombreEntidad + VerboEnPasado` (ej. `PlayerAttacked`).

Para cada evento define:
1.  **Nombre:** (ej. `BattleStarted`).
2.  **Trigger:** ¿Qué acción lo dispara? (ej. "Usuario envía POST /attack").
3.  **Payload JSON:** Estructura de datos necesaria. Mínima información necesaria.

Salida Esperada:
- Tabla con Eventos, Triggers y Payloads.
- Diagrama Mermaid (sequenceDiagram) de un flujo de ataque exitoso:
  User -> Gateway -> (Produce Event) -> Kafka -> (Consume Event) -> Game Engine -> (Update DB).
```

### 1.3 Diseño de API (Schema-First)

Antes de programar el Gateway, definimos los endpoints.

**Prompt para OpenAPI:**

```markdown
TAREA: Diseño de API REST (OpenAPI 3.0)

Genera una especificación OpenAPI (YAML) para el "Gateway Service".
Endpoints requeridos:
1.  `POST /players`: Crear personaje.
2.  `POST /battle/attack`: Realizar un ataque (Input: `attacker_id`, `target_id`).
3.  `GET /players/{id}`: Ver estado actual.

Reglas:
- Usa tipos de datos estrictos.
- Define respuestas 200, 400 y 500.
- Incluye ejemplos en la documentación.
```

### 1.4 Modelado de Datos (SQL-First)

Finalmente, definimos cómo guardamos esto en Postgres.

**Prompt para Diagrama ER:**

```markdown
TAREA: Diseño de Base de Datos PostgreSQL

Diseña el esquema relacional para soportar el dominio.
Requisitos:
- Tablas normalizadas.
- Uso de UUIDs para `id`.
- Timestamps (`created_at`, `updated_at`).
- JSONB si es necesario para datos flexibles (ej. `stats` del jugador).

Salida Esperada:
- Script SQL DDL (`CREATE TABLE...`).
- Explicación de índices necesarios para performance.
```

---

## 📂 Fase 2: El Laboratorio (Infraestructura)

Aquí preparamos el terreno. El objetivo es que la IA nos genere un entorno local completo con un solo comando.

### 2.1 La Sinfonía de Contenedores (Docker)

Necesitamos Kafka y Postgres corriendo sin esfuerzo.

**Prompt para Docker Compose:**

```markdown
TAREA: Configuración de Infraestructura Local (Docker)

Genera un archivo `docker-compose.yml` robusto para desarrollo local.
Servicios requeridos:
1.  **PostgreSQL 16:** Con persistencia de datos (volume) y configuración básica de usuario/pass.
2.  **Kafka (Modo Kraft):** Sin Zookeeper si es posible (versión reciente), o con Zookeeper si es más estable para dev.
3.  **Kafka UI:** Una interfaz visual (ej. Provectus) para ver tópicos y mensajes.
4.  **Init Service:** Un contenedor efímero (`alpine`) que espere a que Postgres y Kafka estén listos (healthchecks).

Salida Esperada:
- Archivo `docker-compose.yml`.
- Comandos explicados para levantar y tumbar el entorno.
```

### 2.2 Automatización (Makefile)

Odiamos escribir comandos largos.

**Prompt para Makefile:**

```markdown
TAREA: Automatización con Makefile

Crea un `Makefile` para gestionar el ciclo de vida del proyecto.
Comandos necesarios:
- `up`: Levantar infraestructura (docker-compose up -d).
- `down`: Apagar infraestructura.
- `logs`: Ver logs de contenedores.
- `proto`: Compilar Protobufs (si decidimos usarlos, dejar placeholder).
- `lint`: Correr `golangci-lint`.
```

---

## 📂 Fase 3: Capa de Servicios A - El Gateway

Ahora sí, escribimos Go. Empezamos por el servicio que recibe al usuario.

### 3.1 Scaffolding Hexagonal

Estructura de carpetas limpia.

**Prompt para Estructura de Proyecto:**

```markdown
TAREA: Inicialización del Proyecto Gateway (Go)

Inicializa un módulo Go llamado `github.com/usuario/rpg-gateway`.
Crea la siguiente estructura de directorios basada en Clean Architecture:

/cmd/api          -> main.go (Entrypoint)
/internal
    /core
        /domain   -> Entidades (Player, Attack)
        /ports    -> Interfaces (PlayerService, EventPublisher)
    /adapters
        /http     -> Echo/Chi Handlers
        /kafka    -> Producer Implementation
        /repo     -> Postgres Implementation (si aplica, o solo en Engine)
/pkg              -> Utilitarios compartidos (Loggers, Errors)

Salida:
- Comandos `mkdir` o script bash para crearla.
- Archivo `go.mod` básico.
```

### 3.2 Handlers HTTP (El Contrato)

Implementamos los endpoints definidos en la Fase 1 (OpenAPI).

**Prompt para Handlers:**

```markdown
TAREA: Implementación de Handlers HTTP

Crea el adaptador HTTP usando la librería estándar o `chi`.
Implementa el endpoint `POST /attack`.

Requisitos:
1.  Recibir JSON body: `{"target_id": "...", "type": "melee"}`.
2.  Validar input (no IDs vacíos).
3.  Llamar al puerto `AttackService.PerformAttack(...)`.
4.  Retornar 202 Accepted (porque el procesamiento será asíncrono).

Nota: Solo crea el código del Handler y la Interfaz del Servicio. No la lógica de negocio real todavía.
```

### 3.3 Publicador de Eventos (Kafka Producer)

El Gateway no procesa el ataque, solo avisa que ocurrió.

**Prompt para Kafka Producer:**

```markdown
TAREA: Implementación del Kafka Producer

Implementa el puerto `EventPublisher` usando `segmentio/kafka-go`.
Función: `PublishAttack(ctx, event DomainEvent) error`.

Requisitos:
1.  Serializar el evento a JSON.
2.  Escribir en el tópico `attacks`.
3.  Manejar contexto para timeouts.
4.  Implementar un mecanismo de "Graceful Shutdown" para el Writer.
```

---

## 📂 Fase 4: Capa de Servicios B - El Motor (Engine)

El corazón del juego. Aquí procesamos lo que el Gateway envió. La magia ocurre **asíncronamente**.

### 4.1 El Consumidor (Kafka Consumer Group)

Necesitamos escuchar el tópico `attacks` continuamente.

**Prompt para Consumer:**

```markdown
TAREA: Implementación de Kafka Consumer Group

Crea un servicio `GameProcessor` que actúe como consumidor de Kafka.
Configuración:
- Group ID: `game-engine-group-1` (para escalar horizontalmente).
- Tópico: `attacks`.

Código Requerido:
1.  Un bucle infinito `for` que lea mensajes usando `reader.FetchMessage`.
2.  Manejo de señales (SIGTERM) para cerrar la conexión limpiamente.
3.  Una función `processMessage` (placeholder por ahora) que se llame por cada evento.
4.  **Importante:** Solo hacer `CommitMessages` si `processMessage` no retorna error.
```

### 4.2 Lógica de Combate y Persistencia

Procesamos el golpe y actualizamos la base de datos.

**Prompt para Lógica de Juego:**

```markdown
TAREA: Implementación de Lógica de Combate

Desarrolla la función `processMessage`.
Flujo:
1.  Deserializar JSON (`AttackEvent`).
2.  **Repo:** Buscar `Attacker` y `Target` en Postgres por ID.
3.  **Dominio:** Calcular daño (Fuerza del Atacante + Random(1-10) - Defensa del Objetivo).
4.  **Dominio:** Restar HP al objetivo.
5.  **Repo:** Actualizar el nuevo HP del objetivo en Postgres (UPDATE users SET hp = ...).

Tip: Usa una transacción de BD si necesitas actualizar múltiples tablas, pero por ahora una simple actualización basta.
Output: Loguear "X golpeó a Y causando Z daño. HP restante: W".
```

---

## 📂 Fase 5: Resiliencia y Observabilidad

Errores van a ocurrir. Necesitamos verlos y recuperarnos.

### 5.1 Logging Estructurado (slog)

Basta de `fmt.Println`. Queremos logs que una máquina pueda leer (JSON).

**Prompt para Logging:**

```markdown
TAREA: Configuración de Structured Logging

Instruye reemplazar todos los logs estándar por `slog` (Go 1.21+).
Requisitos:
1.  Formato JSON por defecto.
2.  Nivel de log configurable por variable de entorno (`LOG_LEVEL`).
3.  Atributos clave en cada log: `service_name`, `trace_id` (si está disponible), `error` (si aplica).

Ejemplo esperado:
`{"time":"...", "level":"INFO", "msg":"attack processed", "damage": 15, "target_id": "..."}`
```

### 5.2 Estrategia de Reintentos (Retries)

Si la base de datos parpadea, no queremos perder el evento del ataque.

**Prompt para Retries:**

```markdown
TAREA: Implementación de Backoff Exponencial

Modifica el Consumer para manejar errores transitorios (ej. conexión DB caída).
Lógica:
1.  Si `processMessage` falla, esperar 100ms y reintentar.
2.  Si falla de nuevo, esperar 200ms, luego 400ms (hasta 3 intentos).
3.  Si falla después de 3 intentos: Loguear ERROR CRÍTICO y descartar mensaje (después implementaremos Dead Letter Queue).
```

---

## 📂 Fase 6: Entrega y Migración (Handover)

Preparar el paquete para el futuro.

### 6.1 Generación de Documentación

Si no está documentado, no existe.

**Prompt para README:**

```markdown
TAREA: Generación de Documentación

Crea un `README.md` profesional para el repositorio.
Secciones:
1.  **Arquitectura:** Diagrama Mermaid simple.
2.  **Quick Start:**
    - `make up` (Levantar entorno).
    - `curl` de ejemplo para crear Player y Atacar.
3.  **Estructura:** Explicación breve de folders `internal/`.
```

### 6.2 El "Save State" (Handover Prompt)

El artefacto final de este tutorial. Un prompt para que *otro* agente entienda todo esto en 1 segundo.

**Prompt de Migración:**

```markdown
TAREA: Generar Prompt de Contexto (Handover)

Escribe un párrafo resumen que describa el estado actual del proyecto técnicamente.
Debe servir como "input" para una nueva sesión de chat con otra IA.

Debe incluir:
- Stack exacto (Go 1.22, Kafka-Go, Pgx).
- Estado de la arquitectura (Gateway HTTP -> Kafka -> Engine Consumer).
- Qué falta por hacer (ej. "Falta agregar sistema de inventario").
```

---

## ✅ ¡Misión Cumplida!
Si has seguido los prompts fase por fase, ahora tienes un sistema distribuidos funcional, documentado y listo para evolucionar.
```
