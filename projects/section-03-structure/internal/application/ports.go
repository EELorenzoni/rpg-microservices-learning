package application

import "github.com/EELorenzoni/rpg-microservices-learning/section-03/internal/domain"

// HeroRepository define el contrato para guardar héroes.
//
// EN NODE.JS/TS: Esto sería 'interface HeroRepository { save(hero: Hero): Promise<void>; }'
//
// CONCEPTO CLAVE: Inversión de Dependencia.
// La capa de Aplicación dice "Necesito que alguien guarde esto", pero no sabe CÓMO.
// No importamos Kafka ni Postgres aquí. Solo definimos qué necesitamos.
type HeroRepository interface {
	// Save guarda un héroe.
	// Recibe un puntero porque podría modificarlo (ej: agregar ID de base de datos),
	// aunque en este caso solo lo leemos.
	Save(hero *domain.Hero) error
}
