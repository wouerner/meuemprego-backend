package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
)

type MockAccessRequestRepository struct {
	mock.Mock
}

func (m *MockAccessRequestRepository) Create(request *domain.AccessRequest) error {
	args := m.Called(request)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	request.ID = 1
	return nil
}

func (m *MockAccessRequestRepository) FindByID(id uint) (*domain.AccessRequest, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AccessRequest), args.Error(1)
}

func (m *MockAccessRequestRepository) FindByCandidateID(candidateID uint) ([]domain.AccessRequest, error) {
	args := m.Called(candidateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.AccessRequest), args.Error(1)
}

func (m *MockAccessRequestRepository) Update(request *domain.AccessRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func TestAccessRequestService_Send_RejectsWhenHunterNotApproved(t *testing.T) {
	requestRepo := new(MockAccessRequestRepository)
	hunterRepo := new(MockHunterRepository)
	candidateRepo := new(MockCandidateRepository)
	svc := service.NewAccessRequestService(requestRepo, hunterRepo, candidateRepo)

	hunter := &domain.Hunter{ID: 1, UserID: uintPointer(1), Status: domain.HunterStatusPendente}
	hunterRepo.On("FindByUserID", uint(1)).Return(hunter, nil)

	res, err := svc.Send(1, domain.CreateAccessRequestDTO{CandidateID: 1, Message: "Olá!"})

	assert.ErrorIs(t, err, service.ErrHunterNotApproved)
	assert.Nil(t, res)
	requestRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestAccessRequestService_Send_Success(t *testing.T) {
	requestRepo := new(MockAccessRequestRepository)
	hunterRepo := new(MockHunterRepository)
	candidateRepo := new(MockCandidateRepository)
	svc := service.NewAccessRequestService(requestRepo, hunterRepo, candidateRepo)

	hunter := &domain.Hunter{
		ID: 1, UserID: uintPointer(1), Status: domain.HunterStatusAprovado,
		Name: "Juliana", Avatar: "avatar.jpg", Headline: "Headhunter",
		Specialties: `["Tech"]`,
	}
	hunterRepo.On("FindByUserID", uint(1)).Return(hunter, nil)
	hunterRepo.On("FindByID", uint(1)).Return(hunter, nil)
	requestRepo.On("Create", mock.AnythingOfType("*domain.AccessRequest")).Return(nil)

	res, err := svc.Send(1, domain.CreateAccessRequestDTO{CandidateID: 5, Message: "Olá!"})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, uint(5), res.CandidateID)
	assert.Equal(t, domain.AccessRequestPending, res.Status)
	assert.Equal(t, "Juliana", res.HunterName)
	requestRepo.AssertExpectations(t)
}

func TestAccessRequestService_Send_InvalidInput(t *testing.T) {
	requestRepo := new(MockAccessRequestRepository)
	hunterRepo := new(MockHunterRepository)
	candidateRepo := new(MockCandidateRepository)
	svc := service.NewAccessRequestService(requestRepo, hunterRepo, candidateRepo)

	res, err := svc.Send(1, domain.CreateAccessRequestDTO{})

	assert.ErrorIs(t, err, service.ErrInvalidAccessRequest)
	assert.Nil(t, res)
	hunterRepo.AssertNotCalled(t, "FindByUserID", mock.Anything)
}

func TestAccessRequestService_ListForCandidate_NoProfile(t *testing.T) {
	requestRepo := new(MockAccessRequestRepository)
	hunterRepo := new(MockHunterRepository)
	candidateRepo := new(MockCandidateRepository)
	svc := service.NewAccessRequestService(requestRepo, hunterRepo, candidateRepo)

	candidateRepo.On("FindByUserID", uint(1)).Return(nil, repository.ErrCandidateNotFound)

	res, err := svc.ListForCandidate(1)

	assert.ErrorIs(t, err, service.ErrCandidateProfileMissing)
	assert.Nil(t, res)
}

func TestAccessRequestService_Respond_NotOwner(t *testing.T) {
	requestRepo := new(MockAccessRequestRepository)
	hunterRepo := new(MockHunterRepository)
	candidateRepo := new(MockCandidateRepository)
	svc := service.NewAccessRequestService(requestRepo, hunterRepo, candidateRepo)

	candidate := &domain.Candidate{ID: 2, UserID: uintPointer(1)}
	req := &domain.AccessRequest{ID: 9, CandidateID: 3, Status: domain.AccessRequestPending}

	candidateRepo.On("FindByUserID", uint(1)).Return(candidate, nil)
	requestRepo.On("FindByID", uint(9)).Return(req, nil)

	res, err := svc.Respond(1, 9, domain.AccessRequestAccepted)

	assert.ErrorIs(t, err, repository.ErrAccessRequestNotFound)
	assert.Nil(t, res)
	requestRepo.AssertNotCalled(t, "Update", mock.Anything)
}

func TestAccessRequestService_Respond_InvalidStatus(t *testing.T) {
	requestRepo := new(MockAccessRequestRepository)
	hunterRepo := new(MockHunterRepository)
	candidateRepo := new(MockCandidateRepository)
	svc := service.NewAccessRequestService(requestRepo, hunterRepo, candidateRepo)

	res, err := svc.Respond(1, 9, "lixo")

	assert.ErrorIs(t, err, service.ErrInvalidRequestStatus)
	assert.Nil(t, res)
	candidateRepo.AssertNotCalled(t, "FindByUserID", mock.Anything)
}

func TestAccessRequestService_Respond_Success(t *testing.T) {
	requestRepo := new(MockAccessRequestRepository)
	hunterRepo := new(MockHunterRepository)
	candidateRepo := new(MockCandidateRepository)
	svc := service.NewAccessRequestService(requestRepo, hunterRepo, candidateRepo)

	candidate := &domain.Candidate{ID: 2, UserID: uintPointer(1)}
	req := &domain.AccessRequest{ID: 9, HunterID: 1, CandidateID: 2, Status: domain.AccessRequestPending}
	hunter := &domain.Hunter{ID: 1, Name: "Juliana", Specialties: `["Tech"]`}

	candidateRepo.On("FindByUserID", uint(1)).Return(candidate, nil)
	requestRepo.On("FindByID", uint(9)).Return(req, nil)
	requestRepo.On("Update", mock.AnythingOfType("*domain.AccessRequest")).Return(nil)
	hunterRepo.On("FindByID", uint(1)).Return(hunter, nil)

	res, err := svc.Respond(1, 9, domain.AccessRequestAccepted)

	assert.NoError(t, err)
	assert.Equal(t, domain.AccessRequestAccepted, res.Status)
	requestRepo.AssertExpectations(t)
}

type MockMetricRepository struct {
	mock.Mock
}

func (m *MockMetricRepository) Create(event *domain.MetricEvent) error {
	args := m.Called(event)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	event.ID = 1
	return nil
}

func (m *MockMetricRepository) FindAll() ([]domain.MetricEvent, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.MetricEvent), args.Error(1)
}

