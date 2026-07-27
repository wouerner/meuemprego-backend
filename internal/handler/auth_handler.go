package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register cria um novo usuário no sistema
// @Summary      Registrar novo usuário
// @Description  Cadastra um novo usuário com nome, email e senha e retorna um token JWT de acesso.
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param        request body domain.RegisterDTO true "Dados de registro do usuário"
// @Success      201  {object}  domain.AuthResponseDTO
// @Failure      400  {object}  map[string]string "Dados inválidos"
// @Failure      409  {object}  map[string]string "Email já cadastrado"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto domain.RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	res, err := h.authService.Register(dto)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Campos obrigatórios ausentes"})
			return
		}
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			respondJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}

	respondJSON(w, http.StatusCreated, res)
}

// Login autentica um usuário existente
// @Summary      Autenticar usuário
// @Description  Valida o email e senha do usuário e retorna o token JWT de autorização.
// @Tags         Autenticação
// @Accept       json
// @Produce      json
// @Param        request body domain.LoginDTO true "Credenciais de acesso"
// @Success      200  {object}  domain.AuthResponseDTO
// @Failure      400  {object}  map[string]string "Requisição inválida"
// @Failure      401  {object}  map[string]string "Credenciais inválidas"
// @Failure      500  {object}  map[string]string "Erro interno do servidor"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto domain.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido fornecido"})
		return
	}

	res, err := h.authService.Login(dto)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			respondJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Erro interno do servidor"})
		return
	}

	respondJSON(w, http.StatusOK, res)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
