package herosrv

import "fmt"

// Delete elimina un héroe por ID.
func (s *Service) Delete(id string) error {
	fmt.Printf("➡️  CORE (Service): Eliminando héroe %s\n", id)

	// 1. Verificar que existe
	hero, err := s.repo.Get(id)
	if err != nil {
		// Publicar evento de fallo
		s.eventBus.Publish(hero, "HeroDeleteFailed")
		return fmt.Errorf("error obteniendo hero: %w", err)
	}

	// 2. Eliminar de DB
	if err := s.repo.Delete(id); err != nil {
		// Publicar evento de fallo
		s.eventBus.Publish(hero, "HeroDeleteFailed")
		return fmt.Errorf("error eliminando en DB: %w", err)
	}

	fmt.Printf("✅ CORE: Hero %s eliminado de la base de datos.\n", id)

	// 3. Publicar evento de éxito
	if err := s.eventBus.Publish(hero, "HeroDeleted"); err != nil {
		fmt.Printf("⚠️ WARN: Hero eliminado en DB, pero falló publicación del evento: %v\n", err)
	} else {
		fmt.Printf("📨 CORE: Evento 'HeroDeleted' publicado correctamente.\n")
	}

	return nil
}
