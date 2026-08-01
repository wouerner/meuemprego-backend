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

type HunterHandler struct {
	hunterService service.HunterService
}

func NewHunterHandler(hunterService service.HunterService) *HunterHandler {
	return &HunterHandler{hunterService: hunterService}
}

// List retorna todos os job hunters cadastrados
// @Summary      Listar job hunters
// @Description  Retorna a lista completa de job hunters cadastrados na plataforma.
// @Tags         Hunters
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.HunterResponseDTO
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /hunters [get]
func (h *HunterHandler) List(w http.ResponseWriter, r *http.Request) {
	hunters, err := h.hunterService.List()
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, hunters)
}

// GetMe retorna o perfil de job hunter do usuário autenticado
// @Summary      Obter perfil de job hunter do usuário autenticado
// @Description  Retorna o perfil de job hunter vinculado ao usuário logado.
// @Tags         Hunters
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.HunterResponseDTO
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      404  {object}  map[string]string "Perfil não encontrado"
// @Router       /hunters/me [get]
func (h *HunterHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	profile, err := h.hunterService.GetProfileByUserID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrHunterNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Perfil de job hunter não encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

// SaveMe cria ou atualiza o perfil de job hunter do usuário autenticado
// @Summary      Salvar perfil de job hunter do usuário autenticado
// @Description  Cria (upsert) o perfil de job hunter vinculado ao usuário logado. Novos perfis nascem com status Pendente.
// @Tags         Hunters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body domain.SaveHunterDTO true "Dados do perfil de job hunter"
// @Success      200  {object}  domain.HunterResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      409  {object}  map[string]string "CPF já cadastrado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /hunters/me [put]
func (h *HunterHandler) SaveMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	var dto domain.SaveHunterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	profile, err := h.hunterService.SaveProfile(userID, dto)
	if err != nil {
		if errors.Is(err, service.ErrInvalidHunterInput) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Campos obrigatórios ausentes"})
			return
		}
		if errors.Is(err, repository.ErrHunterAlreadyExists) {
			respondJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

// SetStatus altera o status de aprovação de um job hunter (moderação de administrador)
// @Summary      Aprovar/rejeitar job hunter
// @Description  Altera o status (Pendente/Aprovado/Rejeitado) de um job hunter pelo ID.
// @Tags         Hunters
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path int true "ID do job hunter"
// @Param        request  body map[string]string true "Novo status (campo status)"
// @Success      200  {object}  domain.HunterResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      404  {object}  map[string]string "Hunter não encontrado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /hunters/{id}/status [patch]
func (h *HunterHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	profile, err := h.hunterService.SetStatus(uint(id), body.Status)
	if err != nil {
		if errors.Is(err, service.ErrInvalidHunterStatus) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, repository.ErrHunterNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Job hunter não encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, profile)
}

// IncrementContacts incrementa o contador de contatos de um job hunter
// @Summary      Incrementar contatos de job hunter
// @Description  Incrementa o contador total de contatos realizados com um job hunter.
// @Tags         Hunters
// @Produce      json
// @Security     BearerAuth
// @Param        id  path int true "ID do job hunter"
// @Success      200  {object}  domain.HunterResponseDTO
// @Failure      404  {object}  map[string]string "Hunter não encontrado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /hunters/{id}/contacts [post]
func (h *HunterHandler) IncrementContacts(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "ID inválido"})
		return
	}

	profile, err := h.hunterService.IncrementContacts(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrHunterNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "Job hunter não encontrado"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}
	respondJSON(w, http.StatusOK, profile)
}
