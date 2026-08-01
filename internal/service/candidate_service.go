package service

import (
	"errors"

	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
)

var ErrInvalidCandidateInput = errors.New("dados de candidato inválidos")

type CandidateService interface {
	SaveProfile(userID uint, dto domain.SaveCandidateDTO) (*domain.CandidateResponseDTO, error)
	GetProfileByUserID(userID uint) (*domain.CandidateResponseDTO, error)
	List() ([]domain.CandidateResponseDTO, error)
	SetApproval(id uint, approved bool) (*domain.CandidateResponseDTO, error)
}

type candidateService struct {
	repo repository.CandidateRepository
}

func NewCandidateService(repo repository.CandidateRepository) CandidateService {
	return &candidateService{repo: repo}
}

func (s *candidateService) SaveProfile(userID uint, dto domain.SaveCandidateDTO) (*domain.CandidateResponseDTO, error) {
	if dto.Name == "" || dto.CPF == "" || dto.Email == "" || dto.Seniority == "" || dto.Area == "" {
		return nil, ErrInvalidCandidateInput
	}

	existing, err := s.repo.FindByUserID(userID)
	if err != nil && !errors.Is(err, repository.ErrCandidateNotFound) {
		return nil, err
	}

	if existing != nil {
		if dto.CPF != existing.CPF {
			other, err := s.repo.FindByCPF(dto.CPF)
			if err == nil && other.ID != existing.ID {
				return nil, repository.ErrCandidateAlreadyExists
			}
		}

		existing.Name = dto.Name
		existing.CPF = dto.CPF
		existing.Email = dto.Email
		if dto.Password != "" {
			existing.Password = dto.Password
		}
		existing.Avatar = dto.Avatar
		existing.Headline = dto.Headline
		existing.Seniority = dto.Seniority
		existing.Area = dto.Area
		existing.CareerGoal = dto.CareerGoal
		existing.ProfessionalMoment = dto.ProfessionalMoment
		existing.RequestHunterContact = dto.RequestHunterContact
		existing.LGPDConsent = dto.LGPDConsent
		existing.LinkedInURL = dto.LinkedInURL
		existing.WhatsAppNumber = dto.WhatsAppNumber

		if err := s.repo.Update(existing); err != nil {
			return nil, err
		}
		res := existing.ToResponse()
		return &res, nil
	}

	candidate := &domain.Candidate{
		UserID:               &userID,
		Name:                 dto.Name,
		CPF:                  dto.CPF,
		Email:                dto.Email,
		Password:             dto.Password,
		Avatar:               dto.Avatar,
		Headline:             dto.Headline,
		Seniority:            dto.Seniority,
		Area:                 dto.Area,
		CareerGoal:           dto.CareerGoal,
		ProfessionalMoment:   dto.ProfessionalMoment,
		RequestHunterContact: dto.RequestHunterContact,
		LGPDConsent:          dto.LGPDConsent,
		IsApproved:           true,
		LinkedInURL:          dto.LinkedInURL,
		WhatsAppNumber:       dto.WhatsAppNumber,
	}

	if err := s.repo.Create(candidate); err != nil {
		return nil, err
	}
	res := candidate.ToResponse()
	return &res, nil
}

func (s *candidateService) GetProfileByUserID(userID uint) (*domain.CandidateResponseDTO, error) {
	candidate, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	res := candidate.ToResponse()
	return &res, nil
}

func (s *candidateService) List() ([]domain.CandidateResponseDTO, error) {
	candidates, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]domain.CandidateResponseDTO, 0, len(candidates))
	for _, c := range candidates {
		dtos = append(dtos, c.ToResponse())
	}
	return dtos, nil
}

func (s *candidateService) SetApproval(id uint, approved bool) (*domain.CandidateResponseDTO, error) {
	candidate, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	candidate.IsApproved = approved
	if err := s.repo.Update(candidate); err != nil {
		return nil, err
	}
	res := candidate.ToResponse()
	return &res, nil
}
