package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// Constantes de configuración
const (
	topic         = "rpg-battles"
	brokerAddress = "localhost:9094"
)

func main() {
	// os.Args es como process.argv en Node.js
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go [produce|consume]")
		os.Exit(1) // process.exit(1)
	}

	mode := os.Args[1]

	switch mode {
	case "produce":
		produceMessages()
	case "consume":
		consumeMessages()
	default:
		fmt.Printf("Modo desconocido: %s. Usa 'produce' o 'consume'\n", mode)
		os.Exit(1)
	}
}

// produceMessages simula un cliente enviando eventos
func produceMessages() {
	// 0. Asegurar que el topic existe antes de escribir
	ensureTopic()

	// 1. Crear el escritor (writer)
	// En Node.js esto sería como inicializar la instancia del cliente de KafkaJS producer.
	// &kafka.Writer{...} crea un puntero a una estructura (como crear un 'new Class' con config object).
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokerAddress),
		Topic: topic,
		// Balancer: Cómo elegimos a qué partición va el mensaje.
		// LeastBytes es inteligente: manda al que tenga menos datos pendientes.
		Balancer: &kafka.LeastBytes{},
	}

	// defer es mágico: asegura que w.Close() se ejecute cuando esta función termine,
	// pase lo que pase (incluso si hay errores). Es como un bloque 'finally' automático.
	defer w.Close()

	fmt.Println("⚔️  Iniciando Productor de Batallas...")

	for i := 1; i <= 5; i++ {
		msgValue := fmt.Sprintf("Heroe ataca a Orco #%d con 50 de daño", i)

		// 2. Escribir mensaje
		// context.Background() es el contexto base, nunca expira.
		// En Go, pasamos 'Context' a todo lo que lleva tiempo (I/O, DB, API calls)
		// para poder cancelarlo si tarda mucho (timeout).
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				// Key: Kafka usa esto para decidir el orden. Mensajes con la misma Key
				// van siempre a la misma partición y mantienen orden estricto.
				Key: []byte(fmt.Sprintf("Key-%d", i)),
				// Value: El contenido real. Kafka solo entiende bytes (Buffer en Node), no strings ni JSON directo.
				Value: []byte(msgValue),
			},
		)

		// En Go, los errores son valores, no excepciones. Siempre chequeamos 'if err != nil'.
		if err != nil {
			log.Fatal("Error enviando mensaje:", err) // log.Fatal imprime y hace os.Exit(1)
		}

		fmt.Printf("Enviado: %s\n", msgValue)
		time.Sleep(1 * time.Second) // Pausa de 1s para ver el efecto en tiempo real
	}

	fmt.Println("✅ Todos los ataques enviados.")
}

// consumeMessages simula un microservicio escuchando eventos
func consumeMessages() {
	ensureTopic()

	// 1. Crear el lector (reader)
	// En KafkaJS esto sería 'consumer.subscribe({ topic })' + 'consumer.run(...)'.
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		// GroupID es CRÍTICO:
		// - Si dos servicios tienen el mismo GroupID, Kafka reparte la carga entre ellos (Load Balancing).
		// - Si tienen diferente GroupID, ambos reciben COPIAS de todos los mensajes (Pub/Sub Broadcast).
		GroupID:     "battle-stats-service",
		MinBytes:    10e3,              // 10KB
		MaxBytes:    10e6,              // 10MB
		StartOffset: kafka.FirstOffset, // FirstOffset = leer desde el principio de la historia (replay)
	})
	defer r.Close()

	fmt.Println("🛡️  Iniciando Consumidor de Batallas (Esperando eventos)...")

	// Loop infinito (como un 'while(true)' o el event loop de Node mantenido vivo)
	for {
		// ReadMessage bloquea la ejecución hasta que llega algo.
		// Es sincrónico pero no bloquea el hilo del sistema gracias a las Goroutines de Go.
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("⚠️  Error leyendo mensaje: %v\n    --> Reintentando en 1s...\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		// string(m.Value) convierte el buffer de bytes a texto legible
		fmt.Printf("Mensaje recibido: %s (offset %d)\n", string(m.Value), m.Offset)
	}
}

// ensureTopic crea el topic si no existe (Admin Client)
func ensureTopic() {
	// kafka.Dial es como abrir un socket TCP crudo
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatal("Error conectando a Kafka para verificar topic:", err)
	}
	defer conn.Close()

	// El controlador es el broker "jefe" encargado de crear topics
	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("Error obteniendo controlador:", err)
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		log.Fatal("Error conectando al controlador:", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1, // Cuantas colas paralelas dividir este topic (paralelismo)
			ReplicationFactor: 1, // Cuantas copias de seguridad (redundancia)
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		// Kafka devuelve error si ya existe, lo cual está bien para nosotros
		fmt.Printf("Nota: Topic '%s' ya existe o acaba de ser creado.\n", topic)
	}
}