func TestMetricService_Track_Success(t *testing.T) {
	metricRepo := new(MockMetricRepository)
	svc := service.NewMetricService(metricRepo)

	metricRepo.On("Create", mock.AnythingOfType("*domain.MetricEvent")).Return(nil)

	res, err := svc.Track(domain.TrackMetricDTO{
		TargetType: "hunter", TargetID: "1", TargetName: "Juliana", Channel: "whatsapp", InitiatedByRole: "candidato",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "hunter", res.TargetType)
	assert.Equal(t, "whatsapp", res.Channel)
	assert.False(t, res.Timestamp.IsZero())
	metricRepo.AssertExpectations(t)
}

func TestMetricService_Track_InvalidChannel(t *testing.T) {
	metricRepo := new(MockMetricRepository)
	svc := service.NewMetricService(metricRepo)

	res, err := svc.Track(domain.TrackMetricDTO{
		TargetType: "hunter", TargetID: "1", TargetName: "Juliana", Channel: "telegram",
	})

	assert.ErrorIs(t, err, service.ErrInvalidMetric)
	assert.Nil(t, res)
	metricRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestMetricService_List(t *testing.T) {
	metricRepo := new(MockMetricRepository)
	svc := service.NewMetricService(metricRepo)

	events := []domain.MetricEvent{
		{ID: 1, TargetType: "hunter", TargetName: "Juliana", Channel: "whatsapp"},
		{ID: 2, TargetType: "candidate", TargetName: "Carlos", Channel: "linkedin"},
	}
	metricRepo.On("FindAll").Return(events, nil)

	res, err := svc.List()

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, "Carlos", res[1].TargetName)
	metricRepo.AssertExpectations(t)
}
