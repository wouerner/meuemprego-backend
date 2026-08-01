package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/service"
)

type MetricHandler struct {
	metricService service.MetricService
}

func NewMetricHandler(metricService service.MetricService) *MetricHandler {
	return &MetricHandler{metricService: metricService}
}

// List retorna todos os eventos de transbordo registrados
// @Summary      Listar eventos de métricas
// @Description  Retorna a lista completa de eventos de transbordo (cliques em WhatsApp/LinkedIn).
// @Tags         Métricas
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.MetricEventResponseDTO
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /metrics [get]
func (h *MetricHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := h.metricService.List()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, events)
}

// Track registra um novo evento de transbordo
// @Summary      Registrar evento de transbordo
// @Description  Registra um clique em WhatsApp ou LinkedIn para as métricas da plataforma.
// @Tags         Métricas
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.TrackMetricDTO true "Dados do evento"
// @Success      201  {object}  domain.MetricEventResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /metrics [post]
func (h *MetricHandler) Track(w http.ResponseWriter, r *http.Request) {
	var dto domain.TrackMetricDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	event, err := h.metricService.Track(dto)
	if err != nil {
		if errors.Is(err, service.ErrInvalidMetric) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusCreated, event)
}
