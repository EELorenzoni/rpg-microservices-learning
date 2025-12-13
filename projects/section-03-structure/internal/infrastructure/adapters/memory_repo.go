package adapters

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-03/internal/domain"
)

// InMemoryHeroRepository simula una base de datos en memoria.
//
// "Implementa" la interfaz application.HeroRepository implícitamente
// porque tiene el método Save con la firma correcta.
type InMemoryHeroRepository struct {
	// Aquí podríamos tener un map[string]*domain.Hero para guardar de verdad
}

// NewInMemoryHeroRepository crea una instancia.
func NewInMemoryHeroRepository() *InMemoryHeroRepository {
	return &InMemoryHeroRepository{}
}

// Save cumple con el contrato de HeroRepository.
func (r *InMemoryHeroRepository) Save(hero *domain.Hero) error {
	// Simulación de guardado (DB o Kafka)
	fmt.Printf("💾 INFRA (Memory): Guardando Hero %s (ID: %s) en 'Base de Datos'...\n", hero.Name, hero.ID)
	return nil
}
