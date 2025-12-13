package application

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-03/internal/domain"
)

// CreateHeroCommand: DTO (Data Transfer Object).
// Son los datos que recibimos de "afuera" (Input).
// En Node.js: 'const { id, name } = req.body;'
type CreateHeroCommand struct {
	ID    string
	Name  string
	Power int
}

// CreateHeroService es nuestro "Caso de Uso" o "Handler".
// Es una struct que contiene sus dependencias.
//
// EN NODE.JS (NestJS):
// @Injectable()
//
//	class CreateHeroService {
//	  constructor(private repo: HeroRepository) {}
//	}
type CreateHeroService struct {
	repo HeroRepository
}

// NewCreateHeroService es el constructor.
// Recibe la INTERFAZ, no una implementación concreta.
// Esto nos permite inyectar un Mock para tests o cambiar Kafka por MySQL sin tocar este código.
func NewCreateHeroService(repo HeroRepository) *CreateHeroService {
	return &CreateHeroService{
		repo: repo,
	}
}

// Run ejecuta la lógica de negocio.
func (s *CreateHeroService) Run(cmd CreateHeroCommand) error {
	fmt.Printf("➡️  APP: Ejecutando caso de uso CreateHero para %s\n", cmd.Name)

	// 1. Llamar al Dominio (Factory)
	hero, err := domain.NewHero(cmd.ID, cmd.Name, cmd.Power)
	if err != nil {
		// Retornamos el error tal cual.
		// En un sistema real, podríamos envolverlo en algo más descriptivo.
		return fmt.Errorf("error creando entidad hero: %w", err)
	}

	// 2. Usar el Puerto (Repository) para guardar
	// No sabemos si esto va a Kafka, a disco o a la nube. No nos importa.
	if err := s.repo.Save(hero); err != nil {
		return fmt.Errorf("error guardando hero en repositorio: %w", err)
	}

	fmt.Printf("✅ APP: Hero %s creado y disparado evento de guardado!\n", hero.Name)
	return nil
}
