# 08 - El Flujo del Mesaje Real: Productor y Consumidor Expertos

> **Rol:** Principal Platform Engineer & Event-Driven Architect.
> **Objetivo:** Dejar de tratar a Kafka como una "caja negra" y entender qué demonios viaja por el cable.

Ya tienes la infraestructura (`07-kafka-production-guide.md`). Ahora hablemos de código.
Vamos a diseccionar UN solo evento: `HeroCreated`. No en abstracto, sino el que tu código Go está enviando ahora mismo.

---

## 1. El Evento Real: ¿Qué viaja por el cable?

Cuando tu `Hero Service` ejecuta:
```go
// internal/core/services/herosrv/create.go
s.eventBus.Publish(hero, "HeroCreated")
```

Lo que Kafka recibe NO es un objeto Go. Son **BYTES**.
Un mensaje de Kafka tiene esta anatomía sagrada:

1.  **Topic**: `hero-events-05` (El buzón).
2.  **Key**: `hero-123` (El ID del héroe).
3.  **Value** (Payload): El JSON serializado.
4.  **Headers**: Metadatos (como headers HTTP: `trace-id`, `source`, etc).

**Ejemplo Real del Payload (JSON)**:
```json
{
  "event_type": "HeroCreated",
  "occurred_at": "2025-01-10T12:45:22Z",
  "data": {
    "id": "hero-123",
    "name": "Arthas",
    "power": 9001
  }
}
```

> **Ingeniería**: ¿Ves el campo `event_type`? Es obligatorio. Un topic suele tener mezclados eventos `HeroCreated`, `HeroUpdated`, `HeroDeleted`. El consumidor necesita saber qué deserializar.

---

## 2. Terminología Kafka (Explicada con Héroes)

Olvídate de las definiciones de libro.

### A. Offset
-   **Mito**: "Es el ID del mensaje".
-   **Realidad**: Es la posición **en esa partición**.
-   **Ejemplo**: "El evento de Arthas está en el casillero 1050 de la fila 0".
-   **Importancia**: Si tu consumer cae y reinicia, le dice a Kafka: *"Me quedé en 1050, dame el 1051"*.

### B. Partición
-   **Situación**: Tienes 3 particiones y mandas 100 héroes.
-   **Key = null**: Kafka tira los héroes al azar (Round Robin). Arthas va a la 0, Jaina a la 1.
-   **Key = heroID**: Kafka usa una fórmula matemática (`hash(heroID) % 3`).
    -   Todos los eventos de Arthas SIEMPRE caerán en la partición 2.
    -   **Garantiza Orden**: No procesarás "Arthas Muere" antes que "Arthas Nace".

### C. Commit
-   **Concepto**: Es el "Check" de la lista de tareas.
-   **Flujo**:
    1.  Leo "Arthas Nace" (Offset 1050).
    2.  Guardo en DB Query.
    3.  **Commit 1050** -> "Kafka, ya terminé con Arthas, anótalo".

---

## 3. Productor Experto: No rompas el sistema

Tu trabajo como productor es elegir bien la **KEY**.

### La Regla de Oro de la Key
> **"Eventos que deben procesarse en orden, DEBEN tener la misma Key."**

Si pones `Key: null`:
-   Ganas velocidad (balanceo perfecto).
-   Pierdes orden (consumer A procesa offset 10 partición 0, consumer B procesa offset 11 partición 1... al mismo tiempo).

**Consecuencias de una mala Key**:
-   **Key = "static" (String fijo)**: Todas las peticiones van a la partición 0. Tienes 10 brokers, pero solo usas 1. Eso es un **Hot Partition**. Tu sistema colapsará.
-   **Key = Random UUID**: Pierdes el orden de eventos de una misma entidad.

**Veredicto**: Usa `hero.ID` como Key. Balancea bien (millones de héroes) y ordena por héroe.

---

## 4. Consumer Experto: Cómo leer

Cuando tu código Go (`segmentio/kafka-go`) lee, recibe un struct `kafka.Message`.
Míralo con ojos de rayos X:

```go
msg.Topic     // "hero-events-05"
msg.Partition // 0
msg.Offset    // 1050
msg.Key       // []byte("hero-123")  <-- CRÍTICO
msg.Value     // []byte(`{"event_type":"HeroCreated"...}`) <-- PAYLOAD
```

### Tu Responsabilidad como Consumer

1.  **Leer la Key**: A veces la info crítica está en la Key y no en el JSON.
2.  **Ignorar lo desconocido**: Si llega `event_type: "HeroDanced"` y no sabes qué hacer, **ignóralo y haz commit**. No crashees. El producer evolucionó y tú no. Es normal.
3.  **Idempotencia**:
    -   Kafka te manda el mensaje offset 1050.
    -   Tú guardas a Arthas en Postgres.
    -   Justo antes de hacer `Commit`, se corta la luz.
    -   Reinicias. Kafka te manda el 1050 **otra vez**.
    -   **¿Qué haces?**
        -   Junior: `INSERT INTO heroes...` -> **Error: Duplicate Key**. Crash.
        -   Senior: `INSERT ... ON CONFLICT DO NOTHING` o `UPSERT`. -> **Éxito**.

---

## 5. Resumen del Flujo de un Mensaje

1.  **API**: Recibe POST.
2.  **Producer**: Serializa JSON. Calcula `hash("hero-123")`. Asigna a Partición 2. Envía.
3.  **Kafka (Broker)**: Recibe. Escribe en disco en cola de Partición 2. Offset asignado: 1050.
4.  **Consumer**: Pide "dame mensajes nuevos de Partición 2". Recibe Offset 1050.
5.  **Lógica**: Deserializa JSON. Verifica si ya existe en DB. Guarda.
6.  **Ack**: Envía "Commit Offset 1050" a Kafka.

Así es como fluye la sangre de tu sistema distribuido.
