package herosrv

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"
)

// CreateHeroCommand: DTO (Data Transfer Object).
type CreateHeroCommand struct {
	ID    string
	Name  string
	Power int
}

// Create ejecuta la lógica de creación de un héroe.
// Renombrado de Run a Create para mayor claridad.
func (s *Service) Create(cmd CreateHeroCommand) error {
	fmt.Printf("➡️  CORE (Service): Orquestando creación de %s\n", cmd.Name)

	// 1. Llamar al Dominio (Factory)
	hero, err := domain.NewHero(cmd.ID, cmd.Name)
	if err != nil {
		return fmt.Errorf("error creando hero: %w", err)
	}

	// 2. Persistencia (Base de Datos)
	if err := s.repo.Save(hero); err != nil {
		return fmt.Errorf("error guardando en DB: %w", err)
	}

	// 3. Notificación (Event Bus)
	if err := s.eventBus.Publish(hero, "HeroCreated"); err != nil {
		fmt.Printf("⚠️ WARN: Héroe guardado, pero falló el evento: %v\n", err)
	}

	fmt.Printf("✅ CORE: Hero %s procesado y guardado.\n", hero.Name)
	return nil
}
