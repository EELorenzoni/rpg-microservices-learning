package herosrv

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"
)

// List retorna todos los héroes.
func (s *Service) List() ([]*domain.Hero, error) {
	fmt.Println("➡️  CORE (Service): Listando todos los héroes")
	return s.repo.List()
}
