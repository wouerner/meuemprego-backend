package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
)

type MockHunterRepository struct {
	mock.Mock
}

func (m *MockHunterRepository) Create(hunter *domain.Hunter) error {
	args := m.Called(hunter)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	hunter.ID = 1
	return nil
}

func (m *MockHunterRepository) FindByID(id uint) (*domain.Hunter, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hunter), args.Error(1)
}

func (m *MockHunterRepository) FindByUserID(userID uint) (*domain.Hunter, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hunter), args.Error(1)
}

func (m *MockHunterRepository) FindByCPF(cpf string) (*domain.Hunter, error) {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Hunter), args.Error(1)
}

func (m *MockHunterRepository) FindAll() ([]domain.Hunter, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Hunter), args.Error(1)
}

func (m *MockHunterRepository) Update(hunter *domain.Hunter) error {
	args := m.Called(hunter)
	return args.Error(0)
}

func validHunterDTO() domain.SaveHunterDTO {
	return domain.SaveHunterDTO{
		Name:              "Juliana Mendes",
		CPF:               "987.654.321-09",
		Email:             "juliana.mendes@career.com",
		Avatar:            "https://example.com/avatar.jpg",
		Headline:          "Executive Headhunter & Coach de Carreira Tech",
		Bio:               "Especialista em recolocação de executivos de TI.",
		Specialties:       []string{"Tecnologia da Informação", "Produto"},
		SenioritiesServed: []string{"Senior", "Especialista"},
		ServiceModel:      "Assessoria Completa",
		LinkedInURL:       "https://www.linkedin.com/in/juliana-mendes",
		WhatsAppNumber:    "5511988887777",
	}
}

func TestHunterService_SaveProfile_CreatesWithPendenteStatus(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	mockRepo.On("FindByUserID", uint(1)).Return(nil, repository.ErrHunterNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Hunter")).Return(nil)

	res, err := svc.SaveProfile(1, validHunterDTO())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Juliana Mendes", res.Name)
	assert.Equal(t, domain.HunterStatusPendente, res.Status)
	assert.Equal(t, 0, res.TotalContactsCount)
	assert.NotNil(t, res.UserID)
	assert.Equal(t, []string{"Tecnologia da Informação", "Produto"}, res.Specialties)
	mockRepo.AssertExpectations(t)
}

func TestHunterService_SaveProfile_UpdatesExisting(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	existing := &domain.Hunter{ID: 1, UserID: uintPointer(1), CPF: "987.654.321-09", Name: "Antiga", Status: domain.HunterStatusAprovado}
	mockRepo.On("FindByUserID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Hunter")).Return(nil)

	res, err := svc.SaveProfile(1, validHunterDTO())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Juliana Mendes", res.Name)
	assert.Equal(t, domain.HunterStatusAprovado, res.Status)
	mockRepo.AssertExpectations(t)
}

func TestHunterService_SaveProfile_RejectsCPFConflict(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	existing := &domain.Hunter{ID: 1, UserID: uintPointer(1), CPF: "111.222.333-44"}
	other := &domain.Hunter{ID: 2, CPF: "987.654.321-09"}

	mockRepo.On("FindByUserID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByCPF", "987.654.321-09").Return(other, nil)

	res, err := svc.SaveProfile(1, validHunterDTO())

	assert.ErrorIs(t, err, repository.ErrHunterAlreadyExists)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestHunterService_SaveProfile_InvalidInput(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	dto := validHunterDTO()
	dto.Name = ""

	res, err := svc.SaveProfile(1, dto)

	assert.ErrorIs(t, err, service.ErrInvalidHunterInput)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "FindByUserID", mock.Anything)
}

func TestHunterService_SetStatus_Valid(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	hunter := &domain.Hunter{ID: 1, Name: "Juliana", Status: domain.HunterStatusPendente}
	mockRepo.On("FindByID", uint(1)).Return(hunter, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Hunter")).Return(nil)

	res, err := svc.SetStatus(1, domain.HunterStatusAprovado)

	assert.NoError(t, err)
	assert.Equal(t, domain.HunterStatusAprovado, res.Status)
	mockRepo.AssertExpectations(t)
}

func TestHunterService_SetStatus_InvalidStatus(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	res, err := svc.SetStatus(1, "Desconhecido")

	assert.ErrorIs(t, err, service.ErrInvalidHunterStatus)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "FindByID", mock.Anything)
}

func TestHunterService_IncrementContacts(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	hunter := &domain.Hunter{ID: 1, TotalContactsCount: 41}
	mockRepo.On("FindByID", uint(1)).Return(hunter, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Hunter")).Return(nil)

	res, err := svc.IncrementContacts(1)

	assert.NoError(t, err)
	assert.Equal(t, 42, res.TotalContactsCount)
	mockRepo.AssertExpectations(t)
}

func TestHunterService_List(t *testing.T) {
	mockRepo := new(MockHunterRepository)
	svc := service.NewHunterService(mockRepo)

	hunters := []domain.Hunter{
		{ID: 1, Name: "Juliana", Status: domain.HunterStatusAprovado},
		{ID: 2, Name: "Roberto", Status: domain.HunterStatusPendente},
	}
	mockRepo.On("FindAll").Return(hunters, nil)

	res, err := svc.List()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Juliana", res[0].Name)
	mockRepo.AssertExpectations(t)
}
