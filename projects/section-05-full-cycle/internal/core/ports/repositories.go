package ports

import "github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/domain"

// HeroRepository define las operaciones de persistencia (Guardar, Leer).
// Es un Puerto "Driven" (Salida).
type HeroRepository interface {
	// Save guarda un héroe.
	// Recibe un puntero porque podría modificarlo (ej: agregar ID de base de datos),
	// aunque en este caso solo lo leemos.
	// Save guarda un héroe.
	// Recibe un puntero porque podría modificarlo (ej: agregar ID de base de datos),
	// aunque en este caso solo lo leemos.
	Save(hero *domain.Hero) error

	// Get busca un héroe por ID.
	Get(id string) (*domain.Hero, error)
}
