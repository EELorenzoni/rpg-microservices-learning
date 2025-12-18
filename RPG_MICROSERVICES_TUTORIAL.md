# 🧙‍♂️ RPG Microservices: From Prompt to Production
> **A Prompt-Driven Engineering Course**

Este documento es un tutorial maestro diseñado para ser consumido tanto por humanos como por Agentes de IA. Su objetivo es guiar la construcción de un sistema de RPG distribuido usando Go, Kafka y Event-Driven Architecture.

---

## 📂 Fase 0: El "Prime" (Configuración del Agente)

Esta fase no escribe código, sino comportamiento. Antes de empezar, copia y pega el siguiente prompt en tu chat con la IA. Esto configurará su "sys-call" mental y su rol como experto.

### 🤖 El Prompt Maestro del Sistema

Copiar el siguiente bloque y enviarlo a la IA para iniciar la sesión:

```markdown
ACTÚA COMO UN "ARQUITECTO SENIOR EN GO E INGENIERO DE PLATFORM ENGINEERING".

## TU OBJETIVO
Guiar al usuario en la construcción de un Backend RPG Distribuido usando Go, Apache Kafka y Event-Driven Architecture. Estamos haciendo "Desarrollo Guiado por Prompts": no solo escribirás código, sino que explicarás los *prompts* necesarios para generar ese código en el futuro o en otros agentes.

## RESTRICCIONES TÉCNICAS
1.  **Lenguaje:** Go 1.22+ (Tipado estricto, features modernas).
2.  **Arquitectura:** Hexagonal / Clean Architecture.
    - `cmd/`: Puntos de entrada.
    - `internal/core/domain`: Lógica y Modelos (Go Puro).
    - `internal/core/ports`: Interfaces (Repository, EventBus).
    - `internal/core/services`: Casos de uso (Vertical Slicing).
    - `internal/adapters`: Implementaciones concretas.
3.  **Comunicación:**
    - Síncrona: REST API (net/http estándar).
    - Asíncrona: Kafka (usando `github.com/segmentio/kafka-go`).
4.  **Platform Engineering:**
    - Kafka centralizado en proyecto separado
    - Admin API para gestión de topics
    - Configuración por variables de entorno (.env)
5.  **Event-Driven:**
    - Todos los eventos con estructura estándar: `event_type`, `occurred_at`, `data`
    - Publicar eventos de éxito Y fallo
    - Dead Letter Queue (DLQ) para resiliencia

## REGLAS DE COMPORTAMIENTO
1.  **Piensa Primero:** Antes de codear, define la estructura de archivos o el flujo lógico.
2.  **Event-First:** Define eventos de dominio antes de escribir handlers.
3.  **Educativo:** Explica *por qué* elegiste un patrón (ej. "¿Por qué usar DLQ?").
4.  **Iterativo:** Comienza con el MVP, luego refactoriza.
5.  **Idioma:** Todas tus explicaciones y comentarios deben ser en Español Latinoamericano.
6.  **SOLID es obligatorio:** Aplicar los 5 principios en todo el código.

## CONTEXTO ACTUAL
Estamos empezando desde cero. Espera instrucciones para la Fase 1.
```

---

## 📂 Fase 1: Platform Engineering - Kafka Centralizado

### 1.1 Creación del Proyecto Platform

**Prompt para Platform Admin:**

```markdown
TAREA: Crear Platform Kafka Admin

Crea un proyecto Go independiente llamado `platform-kafka-admin` que centralice la gestión de Kafka.

Estructura:
/cmd/admin-api/main.go     → API REST para gestionar topics
/internal/core/service.go  → Lógica de admin (CreateTopic, ListTopics, DeleteTopic)
/internal/handlers/http.go → Handlers HTTP (Gin)
/docker-compose.yml        → Kafka + Kafka UI
/.env                      → Variables de entorno
/Makefile                  → Automatización

Requisitos:
1. Kafka en modo KRaft (sin Zookeeper)
2. Admin API en puerto 3000
3. Endpoints:
   - POST /topics (crear topic)
   - GET /topics (listar topics)
   - DELETE /topics/:name (eliminar topic)
4. Configuración estricta por .env (si env var no existe, fallar)
5. Kafka UI en puerto 8080

Configuración de Kafka para desarrollo:
- KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
- KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR: 1
```

### 1.2 Configuración Profesional

**Prompt para Variables de Entorno:**

```markdown
TAREA: Implementar configuración con godotenv

Requisitos:
1. Instalar `github.com/joho/godotenv`
2. Leer archivo .env al inicio de main.go
3. Validar que existan las variables requeridas:
   - KAFKA_BROKER (dirección del broker)
   - ADMIN_PORT (puerto de la API)
4. Si alguna variable está vacía, fallar con log.Fatal explicativo
5. NO usar valores por defecto hardcodeados

Archivo .env debe contener:
KAFKA_BROKER=127.0.0.1:9094
ADMIN_PORT=:3000
```

