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

func (repo *Memory) Get(id string) (*domain.Hero, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	hero, exists := repo.data[id]
	if !exists {
		return nil, fmt.Errorf("hero not found with id %s", id)
	}
	return hero, nil
}

// Update actualiza un héroe existente.
func (repo *Memory) Update(hero *domain.Hero) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.data[hero.ID]; !exists {
		return fmt.Errorf("hero not found with id %s", hero.ID)
	}

	repo.data[hero.ID] = hero
	fmt.Printf("🔄 INFRA (DB): Actualizando Hero %s\n", hero.Name)
	return nil
}

// Delete elimina un héroe por ID.
func (repo *Memory) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.data[id]; !exists {
		return fmt.Errorf("hero not found with id %s", id)
	}

	delete(repo.data, id)
	fmt.Printf("🗑️ INFRA (DB): Eliminando Hero %s. Total records: %d\n", id, len(repo.data))
	return nil
}

// List retorna todos los héroes.
func (repo *Memory) List() ([]*domain.Hero, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	heroes := make([]*domain.Hero, 0, len(repo.data))
	for _, hero := range repo.data {
		heroes = append(heroes, hero)
	}

	fmt.Printf("📋 INFRA (DB): Listando %d héroes\n", len(heroes))
	return heroes, nil
}
