package service

import (
	"errors"

	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
)

var (
	ErrInvalidHunterInput  = errors.New("dados de job hunter inválidos")
	ErrInvalidHunterStatus = errors.New("status inválido para job hunter")
)

type HunterService interface {
	SaveProfile(userID uint, dto domain.SaveHunterDTO) (*domain.HunterResponseDTO, error)
	GetProfileByUserID(userID uint) (*domain.HunterResponseDTO, error)
	GetByID(id uint) (*domain.HunterResponseDTO, error)
	List() ([]domain.HunterResponseDTO, error)
	SetStatus(id uint, status string) (*domain.HunterResponseDTO, error)
	IncrementContacts(id uint) (*domain.HunterResponseDTO, error)
}

type hunterService struct {
	repo repository.HunterRepository
}

func NewHunterService(repo repository.HunterRepository) HunterService {
	return &hunterService{repo: repo}
}

func isValidHunterStatus(status string) bool {
	return status == domain.HunterStatusPendente ||
		status == domain.HunterStatusAprovado ||
		status == domain.HunterStatusRejeitado
}

func (s *hunterService) SaveProfile(userID uint, dto domain.SaveHunterDTO) (*domain.HunterResponseDTO, error) {
	if dto.Name == "" || dto.CPF == "" || dto.Email == "" || dto.ServiceModel == "" {
		return nil, ErrInvalidHunterInput
	}

	existing, err := s.repo.FindByUserID(userID)
	if err != nil && !errors.Is(err, repository.ErrHunterNotFound) {
		return nil, err
	}

	if existing != nil {
		if dto.CPF != existing.CPF {
			other, err := s.repo.FindByCPF(dto.CPF)
			if err == nil && other.ID != existing.ID {
				return nil, repository.ErrHunterAlreadyExists
			}
		}

		existing.ApplyDTO(dto)
		if err := s.repo.Update(existing); err != nil {
			return nil, err
		}
		res := existing.ToResponse()
		return &res, nil
	}

	hunter := &domain.Hunter{UserID: &userID, Status: domain.HunterStatusPendente, Rating: 0}
	hunter.ApplyDTO(dto)

	if err := s.repo.Create(hunter); err != nil {
		return nil, err
	}
	res := hunter.ToResponse()
	return &res, nil
}

func (s *hunterService) GetProfileByUserID(userID uint) (*domain.HunterResponseDTO, error) {
	hunter, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	res := hunter.ToResponse()
	return &res, nil
}

func (s *hunterService) GetByID(id uint) (*domain.HunterResponseDTO, error) {
	hunter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	res := hunter.ToResponse()
	return &res, nil
}

func (s *hunterService) List() ([]domain.HunterResponseDTO, error) {
	hunters, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]domain.HunterResponseDTO, 0, len(hunters))
	for _, h := range hunters {
		dtos = append(dtos, h.ToResponse())
	}
	return dtos, nil
}

func (s *hunterService) SetStatus(id uint, status string) (*domain.HunterResponseDTO, error) {
	if !isValidHunterStatus(status) {
		return nil, ErrInvalidHunterStatus
	}

	hunter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	hunter.Status = status
	if err := s.repo.Update(hunter); err != nil {
		return nil, err
	}
	res := hunter.ToResponse()
	return &res, nil
}

func (s *hunterService) IncrementContacts(id uint) (*domain.HunterResponseDTO, error) {
	hunter, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	hunter.TotalContactsCount++
	if err := s.repo.Update(hunter); err != nil {
		return nil, err
	}
	res := hunter.ToResponse()
	return &res, nil
}
