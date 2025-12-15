package herorepo

import (
	"fmt"
	"sync"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"
)

// Memory es un adaptador "fake" para base de datos.
type Memory struct {
	mu   sync.RWMutex
	data map[string]*domain.Hero
}

// NewMemory crea el repositorio en memoria.
func NewMemory() *Memory {
	return &Memory{
		data: make(map[string]*domain.Hero),
	}
}

// Save simula el INSERT en DB.
func (repo *Memory) Save(hero *domain.Hero) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data[hero.ID] = hero
	fmt.Printf("💾 INFRA (DB): Guardando Hero %s en base de datos (Memoria)... Total records: %d\n", hero.Name, len(repo.data))
	return nil
}

// Get recupera un héroe por ID.
func (repo *Memory) Get(id string) (*domain.Hero, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	hero, exists := repo.data[id]
	if !exists {
		return nil, fmt.Errorf("hero not found with id %s", id)
	}
	return hero, nil
}
