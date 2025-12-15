package herosrv

import (
	"fmt"

	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/core/domain"
	"github.com/EELorenzoni/rpg-microservices-learning/section-04/internal/core/ports"
)

// CreateHeroCommand: DTO (Data Transfer Object).
// Son los datos que recibimos de "afuera" (Input).
// En Node.js: 'const { id, name } = req.body;'
type CreateHeroCommand struct {
	ID    string
	Name  string
	Power int
}

// Service es nuestro "Caso de Uso" o "Handler".
// Es una struct que contiene sus dependencias.
//
// EN NODE.JS (NestJS):
// @Injectable()
//
//	class CreateHeroService {
//	  constructor(private repo: HeroRepository) {}
//	}
type Service struct {
	repo ports.HeroRepository
}

// New crea una instancia del servicio.
// Retorna *Service (un puntero).
//
// 💡 POINTERS (Sintaxis):
// 1. `*Service`: Indico que devuelvo una dirección.
// 2. `&Service{...}`: Creo el objeto y tomo su dirección.
func New(repo ports.HeroRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// Run ejecuta la lógica de negocio.
// Receiver (s *Service):
func (s *Service) Run(cmd CreateHeroCommand) error {
	fmt.Printf("➡️  APP: Ejecutando caso de uso CreateHero para %s\n", cmd.Name)

	// 1. Llamar al Dominio (Factory)
	// 💡 POINTERS: NewHero devuelve *Hero. 'hero' aquí es un puntero (una dirección de memoria: 0x1234abcd).
	hero, err := domain.NewHero(cmd.ID, cmd.Name)
	if err != nil {
		return fmt.Errorf("error creando hero: %w", err)
	}

	// 2. Llamar al Repositorio
	// 💡 POINTERS: repo.Save espera un *Hero. Como 'hero' YA es un *Hero, se lo pasamos directo.
	// Si 'hero' fuera valor, tendríamos que usar `repo.Save(&hero)`.
	if err := s.repo.Save(hero); err != nil {
		return fmt.Errorf("error guardando hero en repositorio: %w", err)
	}

	fmt.Printf("✅ APP: Hero %s creado y disparado evento de guardado!\n", hero.Name)
	return nil
}