---

## 📂 Fase 2: Hero Service - CRUD Event-Driven

### 2.1 Estructura Hexagonal

**Prompt para Estructura del Proyecto:**

```markdown
TAREA: Crear proyecto Hero Service (section-05-full-cycle)

Estructura Hexagonal completa:
/cmd/api/main.go              → HTTP Server
/cmd/consumer/main.go         → Kafka Consumer
/internal/core
    /domain/hero.go           → Entidad Hero + Factory
    /ports/repositories.go    → Interface HeroRepository
    /ports/events.go          → Interface EventBus
    /services/herosrv
        /service.go           → Struct + Dependencies
        /create.go            → Command: Create
        /get.go               → Query: Get
        /update.go            → Command: Update
        /delete.go            → Command: Delete
        /list.go              → Query: List
/internal/handlers/herohdl
    /http.go                  → REST Handlers
    /consumer.go              → Kafka Consumer Handler
/internal/repositories/herorepo
    /memory.go                → In-Memory Repository
    /kafka_repo.go            → Kafka EventBus Implementation

Reglas:
1. Vertical Slicing: Cada operación CRUD en archivo separado
2. CQS: Separar Commands (escribir) de Queries (leer)
3. Dependency Injection: Service recibe interfaces, no implementaciones
```

### 2.2 Domain Layer (Entidad Hero)

**Prompt para Domain:**

```markdown
TAREA: Implementar Entidad Hero con validaciones

Crear domain/hero.go con:

type Hero struct {
    ID        string
    Name      string
    Level     int
    Power     int
    CreatedAt time.Time
}

Reglas:
1. Factory Pattern: NewHero(id, name) que valide:
   - Name no puede estar vacío
   - ID debe ser válido
2. Retornar puntero (*Hero)
3. Errors de dominio predefinidos:
   - ErrHeroNameEmpty
4. Poder inicial: 10, Level inicial: 1
```

### 2.3 Generación Automática de IDs

**Prompt para UUID:**

```markdown
TAREA: Implementar generación automática de IDs

Requisitos:
1. Instalar `github.com/google/uuid`
2. En herosrv/create.go:
   - NO recibir ID en CreateHeroCommand
   - Generar ID con uuid.New().String()
   - Retornar el héroe creado (*domain.Hero, error)
3. En HTTP Handler:
   - Request JSON sin campo "id"
   - Response debe incluir el héroe completo con su ID generado

Ejemplo Response:
{
  "status": "created",
  "hero": {
    "id": "a1b2c3d4-...",
    "name": "Arthas",
    ...
  }
}
```

### 2.4 Event-Driven Architecture

**Prompt para Eventos:**

```markdown
TAREA: Implementar publicación de eventos con estructura estándar

Estructura de eventos:
{
  "event_type": "HeroCreated",
  "occurred_at": "2025-12-18T16:00:00Z",
  "data": {
    "id": "...",
    "name": "...",
    ...
  }
}

Tipos de eventos:
✅ Éxito:
- HeroCreated
- HeroUpdated
- HeroDeleted

❌ Fallo:
- HeroCreateFailed
- HeroUpdateFailed
- HeroDeleteFailed

Reglas:
1. SIEMPRE publicar eventos (tanto éxito como fallo)
2. En caso de error de validación, publicar evento de fallo ANTES de retornar error
3. En caso de éxito, publicar evento DESPUÉS de persistir en DB
4. Logs claros: "✅ Hero guardado en DB" → "📨 Evento 'HeroCreated' publicado correctamente"
```

### 2.5 Dead Letter Queue (DLQ)

**Prompt para Consumer Robusto:**

```markdown
TAREA: Implementar Consumer con DLQ

Crear internal/handlers/herohdl/consumer.go con:

1. DLQ Writer: Productor a topic "hero-events-05-dlq"
2. FetchMessage (NO ReadMessage) para control manual de commits
3. Función processMessage(msg kafka.Message) error que:
   - Retorne error si payload es inválido
4. Lógica de manejo:
   - Si processMessage falla → Enviar mensaje a DLQ con headers:
     * "original-topic"
     * "error-reason"
   - Hacer commit SIEMPRE (para avanzar, no bloquear)

Simulación de poison message:
Si payload == `{"fail":true}`, retornar error para probar DLQ
```

---

## 📂 Fase 3: Routing Inteligente (REST)

**Prompt para Router:**

