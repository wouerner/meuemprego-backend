package handler

import (
	"net/http"
	"time"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check verifica o status de saúde da API
// @Summary      Healthcheck da API
// @Description  Retorna o status operacional da aplicação, timestamp UTC e versão da API.
// @Tags         Saúde
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"version":   "v1",
		"service":   "runter-backend-api",
	})
}
