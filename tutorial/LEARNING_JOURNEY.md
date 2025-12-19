# Diario de Aprendizaje: AI Mastery & Arquitectura de Microservicios
**Proyecto:** RPG Distribuido con Go y Kafka  
**Autor:** Chori  
**Mentor AI:** Antigravity

Este documento sirve como bitácora de aprendizaje. Aquí documentaremos tanto *cómo* usamos la IA efectivamente, como los conceptos técnicos del proyecto.

---

## Parte 1: Dominando la IA Generativa (Meta-Aprendizaje)

### 1.1 El Cambio de Mentalidad
Para aprovechar la IA al máximo, cambiamos el enfoque de "Codificador Solitario" a **"Arquitecto Técnico y Revisor"**.
*   **Antes:** Pensar la lógica -> Escribir la sintaxis -> Debuggear typos.
*   **Ahora:** Definir la arquitectura/objetivo -> Describir el *intent* (intención) a la IA -> Revisar y validar la solución -> Iterar.

### 1.2 Principios de Prompting Efectivo (Ingeniería de Instrucciones)
1.  **Contexto Rico:** Darle un ROL a la IA (ej. "Senior Go Developer").
2.  **Chain of Thought (Cadena de Pensamiento):** Pedirle que primero *planifique* antes de *ejecutar*.
3.  **Iteración:** No aceptar el primer resultado ciegamente. Pedir optimizaciones o explicaciones.

---

## Parte 2: Arquitectura del Proyecto (El Objetivo Técnico)

Construimos un **Sistema de Gestión de Héroes (RPG)** con arquitectura empresarial moderna, event-driven y escalable.

### El Stack Tecnológico Actual

1.  **REST API**: Endpoints HTTP para CRUD de héroes
2.  **Apache Kafka (Event Bus)**: Sistema nervioso del sistema
    - Publica eventos de éxito y fallo
    - Dead Letter Queue (DLQ) para resiliencia
3.  **Platform Engineering**: Kafka centralizado
4.  **Microservicios en Go**:
    - **Hero Service (section-05-full-cycle)**: CRUD completo con eventos
    - **Platform Kafka Admin**: Gestión centralizada de infraestructura
5.  **Patrones Aplicados**:
    - Hexagonal Architecture (Ports & Adapters)
    - SOLID Principles
    - CQS (Command Query Separation)
    - Event-Driven Architecture

---

## Parte 3: Lo que Hemos Construido ✅

### ✅ Completado

#### 1. **Platform Kafka Admin** (`projects/platform-kafka-admin`)
- Infraestructura centralizada de Kafka
- Admin API (REST) para crear/eliminar topics
- Kafka UI para visualización
- Configuración profesional (.env, validación estricta)
- **Guía**: `tutorial/platform-kafka-admin-guide.md`

#### 2. **Hero Service** (`projects/section-05-full-cycle`)
- **CRUD Completo**:
  - Create (POST /heroes) - ID auto-generado con UUID
  - Read (GET /heroes?id=...)
  - Update (PUT /heroes?id=...)
  - Delete (DELETE /heroes?id=...)
  - List (GET /heroes)
- **Event-Driven**:
  - Eventos de éxito: `HeroCreated`, `HeroUpdated`, `HeroDeleted`
  - Eventos de fallo: `HeroCreateFailed`, `HeroUpdateFailed`, `HeroDeleteFailed`
  - Estructura estándar: `event_type`, `occurred_at`, `data`
- **Consumer Robusto**:
  - Dead Letter Queue (DLQ) para mensajes venenosos
  - Logging detallado (Partition, Offset, Key)
  - Nunca se bloquea
- **Arquitectura**:
  - Hexagonal (Ports & Adapters)
  - Vertical Slicing por operación
  - Dependency Injection
- **Tutorial**: `tutorial/05-ciclo-completo-solid.md`

#### 3. **Battle System** (`projects/section-06-battle-system`)
- **Concepto**: Combate asíncrono entre héroes.
- **Mecánica**: Basada en eventos de turnos para mayor escalabilidad.
- **Integración**: HTTP Client a Hero Service + Kafka para orquestación.
- **Guía**: `tutorial/06-battle-system.md`

#### 4. **Conceptos Enseñados**
- Kafka: Topics, Partitions, Replicas, Offsets, Consumer Groups
- Event Sourcing básico
- Consistency models (At-least-once, exactly-once)
- Platform Engineering
- 12-Factor App (Configuración por ENV)

---

## Parte 4: Próximos Pasos (Roadmap)

### 🎯 Fase Siguiente: Battle System (Combate)

El siguiente paso natural es implementar el **sistema de combate asíncrono** que justifica toda la arquitectura de eventos.

#### **Servicio de Combate** (Próximo)
- **Endpoint**: `POST /battles` (Iniciar combate entre 2 héroes)
- **Lógica**:
  - Calcular daño basado en stats
  - Turnos asíncronos vía Kafka
  - Actualizar HP de héroes
- **Eventos**:
  - `BattleStarted`
  - `HeroAttacked` (con daño calculado)
  - `BattleEnded` (ganador/perdedor)
- **Consumer**: Escucha batallas y actualiza estado de héroes

#### **Inventario** (Futuro)
- Sistema de items
- Equipar/desequipar
- Eventos de cambio de stats

#### **Persistencia Real** (Evolución)
- Migrar de Memory a PostgreSQL
- Implementar `herorepo.Postgres`
- Migrations con `goose` o `migrate`

#### **Observabilidad** (Producción)
- Structured logging con `slog`
- Metrics con Prometheus
- Distributed tracing

---

## Aprendizajes Clave

1. **Platform Engineering > Microservices individuales**: Centralizar infraestructura (Kafka) evita caos
2. **Events > Requests**: La comunicación asíncrona desacopla y escala mejor
3. **DLQ es obligatorio**: Los mensajes venenosos NO deben bloquear el sistema
4. **Siempre publicar eventos**: Tanto éxito como fallo (observabilidad completa)
5. **UUIDs > IDs manuales**: Generación automática evita colisiones
6. **SOLID no es teoría**: Es supervivencia en proyectos reales

---

## Métricas del Proyecto

- **Proyectos**: 2 (Platform Admin, Hero Service)
- **Tutoriales**: 8 documentos Markdown
- **Endpoints REST**: 6 (CRUD + List + Platform Admin)
- **Tipos de Eventos Kafka**: 6 (3 success, 3 failure)
- **Patterns**: Hexagonal, SOLID, CQS, Event-Driven, DLQ
- **Lenguaje**: 100% Go
- **Tests**: Pendiente (próxima iteración)

---

¡Felicidades! Has construido una arquitectura profesional desde cero. 🚀
