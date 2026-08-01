package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
)

type AccessRequestHandler struct {
	accessRequestService service.AccessRequestService
}

func NewAccessRequestHandler(accessRequestService service.AccessRequestService) *AccessRequestHandler {
	return &AccessRequestHandler{accessRequestService: accessRequestService}
}

// ListMe retorna as solicitações de acesso recebidas pelo candidato autenticado
// @Summary      Listar solicitações de acesso do candidato
// @Description  Retorna a caixa de entrada de solicitações de acesso de job hunters para o candidato autenticado.
// @Tags         Access Requests
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.AccessRequestResponseDTO
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      404  {object}  map[string]string "Perfil de candidato não encontrado"
// @Router       /access-requests/me [get]
func (h *AccessRequestHandler) ListMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	requests, err := h.accessRequestService.ListForCandidate(userID)
	if err != nil {
		if errors.Is(err, service.ErrCandidateProfileMissing) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, requests)
}

// Send cria uma solicitação de acesso de um job hunter aprovado para um candidato
// @Summary      Enviar solicitação de acesso
// @Description  Envia uma solicitação de acesso de um job hunter (autenticado e aprovado) para um candidato.
// @Tags         Access Requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.CreateAccessRequestDTO true "Dados da solicitação"
// @Success      201  {object}  domain.AccessRequestResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      403  {object}  map[string]string "Hunter não aprovado"
// @Failure      404  {object}  map[string]string "Perfil de hunter não encontrado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /access-requests [post]
func (h *AccessRequestHandler) Send(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	var dto domain.CreateAccessRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	res, err := h.accessRequestService.Send(userID, dto)
	if err != nil {
		if errors.Is(err, service.ErrInvalidAccessRequest) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrHunterProfileMissing) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrHunterNotApproved) {
			respondJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusCreated, res)
}

// Respond atualiza o status de uma solicitação de acesso (aceitar/rejeitar)
// @Summary      Responder solicitação de acesso
// @Description  Aceita ou rejeita uma solicitação de acesso recebida pelo candidato autenticado.
// @Tags         Access Requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path int true "ID da solicitação"
// @Param        request  body domain.RespondAccessRequestDTO true "Novo status (accepted/rejected)"
// @Success      200  {object}  domain.AccessRequestResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      404  {object}  map[string]string "Solicitação não encontrada"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /access-requests/{id} [patch]
func (h *AccessRequestHandler) Respond(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}

	var dto domain.RespondAccessRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	res, err := h.accessRequestService.Respond(userID, uint(id), dto.Status)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequestStatus) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, service.ErrCandidateProfileMissing) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, repository.ErrAccessRequestNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Solicitação de acesso não encontrada"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, res)
}
