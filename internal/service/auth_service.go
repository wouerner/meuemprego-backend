package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wouerner/runter-backend/internal/config"
	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrInvalidInput       = errors.New("dados de entrada inválidos")
)

type AuthService interface {
	Register(dto domain.RegisterDTO) (*domain.AuthResponseDTO, error)
	Login(dto domain.LoginDTO) (*domain.AuthResponseDTO, error)
	GetUserByID(id uint) (*domain.UserResponseDTO, error)
}

type authService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *authService) Register(dto domain.RegisterDTO) (*domain.AuthResponseDTO, error) {
	if dto.Email == "" || dto.Password == "" || dto.Name == "" {
		return nil, ErrInvalidInput
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponseDTO{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) Login(dto domain.LoginDTO) (*domain.AuthResponseDTO, error) {
	if dto.Email == "" || dto.Password == "" {
		return nil, ErrInvalidInput
	}

	user, err := s.repo.FindByEmail(dto.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponseDTO{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) GetUserByID(id uint) (*domain.UserResponseDTO, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	res := user.ToResponse()
	return &res, nil
}

func (s *authService) generateJWT(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.cfg.JWTExpirationHours)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
