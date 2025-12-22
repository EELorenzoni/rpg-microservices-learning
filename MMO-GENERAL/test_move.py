import socket
import struct
import time

# =============================================================================
# EXPLICACIÓN DEL TEST DE MOVIMIENTO (REPLICACIÓN)
# =============================================================================
# En un MMO, el servidor es el "juez" (Autoridad).
# Este script simula lo que hará Unreal Engine cuando tu personaje se mueva.
#
# OBJETIVOS DEL TEST:
# 1. Deserialización: ¿Puede el servidor Go recibir bytes y convertirlos en números (X, Y, Z)?
# 2. Handshake: Validar que el servidor reconoce al jugador antes de aceptar su movimiento.
# 3. Protocolo Binario: Validar que el orden de los bytes (Big Endian) es el correcto.
# =============================================================================

SERVER_IP = "192.168.0.100"
SERVER_PORT = 8080

def test_movement():
    # Creamos un socket UDP (el mismo protocolo que usa Unreal para red rápida)
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.settimeout(2.0)

    # -------------------------------------------------------------------------
    # PASO 1: HANDSHAKE (Saludo inicial)
    # -------------------------------------------------------------------------
    # Enviamos 13 bytes vacíos pero con Tipo 0 para pedir un PlayerID.
    # struct.pack ">BIQ" significa: 
    #   > : Big Endian (Estándar de red)
    #   B : unsigned char (1 byte) - Tipo de paquete (0: Handshake)
    #   I : unsigned int  (4 bytes) - Secuencia (0 por ahora)
    #   Q : unsigned long long (8 bytes) - PlayerID (0 porque somos nuevos)
    packet = struct.pack(">BIQ", 0, 0, 0)
    
    print(f"📡 Paso 1: Enviando Handshake a {SERVER_IP}...")
    sock.sendto(packet, (SERVER_IP, SERVER_PORT))
    
    # El servidor nos responde con nuestro ID asignado
    data, addr = sock.recvfrom(1024)
    res_type, res_seq, player_id = struct.unpack(">BIQ", data[:13])
    print(f"✅ Handshake exitoso. El servidor nos asignó el ID: {player_id}")

    # -------------------------------------------------------------------------
    # PASO 2: ENVÍO DE COORDENADAS (Movimiento)
    # -------------------------------------------------------------------------
    # Simulamos que el personaje está en una posición específica de Unreal.
    # Payload: X, Y, Z (Posición) y Yaw (Rotación horizontal).
    # Cada número es un 'float' (4 bytes). 4 floats = 16 bytes.
    x, y, z, yaw = 100.5, 200.0, 50.25, 90.0
    
    # Construimos la CABECERA (13 bytes) indicando que es tipo 1 (Move)
    move_header = struct.pack(">BIQ", 1, 1, player_id)
    
    # Construimos el PAYLOAD (16 bytes) con las coordenadas
    # 'ffff' significa formatear 4 floats
    move_payload = struct.pack(">ffff", x, y, z, yaw)
    
    # El paquete total que viaja por el cable son 29 bytes (13 + 16)
    full_packet = move_header + move_payload
    
    print(f"\n📡 Paso 2: Enviando posición simulada (X:{x}, Y:{y}, Z:{z})...")
    sock.sendto(full_packet, (SERVER_IP, SERVER_PORT))
    
    print("\n-------------------------------------------------------------")
    print("💎 PRUEBA TERMINADA")
    print("Si el Debugger de Antigravity saltó, fíjate en Go si los datos")
    print(f"coinciden con X:{x}, Y:{y}, Z:{z}")
    print("-------------------------------------------------------------")
    
    sock.close()

if __name__ == "__main__":
    test_movement()
