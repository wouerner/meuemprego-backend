package handler

import (
	"net/http"

	"github.com/wouerner/runter-backend/internal/middleware"
	"github.com/wouerner/runter-backend/internal/service"
)

type UserHandler struct {
	authService service.AuthService
}

func NewUserHandler(authService service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// GetMe retorna os dados do usuário autenticado
// @Summary      Obter perfil do usuário autenticado
// @Description  Retorna as informações do usuário atual extraídas a partir do token JWT enviado no cabeçalho Authorization.
// @Tags         Usuários
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  domain.UserResponseDTO
// @Failure      401  {object}  map[string]string "Não autorizado"
// @Failure      404  {object}  map[string]string "Usuário não encontrado"
// @Router       /users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		respondJSON(w, http.StatusUnauthorized, map[string]string{"error": "Usuário não autenticado"})
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "Usuário não encontrado"})
		return
	}

	respondJSON(w, http.StatusOK, user)
}
