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

// CreateTopic creates a new topic.
func (s *AdminService) CreateTopic(name string, partitions, replicas int) error {
	conn, err := kafka.Dial("tcp", s.brokerAddress)
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             name,
		NumPartitions:     partitions,
		ReplicationFactor: replicas,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	return nil
}

// ListTopics returns all topic names.
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
