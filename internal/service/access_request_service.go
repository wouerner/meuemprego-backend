package service

import (
	"errors"

	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
)

var (
	ErrHunterNotApproved       = errors.New("seu perfil precisa ser aprovado pelo administrador antes de enviar solicitações de acesso aos candidatos")
	ErrCandidateProfileMissing = errors.New("perfil de candidato não encontrado para o usuário autenticado")
	ErrHunterProfileMissing    = errors.New("perfil de job hunter não encontrado para o usuário autenticado")
	ErrInvalidAccessRequest    = errors.New("solicitação de acesso inválida")
	ErrInvalidRequestStatus    = errors.New("status inválido para a solicitação de acesso")
)

type AccessRequestService interface {
	Send(userID uint, dto domain.CreateAccessRequestDTO) (*domain.AccessRequestResponseDTO, error)
	ListForCandidate(userID uint) ([]domain.AccessRequestResponseDTO, error)
	Respond(userID uint, id uint, status string) (*domain.AccessRequestResponseDTO, error)
}

type accessRequestService struct {
	requestRepo    repository.AccessRequestRepository
	hunterRepo     repository.HunterRepository
	candidateRepo  repository.CandidateRepository
}

func NewAccessRequestService(
	requestRepo repository.AccessRequestRepository,
	hunterRepo repository.HunterRepository,
	candidateRepo repository.CandidateRepository,
) AccessRequestService {
	return &accessRequestService{requestRepo: requestRepo, hunterRepo: hunterRepo, candidateRepo: candidateRepo}
}

func isValidAccessRequestStatus(status string) bool {
	return status == domain.AccessRequestPending ||
		status == domain.AccessRequestAccepted ||
		status == domain.AccessRequestRejected
}

func (s *accessRequestService) toResponse(req *domain.AccessRequest) (*domain.AccessRequestResponseDTO, error) {
	hunter, err := s.hunterRepo.FindByID(req.HunterID)
	if err != nil {
		return nil, err
	}

	return &domain.AccessRequestResponseDTO{
		ID:                req.ID,
		HunterID:          req.HunterID,
		CandidateID:       req.CandidateID,
		HunterName:        hunter.Name,
		HunterAvatar:      hunter.Avatar,
		HunterHeadline:    hunter.Headline,
		HunterSpecialties: hunter.SpecialtiesSlice(),
		Message:           req.Message,
		Status:            req.Status,
		RequestedAt:       req.CreatedAt,
	}, nil
}

func (s *accessRequestService) Send(userID uint, dto domain.CreateAccessRequestDTO) (*domain.AccessRequestResponseDTO, error) {
	if dto.CandidateID == 0 || dto.Message == "" {
		return nil, ErrInvalidAccessRequest
	}

	hunter, err := s.hunterRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrHunterProfileMissing
	}
	if hunter.Status != domain.HunterStatusAprovado {
		return nil, ErrHunterNotApproved
	}

	req := &domain.AccessRequest{
		HunterID:    hunter.ID,
		CandidateID: dto.CandidateID,
		Message:     dto.Message,
		Status:      domain.AccessRequestPending,
	}

	if err := s.requestRepo.Create(req); err != nil {
		return nil, err
	}
	return s.toResponse(req)
}

func (s *accessRequestService) ListForCandidate(userID uint) ([]domain.AccessRequestResponseDTO, error) {
	candidate, err := s.candidateRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrCandidateProfileMissing
	}

	requests, err := s.requestRepo.FindByCandidateID(candidate.ID)
	if err != nil {
		return nil, err
	}

	dtos := make([]domain.AccessRequestResponseDTO, 0, len(requests))
	for i := range requests {
		res, err := s.toResponse(&requests[i])
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, *res)
	}
	return dtos, nil
}

func (s *accessRequestService) Respond(userID uint, id uint, status string) (*domain.AccessRequestResponseDTO, error) {
	if !isValidAccessRequestStatus(status) {
		return nil, ErrInvalidRequestStatus
	}

	candidate, err := s.candidateRepo.FindByUserID(userID)
	if err != nil {
		return nil, ErrCandidateProfileMissing
	}

	req, err := s.requestRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.CandidateID != candidate.ID {
		return nil, repository.ErrAccessRequestNotFound
	}

	req.Status = status
	if err := s.requestRepo.Update(req); err != nil {
		return nil, err
	}
	return s.toResponse(req)
}
