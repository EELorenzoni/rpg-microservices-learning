package ports

import "github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"

// 💡 SOLID (ISP - Interface Segregation Principle):
// Separamos "Guardar Dato" (Repository) de "Publicar Evento" (EventBus).
// El repositorio no tiene por qué saber de Kafka. El bus de eventos sí.

// EventBus define el contrato para publicar eventos de dominio.
type EventBus interface {
	Publish(hero *domain.Hero, eventType string) error
}
