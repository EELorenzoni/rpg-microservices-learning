package herosrv

import (
	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/ports"
)

// Service es nuestro "Caso de Uso" o "Handler".
// Es una struct que contiene sus dependencias.
//
// 💡 SOLID (SRP - Single Responsibility Principle):
// Este servicio tiene UNA responsabilidad: "Orquestar la creación de un héroe".
// No sabe CÓMO se guarda (Repo) ni CÓMO se notifica (EventBus).
type Service struct {
	repo     ports.HeroRepository
	eventBus ports.EventBus
}

// New crea una instancia del servicio.
// 💡 SOLID (DIP - Dependency Inversion Principle):
// Dependemos de ABSTRACCIONES (Interfaces ports.HeroRepository, ports.EventBus),
// no de concreciones (structs Kafka o Memory).
func New(repo ports.HeroRepository, eventBus ports.EventBus) *Service {
	return &Service{
		repo:     repo,
		eventBus: eventBus,
	}
}
