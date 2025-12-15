package herosrv

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"
)

// Get recupera un héroe.
// Renombrado de GetHero a Get.
func (s *Service) Get(id string) (*domain.Hero, error) {
	fmt.Printf("➡️  CORE (Service): Buscando héroe %s\n", id)
	return s.repo.Get(id)
}
