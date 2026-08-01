package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/middleware"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
)

type CandidateHandler struct {
	candidateService service.CandidateService
}

func NewCandidateHandler(candidateService service.CandidateService) *CandidateHandler {
	return &CandidateHandler{candidateService: candidateService}
}

// List retorna todos os candidatos cadastrados
// @Summary      Listar candidatos
// @Description  Retorna a lista completa de candidatos cadastrados na plataforma.
// @Tags         Candidatos
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.CandidateResponseDTO
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /candidates [get]
func (h *CandidateHandler) List(w http.ResponseWriter, r *http.Request) {
	candidates, err := h.candidateService.List()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, candidates)
}

// GetMe retorna o perfil de candidato do usuário autenticado
// @Summary      Obter perfil de candidato do usuário autenticado
// @Description  Retorna o perfil de candidato vinculado ao usuário logado.
// @Tags         Candidatos
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.CandidateResponseDTO
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      404  {object}  map[string]string "Perfil não encontrado"
// @Router       /candidates/me [get]
func (h *CandidateHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	profile, err := h.candidateService.GetProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrCandidateNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Perfil de candidato não encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

// SaveMe cria ou atualiza o perfil de candidato do usuário autenticado
// @Summary      Salvar perfil de candidato do usuário autenticado
// @Description  Cria (upsert) o perfil de candidato vinculado ao usuário logado.
// @Tags         Candidatos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.SaveCandidateDTO true "Dados do perfil de candidato"
// @Success      200  {object}  domain.CandidateResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      409  {object}  map[string]string "CPF já cadastrado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /candidates/me [put]
func (h *CandidateHandler) SaveMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	var dto domain.SaveCandidateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	profile, err := h.candidateService.SaveProfile(userID, dto)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCandidateInput) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Campos obrigatórios ausentes"})
			return
		}
		if errors.Is(err, repository.ErrCandidateAlreadyExists) {
			respondJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

// SetApproval altera a aprovação de um candidato (moderação de administrador)
// @Summary      Aprovar/rejeitar candidato
// @Description  Altera o flag de aprovação de um candidato pelo ID.
// @Tags         Candidatos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path int true "ID do candidato"
// @Param        request  body map[string]bool true "Novo estado de aprovação (campo approved)"
// @Success      200  {object}  domain.CandidateResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      404  {object}  map[string]string "Candidato não encontrado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /candidates/{id} [patch]
func (h *CandidateHandler) SetApproval(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}

	var body struct {
		Approved bool `json:"approved"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	profile, err := h.candidateService.SetApproval(uint(id), body.Approved)
	if err != nil {
		if errors.Is(err, repository.ErrCandidateNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Candidato não encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

func userIDFromContext(r *http.Request) (uint, bool) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	return userID, ok
}
