package herohdl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/EELorenzoni/rpg-microservices-learning/section-05/internal/core/services/herosrv"
)

// HTTPHandler es un "Driving Adapter" para Web API.
type HTTPHandler struct {
	service *herosrv.Service
}

// NewHTTPHandler crea el handler.
func NewHTTPHandler(service *herosrv.Service) *HTTPHandler {
	return &HTTPHandler{
		service: service,
	}
}

// CreateHero maneja POST /heroes.
func (h *HTTPHandler) CreateHero(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parsear Input (JSON)
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("\n🌐 HANDLER (HTTP): Recibido POST /heroes -> %v\n", req)

	// 2. Map Input -> Command
	cmd := herosrv.CreateHeroCommand{
		Name:  req.Name,
		Power: 90,
	}

	// 3. Llamar Servicio
	start := time.Now()
	hero, err := h.service.Create(cmd)

	// 4. Mapear Output -> HTTP Response
	if err != nil {
		fmt.Printf("❌ HANDLER: Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "created",
		"time":   time.Since(start).String(),
		"hero":   hero,
	})
}

// GetHero maneja GET /heroes?id=...
func (h *HTTPHandler) GetHero(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	fmt.Printf("\n🌐 HANDLER (HTTP): Recibido GET /heroes?id=%s\n", id)

	// Renombrado: GetHero -> Get
	hero, err := h.service.Get(id)
	if err != nil {
		fmt.Printf("❌ HANDLER: Not Found: %v\n", err)
		http.Error(w, "Hero not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hero)
}

// UpdateHero maneja PUT /heroes?id=...
func (h *HTTPHandler) UpdateHero(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("\n🌐 HANDLER (HTTP): Recibido PUT /heroes?id=%s\n", id)

	cmd := herosrv.UpdateHeroCommand{
		ID:   id,
		Name: req.Name,
	}

	if err := h.service.Update(cmd); err != nil {
		fmt.Printf("❌ HANDLER: Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "updated",
	})
}

// DeleteHero maneja DELETE /heroes?id=...
func (h *HTTPHandler) DeleteHero(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		return
	}

	fmt.Printf("\n🌐 HANDLER (HTTP): Recibido DELETE /heroes?id=%s\n", id)

	if err := h.service.Delete(id); err != nil {
		fmt.Printf("❌ HANDLER: Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "deleted",
	})
}

// ListHeroes maneja GET /heroes (sin query params)
func (h *HTTPHandler) ListHeroes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("\n🌐 HANDLER (HTTP): Recibido GET /heroes (list)")

	heroes, err := h.service.List()
	if err != nil {
		fmt.Printf("❌ HANDLER: Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(heroes)
}
