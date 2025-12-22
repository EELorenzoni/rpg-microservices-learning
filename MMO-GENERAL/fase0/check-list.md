# Checklist de Ejecución: Fase 0 - Fundaciones

Este documento registra el progreso real del proyecto y sirve como guía de tareas pendientes para alcanzar el hito de Comunicación de Baja Latencia.

## 🟢 1. Infraestructura y Acceso (Hardware/OS)

| Tarea | Estado | Observaciones |
| :--- | :---: | :--- |
| Configuración de Hardware "CHORI" | ✅ DONE | AMD Ryzen / 16GB RAM / SSD 500GB. |
| Instalación de Ubuntu Server 24.04 | ✅ DONE | Instalación limpia y optimizada. |
| Configuración de Red Local (IP Estática) | ✅ DONE | IP: 192.168.0.100. |
| Acceso Remoto vía SSH (Mac -> PC) | ✅ DONE | Llaves y túnel de comunicación validados. |
| Despliegue de Docker & Docker Compose | ✅ DONE | Motor de servicios de soporte activo. |

## 🟡 2. Entorno de Desarrollo (Servidor)

- [x] **Configuración VS Code Remote SSH**: Conectar VS Code (Mac) al servidor (Ubuntu) para editar remotamente.
- [x] **Instalación de Go 1.21+**: Necesario para compilar el Game Server nativamente en el servidor.
- [x] **Estructura de Carpetas**: Crear `/home/superchori/mmo-server/cmd` y `/internal`.
- [x] **Inicialización de Módulo**: Ejecutar `go mod init` para gestión de dependencias.

## 🔴 3. Desarrollo del Game Server (Go)

- [x] **Socket Listener (UDP)**: Abrir puerto 8080 y manejar el buffer de entrada.
- [x] **Bucle de Simulación (Tick Rate)**: Implementar el Ticker a 33.3ms constante.
- [x] **Gestor de Conexiones (RAM)**: Mapa concurrente para registrar PlayerID y direcciones IP.
- [x] **Lógica de Handshake**: Responder al cliente con un Ack y asignar un ID de sesión.
- [x] **Replicación de Movimiento**: Recibir coordenadas y hacer broadcast al resto de los conectados.

## ⚪ 4. Integración con el Cliente (Unreal Engine)

- [ ] **Módulo de Red C++**: Crear la clase `MMONetworkClient` usando `FUdpSocketBuilder`.
- [ ] **Serialización Binaria**: Función para empaquetar coordenadas en bytes (no JSON).
- [ ] **Interpolación de Movimiento**: Lógica en UE para suavizar la posición entre paquetes recibidos.

## 📈 5. Criterios de Validación Final

- [ ] **Prueba de Latencia**: Verificar <10ms en red local.
- [ ] **Estabilidad del Tick**: Confirmar que el servidor mantiene los 30Hz sin variaciones (jitter).
- [ ] **Prueba de Concurrencia**: Conectar 2 instancias de Unreal y que se vean moverse mutuamente.

---

> **Próxima Acción Inmediata:** Definir las estructuras C++ en Unreal Engine y crear el componente `MMONetworkClient` para iniciar el Handshake con el servidor.