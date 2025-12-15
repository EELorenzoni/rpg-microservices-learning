package domain

import (
	"errors"
	"fmt"
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

// NewHero es un "Factory" que crea un héroe válido.
// Retorna *Hero (un PUNTERO a Hero).
//
// 💡 POINTERS (Sintaxis):
// 1. TIPO '*Hero': El asterisco en la firma `(*Hero, error)` define que devolvemos una DIRECCIÓN de memoria, no el valor bruto.
// 2. OPERADOR '&': `return &Hero{...}` usa el "ampersand" para obtener la dirección de memoria del struct recién creado.
//   - Sin '&': Crearíamos el struct en el stack y devolveríamos una COPIA.
//   - Con '&': Go mueve el struct al "Heap" (memoria compartida) y nos da su ID (dirección 0x...).
//
// 💡 WHY POINTERS?
// 1. EFICIENCIA: Evitamos copiar structs grandes.
// 2. IDENTIDAD: Referenciamos al MISMO objeto único.
func NewHero(id string, name string) (*Hero, error) {
	if name == "" {
		return nil, fmt.Errorf("hero name cannot be empty")
	}

	// &Hero{...} <- "Genera el struct y dame su dirección (&)"
	return &Hero{
		ID:        id,
		Name:      name,
		Level:     1, // Default value
		Power:     10,
		CreatedAt: time.Now(),
	}, nil
}

// LevelUp aumenta el nivel y poder del héroe.
// Usa un "Pointer Receiver" (h *Hero).
//
// 💡 POINTERS (Sintaxis):
// - `(h *Hero)`: Aquí el '*' indica que 'h' NO es el objeto, es la LLAVE (dirección) para acceder al objeto.
// - Acceso transparente: Go nos permite hacer `h.Level` sin escribir `(*h).Level`. Es azúcar sintáctico.
//
// 💡 WHY POINTERS?
// MUTABILIDAD: Al tener la dirección, podemos modificar el valor real en memoria.
func (h *Hero) LevelUp() {
	h.Level++
	h.Power += 10
}
