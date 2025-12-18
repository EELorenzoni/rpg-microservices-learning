package herosrv

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"
	"github.com/google/uuid"
)

// CreateHeroCommand: DTO (Data Transfer Object).
type CreateHeroCommand struct {
	Name  string
	Power int
}

// Create ejecuta la lógica de creación de un héroe.
// Renombrado de Run a Create para mayor claridad.
func (s *Service) Create(cmd CreateHeroCommand) (*domain.Hero, error) {
	fmt.Printf("➡️  CORE (Service): Orquestando creación de %s\n", cmd.Name)

	// 1. Generar ID único
	heroID := uuid.New().String()

	// 2. Llamar al Dominio (Factory)
	hero, err := domain.NewHero(heroID, cmd.Name)
	if err != nil {
		// Publicar evento de fallo
		s.eventBus.Publish(&domain.Hero{ID: heroID, Name: cmd.Name}, "HeroCreateFailed")
		return nil, fmt.Errorf("error creando hero: %w", err)
	}

	// 3. Persistencia (Base de Datos)
	if err := s.repo.Save(hero); err != nil {
		// Publicar evento de fallo
		s.eventBus.Publish(hero, "HeroCreateFailed")
		return nil, fmt.Errorf("error guardando en DB: %w", err)
	}

	fmt.Printf("✅ CORE: Hero %s guardado en la base de datos.\n", hero.Name)

	// 4. Publicar evento de éxito
	if err := s.eventBus.Publish(hero, "HeroCreated"); err != nil {
		fmt.Printf("⚠️ WARN: Hero guardado en DB, pero falló publicación del evento: %v\n", err)
	} else {
		fmt.Printf("📨 CORE: Evento 'HeroCreated' publicado correctamente.\n")
	}

	return hero, nil
}
