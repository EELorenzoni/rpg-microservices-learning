package herosrv

import (
	"fmt"
)

// UpdateHeroCommand: DTO para actualización.
type UpdateHeroCommand struct {
	ID   string
	Name string
}

// Update actualiza un héroe existente.
func (s *Service) Update(cmd UpdateHeroCommand) error {
	fmt.Printf("➡️  CORE (Service): Actualizando héroe %s\n", cmd.ID)

	// 1. Obtener héroe existente
	hero, err := s.repo.Get(cmd.ID)
	if err != nil {
		// Publicar evento de fallo
		s.eventBus.Publish(hero, "HeroUpdateFailed")
		return fmt.Errorf("error obteniendo hero: %w", err)
	}

	// 2. Actualizar campos
	if cmd.Name != "" {
		hero.Name = cmd.Name
	}

	// 3. Persistir cambios
	if err := s.repo.Update(hero); err != nil {
		// Publicar evento de fallo
		s.eventBus.Publish(hero, "HeroUpdateFailed")
		return fmt.Errorf("error actualizando en DB: %w", err)
	}

	fmt.Printf("✅ CORE: Hero %s actualizado en la base de datos.\n", hero.Name)

	// 4. Publicar evento de éxito
	if err := s.eventBus.Publish(hero, "HeroUpdated"); err != nil {
		fmt.Printf("⚠️ WARN: Hero actualizado en DB, pero falló publicación del evento: %v\n", err)
	} else {
		fmt.Printf("📨 CORE: Evento 'HeroUpdated' publicado correctamente.\n")
	}

	return nil
}
