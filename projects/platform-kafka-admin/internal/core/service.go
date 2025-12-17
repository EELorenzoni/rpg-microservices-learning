package core

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

// AdminService manages Kafka topics.
type AdminService struct {
	brokerAddress string
}

func NewAdminService(brokerAddress string) *AdminService {
	return &AdminService{
		brokerAddress: brokerAddress,
	}
}

// CreateTopic creates a new topic with optional configuration.
// 🎓 KAFKA EXPERT:
// Configuration (Config entries) allows us to define:
// - Retention (How long to keep messages?)
// - Cleanup Policy (Delete old logs or Compact them?)
func (s *AdminService) CreateTopic(name string, partitions, replicas int, config map[string]string) error {
	conn, err := kafka.Dial("tcp", s.brokerAddress)
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	// Configuración del Topic
	topicConfig := kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     partitions,
		ReplicationFactor: replicas,
		ConfigEntries:     []kafka.ConfigEntry{},
	}

	// 🎓 Map config map to Kafka ConfigEntry
	for k, v := range config {
		topicConfig.ConfigEntries = append(topicConfig.ConfigEntries, kafka.ConfigEntry{
			ConfigName:  k,
			ConfigValue: v,
		})
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}

// ListTopics returns all topic names.
// 🎓 KAFKA EXPERT:
// Leemos la metadata del cluster (`ReadPartitions`).
// Kafka almacena qué topics existen y dónde están sus particiones.
func (s *AdminService) ListTopics() ([]string, error) {
	conn, err := kafka.Dial("tcp", s.brokerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, fmt.Errorf("failed to read partitions: %w", err)
	}

	// Deduplicate topics
	topicMap := make(map[string]bool)
	for _, p := range partitions {
		topicMap[p.Topic] = true
	}

	var topics []string
	for t := range topicMap {
		topics = append(topics, t)
	}
	return topics, nil
}

// DeleteTopic deletes a topic.
func (s *AdminService) DeleteTopic(name string) error {
	conn, err := kafka.Dial("tcp", s.brokerAddress)
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	err = conn.DeleteTopics(name)
	if err != nil {
		return fmt.Errorf("failed to delete topic: %w", err)
	}
	return nil
}
