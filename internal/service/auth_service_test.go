package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wouerner/runter-backend/internal/config"
	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository simula a camada de banco de dados sem conectar no Postgres
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	user.ID = 1
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "secret", JWTExpirationHours: 24}
	authSvc := service.NewAuthService(mockRepo, cfg)

	dto := domain.RegisterDTO{
		Name:     "João Silva",
		Email:    "joao@exemplo.com",
		Password: "senha123_segura",
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	res, err := authSvc.Register(dto)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, "joao@exemplo.com", res.User.Email)
	assert.Equal(t, "João Silva", res.User.Name)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "secret", JWTExpirationHours: 24}
	authSvc := service.NewAuthService(mockRepo, cfg)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123_segura"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       1,
		Name:     "João Silva",
		Email:    "joao@exemplo.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("FindByEmail", "joao@exemplo.com").Return(existingUser, nil)

	loginDTO := domain.LoginDTO{
		Email:    "joao@exemplo.com",
		Password: "senha123_segura",
	}

	res, err := authSvc.Login(loginDTO)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
	assert.Equal(t, uint(1), res.User.ID)
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	cfg := &config.Config{JWTSecret: "secret", JWTExpirationHours: 24}
	authSvc := service.NewAuthService(mockRepo, cfg)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha_correta"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		ID:       1,
		Email:    "joao@exemplo.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("FindByEmail", "joao@exemplo.com").Return(existingUser, nil)

	loginDTO := domain.LoginDTO{
		Email:    "joao@exemplo.com",
		Password: "senha_errada",
	}

	res, err := authSvc.Login(loginDTO)

	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}
