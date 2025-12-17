package handlers

import (
	"net/http"

	"github.com/EELorenzoni/rpg-microservices-learning/platform-kafka-admin/internal/core"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	service *core.AdminService
}

func NewAdminHandler(service *core.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// CreateTopicRequest DTO
type CreateTopicRequest struct {
	Name       string            `json:"name" binding:"required"`
	Partitions int               `json:"partitions"`
	Replicas   int               `json:"replicas"`
	Config     map[string]string `json:"config"` // 🎓 Optional: retention.ms, etc.
}

func (h *AdminHandler) CreateTopic(c *gin.Context) {
	var req CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Defaults
	if req.Partitions <= 0 {
		req.Partitions = 1
	}
	if req.Replicas <= 0 {
		req.Replicas = 1
	}

	if err := h.service.CreateTopic(req.Name, req.Partitions, req.Replicas, req.Config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "topic created", "topic": req})
}

func (h *AdminHandler) ListTopics(c *gin.Context) {
	topics, err := h.service.ListTopics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"topics": topics})
}

func (h *AdminHandler) DeleteTopic(c *gin.Context) {
	name := c.Param("name")
	if err := h.service.DeleteTopic(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "topic deleted", "name": name})
}
