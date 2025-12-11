package main

import (
	// "context" es vital en Go para controlar tiempos de espera (timeouts)
	// y cancelaciones de tareas largas. Kafka-go lo usa intensivamente.
	"context"

	// "fmt" se usa para imprimir texto en la consola (stdout).
	"fmt"

	// "log" es similar a fmt, pero agrega fecha/hora y permite salir del programa
	// en caso de errores fatales (log.Fatal).
	"log"

	// "os" nos da acceso al sistema operativo, por ejemplo, para leer
	// los argumentos que escribes en la terminal (ej: 'produce' o 'consume').
	"os"

	// "time" nos sirve para pausar la ejecución (Sleep) y manejar duraciones.
	"time"

	// Importamos la librería externa de Kafka.
	// Asegúrate de correr 'go get github.com/segmentio/kafka-go' antes.
	"github.com/segmentio/kafka-go"
)

// --- CONFIGURACIÓN ---
// Definimos constantes para no repetir textos y facilitar cambios.
// En un entorno real, esto vendría de variables de entorno (.env).
const (
	topic         = "rpg-battles"    // El nombre del "chat" o canal donde van los datos.
	brokerAddress = "localhost:9094" // Dirección del servidor Kafka (Broker).
)

func main() {
	// Verificamos si el usuario ingresó algún comando.
	// os.Args[0] es el nombre del archivo, os.Args[1] es el primer argumento.
	if len(os.Args) < 2 {
		fmt.Println("Uso correcto: go run main.go [produce|consume]")
		os.Exit(1) // Salimos con código de error 1.
	}

	// Capturamos el modo que el usuario quiere ejecutar.
	mode := os.Args[1]

	switch mode {
	case "produce":
		produceMessages() // Ejecuta la función de enviar datos.
	case "consume":
		consumeMessages() // Ejecuta la función de leer datos.
	default:
		// Si escriben algo que no entendemos, avisamos y salimos.
		fmt.Printf("Modo desconocido: %s. Usa 'produce' o 'consume'\n", mode)
		os.Exit(1)
	}
}

// --- PRODUCTOR (El que envía los mensajes) ---
func produceMessages() {
	// Paso opcional pero recomendado: Asegurar que el topic exista antes de escribir.
	ensureTopic()

	// 1. CONFIGURACIÓN DEL WRITER (ESCRITOR)
	// El Writer es un componente de alto nivel que gestiona conexiones, reintentos
	// y balanceo de carga automáticamente.
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokerAddress), // A qué servidor conectarse.
		Topic: topic,                    // A qué tópico escribir.
		// Balancer: Decide a qué partición enviar el mensaje.
		// LeastBytes intenta enviar el mensaje a la partición que tenga menos datos,
		// ayudando a distribuir la carga equitativamente.
		Balancer: &kafka.LeastBytes{},
	}

	// 'defer' asegura que w.Close() se ejecute justo antes de que la función termine.
	// Es vital para cerrar conexiones de red y liberar memoria.
	defer w.Close()

	fmt.Println("⚔️  Iniciando Productor de Batallas...")

	// Simulamos el envío de 5 eventos.
	for i := 1; i <= 5; i++ {
		// Creamos el contenido del mensaje (payload).
		msgValue := fmt.Sprintf("Heroe ataca a Orco #%d con 50 de daño", i)

		// 2. ESCRIBIR EL MENSAJE
		// WriteMessages envía uno o más mensajes al broker.
		// context.Background() indica que no hay un tiempo límite específico para esta operación.
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				// Key: Kafka usa la llave para decidir el orden y la partición.
				// Mensajes con la misma Key siempre van a la misma partición en orden.
				Key: []byte(fmt.Sprintf("Key-%d", i)),

				// Value: Es la información real que queremos transmitir.
				// Kafka solo entiende bytes, por eso convertimos el string a []byte.
				Value: []byte(msgValue),
			},
		)

		// Manejo de errores básico.
		if err != nil {
			log.Fatal("Error fatal enviando mensaje:", err)
		}

		fmt.Printf("Enviado: %s\n", msgValue)

		// Esperamos 1 segundo para no saturar la pantalla y simular eventos en tiempo real.
		time.Sleep(1 * time.Second)
	}

	fmt.Println("✅ Todos los ataques enviados.")
}

// --- CONSUMIDOR (El que lee los mensajes) ---
func consumeMessages() {
	// También nos aseguramos que el topic exista, por si arrancamos el consumidor primero.
	ensureTopic()

	// 1. CONFIGURACIÓN DEL READER (LECTOR)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,

		// GroupID: ESTO ES CRUCIAL.
		// Identifica a este consumidor como parte de un equipo "battle-stats-service".
		// Kafka recuerda qué mensajes ya leyó este grupo para no repetirlos si el programa se reinicia.
		GroupID: "battle-stats-service",

		// Optimizaciones de red:
		MinBytes: 10e3, // Esperar a tener 10KB de datos antes de responder (menos tráfico).
		MaxBytes: 10e6, // Máximo 10MB por paquete.

		// StartOffset: Solo aplica si es un grupo nuevo sin historial.
		// FirstOffset significa "leer desde el mensaje más antiguo disponible".
		StartOffset: kafka.FirstOffset,
	})

	// Cerramos la conexión al terminar (aunque en un loop infinito, esto ocurre al matar el proceso).
	defer r.Close()

	fmt.Println("🛡️  Iniciando Consumidor de Batallas (Esperando eventos)...")

	// 2. BUCLE INFINITO DE LECTURA
	// Los consumidores suelen estar siempre encendidos escuchando.
	for {
		// ReadMessage es BLOQUEANTE.
		// El código se detiene en esta línea hasta que Kafka tenga un mensaje nuevo.
		m, err := r.ReadMessage(context.Background())

		if err != nil {
			// Si hay error (ej. desconexión momentánea), logueamos y esperamos un poco antes de reintentar.
			fmt.Printf("⚠️  Error leyendo mensaje: %v\n    --> Reintentando en 1s...\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Imprimimos el mensaje.
		// m.Offset es como el ID secuencial del mensaje dentro de la partición.
		fmt.Printf("Mensaje recibido: %s (offset %d)\n", string(m.Value), m.Offset)
	}
}

// --- UTILIDAD: CREACIÓN DE TOPICS ---
// Esta función usa una conexión "cruda" (bajo nivel) para administrar Kafka.
func ensureTopic() {
	// 1. Conexión inicial a cualquier broker disponible.
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatal("Error conectando a Kafka para verificar topic:", err)
	}
	defer conn.Close()

	// 2. Preguntar quién es el "Controller" (el jefe del clúster).
	// Solo el Controller tiene permiso para crear o borrar tópicos.
	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("Error obteniendo controlador:", err)
	}

	// 3. Conectarse directamente al Controller.
	var controllerConn *kafka.Conn
	// Construimos la dirección del controller usando su Host y Puerto.
	controllerConn, err = kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		log.Fatal("Error conectando al controlador:", err)
	}
	defer controllerConn.Close()

	// 4. Definir la configuración del tópico.
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1, // 1 Partición = Sin paralelismo de lectura (orden total garantizado).
			ReplicationFactor: 1, // 1 Copia = Sin redundancia (si el nodo cae, perdemos datos).
		},
	}

	// 5. Intentar crear el tópico.
	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		// Kafka devuelve error si el topic ya existe. En este ejemplo simple,
		// ignoramos ese error específico asumiendo que "si falla, es que ya estaba ahí".
		fmt.Printf("Nota: Topic '%s' ya existe o acaba de ser creado.\n", topic)
	}
}
