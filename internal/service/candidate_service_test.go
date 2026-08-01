package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
)

type MockCandidateRepository struct {
	mock.Mock
}

func (m *MockCandidateRepository) Create(candidate *domain.Candidate) error {
	args := m.Called(candidate)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	candidate.ID = 1
	return nil
}

func (m *MockCandidateRepository) FindByID(id uint) (*domain.Candidate, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Candidate), args.Error(1)
}

func (m *MockCandidateRepository) FindByUserID(userID uint) (*domain.Candidate, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Candidate), args.Error(1)
}

func (m *MockCandidateRepository) FindByCPF(cpf string) (*domain.Candidate, error) {
	args := m.Called(cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Candidate), args.Error(1)
}

func (m *MockCandidateRepository) FindAll() ([]domain.Candidate, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Candidate), args.Error(1)
}

func (m *MockCandidateRepository) Update(candidate *domain.Candidate) error {
	args := m.Called(candidate)
	return args.Error(0)
}

func validCandidateDTO() domain.SaveCandidateDTO {
	return domain.SaveCandidateDTO{
		Name:                 "Carlos Eduardo",
		CPF:                  "529.982.247-25",
		Email:                "carlos.eduardo@email.com",
		Avatar:               "https://example.com/avatar.jpg",
		Headline:             "Engenheiro de Software Senior",
		Seniority:            "Senior",
		Area:                 "Tecnologia da Informação",
		CareerGoal:           "Transição para Tech Lead",
		ProfessionalMoment:   "Aberto a Propostas",
		RequestHunterContact: true,
		LGPDConsent:          true,
		LinkedInURL:          "https://www.linkedin.com/in/carlos-eduardo",
		WhatsAppNumber:       "5511998765432",
	}
}

func TestCandidateService_SaveProfile_CreatesWhenNoneExists(t *testing.T) {
	mockRepo := new(MockCandidateRepository)
	svc := service.NewCandidateService(mockRepo)

	mockRepo.On("FindByUserID", uint(1)).Return(nil, repository.ErrCandidateNotFound)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Candidate")).Return(nil)

	res, err := svc.SaveProfile(1, validCandidateDTO())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Carlos Eduardo", res.Name)
	assert.True(t, res.IsApproved)
	assert.NotNil(t, res.UserID)
	mockRepo.AssertExpectations(t)
}

func TestCandidateService_SaveProfile_UpdatesExisting(t *testing.T) {
	mockRepo := new(MockCandidateRepository)
	svc := service.NewCandidateService(mockRepo)

	existing := &domain.Candidate{ID: 1, UserID: uintPointer(1), CPF: "529.982.247-25", Name: "Antigo"}
	mockRepo.On("FindByUserID", uint(1)).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Candidate")).Return(nil)

	res, err := svc.SaveProfile(1, validCandidateDTO())

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "Carlos Eduardo", res.Name)
	mockRepo.AssertExpectations(t)
}

func TestCandidateService_SaveProfile_RejectsCPFConflict(t *testing.T) {
	mockRepo := new(MockCandidateRepository)
	svc := service.NewCandidateService(mockRepo)

	existing := &domain.Candidate{ID: 1, UserID: uintPointer(1), CPF: "111.222.333-44"}
	other := &domain.Candidate{ID: 2, CPF: "529.982.247-25"}

	mockRepo.On("FindByUserID", uint(1)).Return(existing, nil)
	mockRepo.On("FindByCPF", "529.982.247-25").Return(other, nil)

	res, err := svc.SaveProfile(1, validCandidateDTO())

	assert.ErrorIs(t, err, repository.ErrCandidateAlreadyExists)
	assert.Nil(t, res)
	mockRepo.AssertExpectations(t)
}

func TestCandidateService_SaveProfile_InvalidInput(t *testing.T) {
	mockRepo := new(MockCandidateRepository)
	svc := service.NewCandidateService(mockRepo)

	dto := validCandidateDTO()
	dto.Name = ""

	res, err := svc.SaveProfile(1, dto)

	assert.ErrorIs(t, err, service.ErrInvalidCandidateInput)
	assert.Nil(t, res)
	mockRepo.AssertNotCalled(t, "FindByUserID", mock.Anything)
}

func TestCandidateService_GetProfileByUserID(t *testing.T) {
	mockRepo := new(MockCandidateRepository)
	svc := service.NewCandidateService(mockRepo)

	candidate := &domain.Candidate{ID: 1, UserID: uintPointer(1), Name: "Carlos Eduardo"}
	mockRepo.On("FindByUserID", uint(1)).Return(candidate, nil)

	res, err := svc.GetProfileByUserID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), res.ID)
	mockRepo.AssertExpectations(t)
}

func TestCandidateService_List(t *testing.T) {
	mockRepo := new(MockCandidateRepository)
	svc := service.NewCandidateService(mockRepo)

	candidates := []domain.Candidate{
		{ID: 1, Name: "Carlos"},
		{ID: 2, Name: "Mariana"},
	}
	mockRepo.On("FindAll").Return(candidates, nil)

	res, err := svc.List()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Carlos", res[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestCandidateService_SetApproval(t *testing.T) {
	mockRepo := new(MockCandidateRepository)
	svc := service.NewCandidateService(mockRepo)

	candidate := &domain.Candidate{ID: 1, Name: "Carlos", IsApproved: true}
	mockRepo.On("FindByID", uint(1)).Return(candidate, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Candidate")).Return(nil)

	res, err := svc.SetApproval(1, false)

	assert.NoError(t, err)
	assert.False(t, res.IsApproved)
	mockRepo.AssertExpectations(t)
}

func uintPointer(v uint) *uint {
	return &v
}
