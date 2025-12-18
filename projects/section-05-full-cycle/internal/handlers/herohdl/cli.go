package herohdl

import (
	"fmt"
	"time"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/services/herosrv"
)

// CLIHandler es un "Driving Adapter" (Adaptador de Entrada).
// Su responsabilidad es recibir INPUTS (de terminal, HTTP, etc.)
// y llamar al SERVICIO.
type CLIHandler struct {
	service *herosrv.Service
}

// NewCLIHandler crea el handler inyectándole el servicio.
func NewCLIHandler(service *herosrv.Service) *CLIHandler {
	return &CLIHandler{
		service: service,
	}
}

// CreateHeroSimulated simula que un usuario tipea un comando en la terminal.
// Recibe "strings" crudos (simulando argv) y orquesta la llamada.
func (h *CLIHandler) CreateHeroSimulated(name string) {
	fmt.Printf("\n🎮 HANDLER (CLI): Recibido input usuario -> Name: %s\n", name)

	// 1. DTO/Command Mapping: Convertir input externo a Estructura de Dominio (Command)
	cmd := herosrv.CreateHeroCommand{
		Name:  name,
		Power: 90, // En un caso real, esto podría venir de un flag o default
	}

	// 2. Llamar al Servicio (Use Case)
	start := time.Now()
	hero, err := h.service.Create(cmd)

	// 3. Manejar Respuesta (Output)
	if err != nil {
		fmt.Printf("❌ HANDLER: Error: %v\n", err)
	} else {
		duration := time.Since(start)
		fmt.Printf("✅ HANDLER: Héroe %s creado con ID %s en %v.\n", hero.Name, hero.ID, duration)
	}
}
