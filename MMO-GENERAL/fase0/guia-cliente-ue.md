# Guía de Implementación: Cliente Unreal Engine (Fase 0)

Para conectar Unreal Engine con nuestro servidor de Go, utilizaremos los sockets nativos de UE (`Networking` y `Sockets` modules).

## 1. Configuración de Estructuras (C++)

Debemos asegurar que el empaquetado sea idéntico al de Go (Big Endian).

```cpp
#pragma pack(push, 1) // Asegura que no haya padding entre campos
struct FMMOPacketHeader {
    uint8 Type;        // 0: Handshake, 1: Move
    uint32 Sequence;   // ID incremental
    uint64 PlayerID;   // Asignado por el servidor
};

struct FMMOMovePayload {
    float X;
    float Y;
    float Z;
    float Yaw;
};
#pragma pack(pop)
```

## 2. Creación del Socket (Snippet)

En tu servicio o clase de red (`AMMONetworkClient`), inicializa el socket así:

```cpp
#include "Networking.h"
#include "Sockets.h"
#include "SocketSubsystem.h"

// ... en el constructor o Initialice()
FSocket* LocalSocket = FUdpSocketBuilder(TEXT("MMOGameSocket"))
    .AsNonBlocking()
    .AsReusable()
    .WithReceiveBufferSize(1024 * 1024);

FIPv4Address ServerIP;
FIPv4Address::Parse(TEXT("192.168.0.100"), ServerIP);
TSharedRef<FInternetAddr> ServerAddr = ISocketSubsystem::Get(PLATFORM_SOCKETSUBSYSTEM)->CreateInternetAddr();
ServerAddr->SetIp(ServerIP.Value);
ServerAddr->SetPort(8080);
```

## 3. Envío de Handshake

```cpp
void SendHandshake() {
    FMMOPacketHeader HandshakePacket;
    HandshakePacket.Type = 0;
    HandshakePacket.Sequence = 0;
    HandshakePacket.PlayerID = 0; // Inicialmente 0

    int32 BytesSent = 0;
    LocalSocket->SendTo((uint8*)&HandshakePacket, sizeof(HandshakePacket), BytesSent, *ServerAddr);
    
    UE_LOG(LogTemp, Log, TEXT("MMO: Handshake enviado al servidor."));
}
```

## 4. Notas Importantes para Unreal
1. **Módulos**: Asegúrate de agregar `"Networking"` y `"Sockets"` en tu archivo `.Build.cs`.
2. **Big Endian**: Go usa Big Endian por defecto en la red. Si estás en una CPU Little Endian (Mac/PC), deberás usar funciones como `FGenericPlatformMemory::WriteUintX` o `ByteSwap` para que los datos lleguen correctamente al servidor si no usas serializadores de UE.
3. **Frecuencia**: Ejecuta el envío de paquetes en un `FTimerHandle` o en el `Tick` de un componente específico limitado a 30 FPS para no saturar el socket.
