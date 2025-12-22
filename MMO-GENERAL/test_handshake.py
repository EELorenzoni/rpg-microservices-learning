import socket
import struct

# =============================================================================
# EXPLICACIÓN DEL TEST DE HANDSHAKE (SALUDO INICIAL)
# =============================================================================
# Antes de que un jugador pueda moverse, el servidor debe saber quién es.
# Este script simula el "primer contacto" del cliente con el servidor.
#
# OBJETIVOS DEL TEST:
# 1. Registro: ¿El servidor guarda la IP y puerto del nuevo jugador?
# 2. Asignación de ID: ¿El servidor devuelve un PlayerID único (ej. 1001)?
# 3. Respuesta UDP: ¿El servidor es capaz de responder al puerto correcto del cliente?
# =============================================================================

# Configuración del servidor Ubuntu
SERVER_IP = "192.168.0.100"
SERVER_PORT = 8080

def test_handshake():
    # 1. Crear el socket UDP (User Datagram Protocol)
    # Es el protocolo estándar para juegos MMO por su baja latencia.
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.settimeout(2.0)

    # 2. Preparar el paquete de Handshake binario
    # Estructura de cabecera definida en la Fase 0:
    # byte 0: Tipo de paquete (0 = Handshake)
    # bytes 1-4: Secuencia (0 en el primer saludo)
    # bytes 5-12: PlayerID (0 porque aún no tenemos uno)
    #
    # '>BIQ' usa el estándar "Big Endian" para que el servidor lo entienda.
    packet_type = 0
    sequence = 0
    player_id = 0
    
    packet = struct.pack(">BIQ", packet_type, sequence, player_id)
    
    print(f"📡 Enviando Handshake a {SERVER_IP}:{SERVER_PORT}...")
    print(f"📦 Contenido del paquete (hex): {packet.hex()}")

    try:
        # 3. Enviar los 13 bytes al servidor
        sock.sendto(packet, (SERVER_IP, SERVER_PORT))

        # 4. Esperar la respuesta (bloqueante hasta recibir datos o timeout)
        data, addr = sock.recvfrom(1024)
        print(f"✅ Respuesta recibida de {addr}")
        
        # 5. Deserializar la respuesta del servidor
        # El servidor nos devuelve el mismo formato de 13 bytes,
        # pero con el PlayerID que nos ha asignado.
        res_type, res_seq, res_id = struct.unpack(">BIQ", data[:13])
        
        print("\n--- RESULTADO DEL SERVIDOR ---")
        print(f"Tipo de Paquete: {res_type} (Confirmación de Handshake)")
        print(f"Secuencia: {res_seq}")
        print(f"TU PLAYER ID ASIGNADO: {res_id}")
        print("-----------------------------\n")

    except socket.timeout:
        print("❌ ERROR: El servidor no respondió. ¿Ejecutaste el servidor en modo Debug (F5)?")
    finally:
        sock.close()

if __name__ == "__main__":
    test_handshake()