```markdown
TAREA: Implementar routing RESTful inteligente

En cmd/api/main.go, crear lógica de routing:

Endpoint: /heroes

Lógica:
- Si query param "id" está presente:
  → Operaciones sobre UN héroe
  - GET    → GetHero
  - PUT    → UpdateHero
  - DELETE → DeleteHero

- Si query param "id" NO está presente:
  → Operaciones sobre la COLECCIÓN
  - POST → CreateHero
  - GET  → ListHeroes

Ejemplo:
POST /heroes {"name":"Arthas"}           → CreateHero
GET /heroes                              → ListHeroes
GET /heroes?id=abc-123                   → GetHero
PUT /heroes?id=abc-123 {"name":"Updated"} → UpdateHero
DELETE /heroes?id=abc-123                → DeleteHero
```

---

## 📂 Fase 4: Tutoriales Avanzados

### 4.1 Guía de Producción

**Prompt para Tutorial 07:**

```markdown
TAREA: Crear tutorial "Kafka en Producción"

Documento: 07-kafka-production-guide.md

Secciones:
1. Parámetros exhaustivos de topics:
   - min.insync.replicas
   - retention.ms
   - cleanup.policy (delete vs compact)
   - compression.type
2. Viaje de un evento (Producer → Broker → Consumer)
   - Diagrama de secuencia Mermaid
3. Consumer Groups explicados
   - Rebalanceo
   - Asignación de particiones
4. Semántica de entrega:
   - At-least-once
   - At-most-once
   - Exactly-once (con limitaciones reales)
5. Estrategias de error (DLQ, Retries)
6. Checklist de producción

Tono: Ingeniero Senior, sin marketing, con experiencia real operando Kafka.
```

### 4.2 Análisis de Flujo de Mensajes

**Prompt para Tutorial 08:**

```markdown
TAREA: Crear tutorial "Flujo del Mensaje Real"

Documento: 08-kafka-event-flow.md

Explicar usando el evento HeroCreated:
1. Anatomía de un mensaje Kafka:
   - Topic, Key, Value, Headers, Partition, Offset
2. Por qué la Key importa (ordenamiento, hot partitions)
3. Responsabilidad del Consumer:
   - Idempotencia (UPSERT, no INSERT)
   - Manejo de duplicados
4. Ejemplo JSON real del evento

Incluir advertencias:
- Orden solo existe DENTRO de una partición
- Duplicados son inevitables (network failures)
- Consumer debe ser idempotente
```

---

## 📂 Fase 5: Documentación Completa

**Prompt para Tutorial 05:**

```markdown
TAREA: Crear tutorial completo "Ciclo Completo y SOLID"

Documento: 05-ciclo-completo-solid.md

Contenidos:
1. Diagrama de secuencia Mermaid (dark theme):
   - Flujo Create con éxito
   - Flujo Create con fallo
   - Consumer con DLQ
2. Estructura del proyecto (Vertical Slicing)
3. Tabla SOLID con ejemplos concretos del código
4. Sección de pruebas con comandos curl:
   - Crear héroe (sin ID)
   - Listar héroes
   - Consultar uno
   - Actualizar
   - Eliminar
   - Probar fallo (name vacío)
   - Ver DLQ en acción
5. Logs esperados del Consumer

Nota: IDs se generan automáticamente, no se envían en requests
```

---

## 📂 Fase 6: Próximos Pasos (Battle System)

**Prompt para Sistema de Combate:**

```markdown
TAREA: Diseñar Sistema de Combate Asíncrono

Próximo servicio: Battle Service

API:
POST /battles
{
  "attacker_id": "uuid",
  "defender_id": "uuid"
}

Flujo:
1. API valida que ambos héroes existan
2. Publica evento "BattleStarted"
3. Consumer calcula:
   - Daño = Attacker.Power + Random(1-10) - Defender.Level
   - Actualiza Defender.HP
4. Publica "HeroAttacked" con resultado
5. Si Defender.HP <= 0, publica "HeroDefeated"

Eventos:
- BattleStarted
- HeroAttacked (con daño)
- BattleEnded (ganador)

Retos:
- Concurrent battles del mismo héroe
- Optimistic locking en DB
- Registro de historial de batallas
```

---

## ✅ ¡Misión Cumplida!

Si has seguido los prompts fase por fase, ahora tienes:
- ✅ Platform Engineering (Kafka centralizado)
- ✅ Hero Service (CRUD completo + eventos)
- ✅ Event-Driven Architecture
- ✅ Dead Letter Queue
- ✅ Documentación profesional
- ✅ Arquitectura Hexagonal + SOLID

**Próximo nivel:**
- [ ] Battle System
- [ ] PostgreSQL (migrations, queries)
- [ ] Tests (unit + integration)
- [ ] Observabilidad (slog, metrics, tracing)
