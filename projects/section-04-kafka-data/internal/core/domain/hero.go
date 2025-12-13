package domain

import (
	"errors"
	"time"
)

// Errores de dominio
// En Node.js harías: class InvalidHeroError extends Error {}
// En Go, los errores son variables simples.
var (
	ErrHeroNameEmpty = errors.New("hero name cannot be empty")
	ErrHeroPowerLow  = errors.New("hero power must be at least 1")
)

// Hero representa a nuestro protagonista.
//
// EN NODE.JS: Esto sería una 'class Hero { constructor(id, name, level) { ... } }'
// o una interfaz de TypeScript 'interface Hero { ... }'.
//
// TAGS (`json:"..."`): Son metadatos. Le dicen a Go: "Cuando conviertas esto a JSON,
// usa 'id' minúscula en lugar de 'ID' mayúscula".
type Hero struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	Power     int       `json:"power"`
	CreatedAt time.Time `json:"created_at"`
}

// NewHero es una "Factory Function".
// Go no tiene constructores (como 'new Hero()').
// Por convención, creamos funciones que empiezan con 'New...'.
func NewHero(id string, name string, power int) (*Hero, error) {
	// Validación: Reglas de negocio puras.
	if name == "" {
		return nil, ErrHeroNameEmpty
	}
	if power < 1 {
		return nil, ErrHeroPowerLow
	}

	// Retornamos un PUNTERO (&Hero).
	// ¿Por qué?
	// - Si retornamos 'Hero' (sin asterisco), Go hace una COPIA del objeto.
	// - Si retornamos '*Hero', retornamos una REFERENCIA a la memoria (como pasar objetos en JS).
	// Queremos que el Héroe sea único y mutable, así que usamos puntero.
	return &Hero{
		ID:        id,
		Name:      name,
		Level:     1, // Nivel inicial siempre 1
		Power:     power,
		CreatedAt: time.Now(),
	}, nil
}

// LevelUp sube de nivel al héroe.
//
// "(h *Hero)" es el RECEIVER.
// Significa: "Esta función es un método que pertenece a la estructura Hero".
// En Node.js: 'Hero.prototype.levelUp = function() { ... }'
//
// Usamos puntero (*Hero) porque queremos MODIFICAR al héroe (h.Level++).
// Si usáramos (h Hero), modificaríamos una copia y el original no cambiaría.
func (h *Hero) LevelUp() {
	h.Level++
	h.Power += 10 // Bonus por subir de nivel
}
