# Unidad 2: Conceptos Teóricos - Kafka y Go

Antes de escribir código, entendamos las herramientas. Vamos a explicar esto "para Juniors, diseñado por Seniors".

## 🧠 Kafka: No es una base de datos, es un tronco (Log)
Imagina Kafka no como una caja donde guardas cosas (database), sino como una **cinta transportadora** infinita o una tubería.

### Conceptos Clave
1.  **Eventos**: Son mensajes. Algo que pasó. Ej: `{ "tipo": "DAÑO", "valor": 50, "target": "Orco" }`.
2.  **Topics (Tópicos)**: Son las etiquetas de la tubería. Un topic podría ser `world-events`. Todos los eventos del mundo van ahí.
3.  **Producer (Productor)**: El que grita el mensaje. "¡He golpeado al orco!".
4.  **Consumer (Consumidor)**: El que escucha. "Oh, alguien golpeó al orco, le bajaré la vida".

> **Analogía**: Twitter (X).
> - **Producer**: Tú escribiendo un tweet.
> - **Topic**: El hashtag #RPG.
> - **Consumer**: Alguien siguiendo ese hashtag.

## 🐹 Go: Concurrencia Nativa
Go es perfecto para esto porque maneja "hacer muchas cosas a la vez" de forma nativa.

- **Goroutines**: Son como hilos (threads) pero ultra ligeros. Podemos tener miles de "trabajadores" escuchando eventos sin que la computadora sude.
- **Channels**: Son tuberías internas de Go.

### ¿Cómo se unen?
Nuestra arquitectura será así:

```mermaid
graph TD
    User((Usuario))
    
    subgraph "Tu Computadora (Localhost)"
        Producer[Go Producer]
        Consumer[Go Consumer]
        
        subgraph "Docker"
            Kafka{"Kafka (KRaft)"}
        end
    end

    User -- "1. go run... produce" --> Producer
    User -- "2. go run... consume" --> Consumer
    Producer -- "3. Envía Evento" --> Kafka
    Kafka -- "4. Notifica" --> Consumer
    
    style Kafka fill:#f9f,stroke:#333,stroke-width:2px
    style Producer fill:#bbf,stroke:#333
    style Consumer fill:#bfb,stroke:#333
```

1.  Un microservicio (Producer) recibe una acción (o comando CLI).
2.  Envía el evento a Kafka.
3.  Otro microservicio (Consumer) ve el evento en Kafka y reacciona.

### Flujo de Mensajes (Sequence Diagram)

```mermaid
sequenceDiagram
    actor U as Usuario
    participant P as Producer (Go)
    participant K as Kafka (Docker)
    participant C as Consumer (Go)

    Note over U, C: Flujo de una Batalla RPG

    U->>P: Ejecuta comando attack
    activate P
    P->>P: Crea mensaje JSON (Heroe ataca...)
    P->>K: PUSH "Battle Event" (Topic: rpg-battles)
    deactivate P
    
    Note right of K: Kafka guarda el evento en disco
    
    activate C
    K-->>C: PULL (Nuevo Evento Disponible)
    C->>C: Procesa daño (Log)
    C-->>U: Imprime resultado en consola
    deactivate C
```

---
**🚀 Siguiente Paso**: Vamos a implementar nuestro primer Productor y Consumidor en Go para ver esto en acción.
