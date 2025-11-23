package handler

import (
	"net/http"

	"avito-backend-trainee-assignment-autumn-2025/pkg/logger"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check handles GET /health
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Health check requested")
	w.WriteHeader(http.StatusOK)
}
