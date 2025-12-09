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

> **Nota:** Para asegurar una comprensión profunda, incluiremos diagramas visuales (Arquitectura, Secuencia y Estados) en cada fase crítica.

Construiremos un **Juego de Rol (RPG)** basado en texto, pero con una arquitectura empresarial moderna y escalable.

### El Stack Tecnológico
1.  **REST API (Gateway):** La puerta de entrada. Recibe las acciones del jugador (ej. `POST /attack`).
2.  **Apache Kafka (Event Bus):** El sistema nervioso. Desacopla los servicios. Si un jugador ataca, se emite un evento `PlayerAttacked`.
3.  **Microservicios en Go:**
    *   **Gateway Service:** Recibe HTTP, valida y publica en Kafka.
    *   **Game Engine Service:** Escucha eventos, calcula daño/lógica, actualiza estado.
    *   **Notification Service:** (Opcional) Escucha eventos para notificar al usuario (WebSocket/Email).

---

## Parte 3: Diseño del RPG (La Motivación)

**Concepto:** "Go Warriors: The Distributed Dungeon"

### Mecánicas Básicas para el MVP (Producto Mínimo Viable)
*   **Creación de Personaje:** Endpoint para crear un héroe.
*   **Exploración:** Moverse entre "habitaciones" (nodos).
*   **Combate:** Sistema de turnos asíncrono (gracias a Kafka).

---

## Próximos Pasos (To-Do)

- [ ] **Diseño Visual & Arquitectura:**
    - [ ] **Diagrama de Arquitectura:** Vista de alto nivel de los microservicios y Kafka (Mermaid.js o ASCII).
    - [ ] **Diagrama de Secuencia:** Detalle del flujo de eventos (ej. `HTTP Request` -> `Kafka Produce` -> `Consume`).
    - [ ] **Diagrama de Estados:** Ciclo de vida del combate o del personaje si amerita.
- [ ] **Configuración de Entorno:**
    - [ ] Instalar Command Line Tools (fix `git`).
    - [ ] Instalar Go (Golang).
    - [ ] Instalar Docker (para correr Kafka fácilmente).
- [ ] **Diseño de APIs:** Definir los endpoints en OpenAPI/Swagger.
- [ ] **Setup de Kafka:** Levantar un cluster local.
