package repository

import (
	"errors"
	"time"

	"github.com/wouerner/runter-backend/internal/domain"
	"gorm.io/gorm"
)

var ErrAccessRequestNotFound = errors.New("solicitação de acesso não encontrada")

type AccessRequestRepository interface {
	Create(request *domain.AccessRequest) error
	FindByID(id uint) (*domain.AccessRequest, error)
	FindByCandidateID(candidateID uint) ([]domain.AccessRequest, error)
	FindByHunterID(hunterID uint) ([]domain.AccessRequest, error)
	CountByHunterIDSince(hunterID uint, since time.Time) (int64, error)
	Update(request *domain.AccessRequest) error
}

type gormAccessRequestRepository struct {
	db *gorm.DB
}

func NewAccessRequestRepository(db *gorm.DB) AccessRequestRepository {
	return &gormAccessRequestRepository{db: db}
}

func (r *gormAccessRequestRepository) Create(request *domain.AccessRequest) error {
	return r.db.Create(request).Error
}

func (r *gormAccessRequestRepository) FindByID(id uint) (*domain.AccessRequest, error) {
	var request domain.AccessRequest
	err := r.db.First(&request, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAccessRequestNotFound
		}
		return nil, err
	}
	return &request, nil
}

func (r *gormAccessRequestRepository) FindByCandidateID(candidateID uint) ([]domain.AccessRequest, error) {
	var requests []domain.AccessRequest
	err := r.db.Where("candidate_id = ?", candidateID).Order("created_at DESC").Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *gormAccessRequestRepository) FindByHunterID(hunterID uint) ([]domain.AccessRequest, error) {
	var requests []domain.AccessRequest
	err := r.db.Where("hunter_id = ?", hunterID).Order("created_at DESC").Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *gormAccessRequestRepository) CountByHunterIDSince(hunterID uint, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&domain.AccessRequest{}).
		Where("hunter_id = ? AND created_at >= ?", hunterID, since).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *gormAccessRequestRepository) Update(request *domain.AccessRequest) error {
	return r.db.Save(request).Error
}
