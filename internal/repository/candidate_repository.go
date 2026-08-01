package repository

import (
	"errors"

	"github.com/wouerner/runter-backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrCandidateNotFound      = errors.New("candidato não encontrado")
	ErrCandidateAlreadyExists = errors.New("candidato com este CPF já cadastrado")
)

type CandidateRepository interface {
	Create(candidate *domain.Candidate) error
	FindByID(id uint) (*domain.Candidate, error)
	FindByUserID(userID uint) (*domain.Candidate, error)
	FindByCPF(cpf string) (*domain.Candidate, error)
	FindAll() ([]domain.Candidate, error)
	Update(candidate *domain.Candidate) error
}

type gormCandidateRepository struct {
	db *gorm.DB
}

func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &gormCandidateRepository{db: db}
}

func (r *gormCandidateRepository) Create(candidate *domain.Candidate) error {
	var existing domain.Candidate
	if err := r.db.Where("cpf = ?", candidate.CPF).First(&existing).Error; err == nil {
		return ErrCandidateAlreadyExists
	}

	return r.db.Create(candidate).Error
}

func (r *gormCandidateRepository) FindByID(id uint) (*domain.Candidate, error) {
	var candidate domain.Candidate
	err := r.db.First(&candidate, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCandidateNotFound
		}
		return nil, err
	}
	return &candidate, nil
}

func (r *gormCandidateRepository) FindByUserID(userID uint) (*domain.Candidate, error) {
	var candidate domain.Candidate
	err := r.db.Where("user_id = ?", userID).First(&candidate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCandidateNotFound
		}
		return nil, err
	}
	return &candidate, nil
}

func (r *gormCandidateRepository) FindByCPF(cpf string) (*domain.Candidate, error) {
	var candidate domain.Candidate
	err := r.db.Where("cpf = ?", cpf).First(&candidate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCandidateNotFound
		}
		return nil, err
	}
	return &candidate, nil
}

func (r *gormCandidateRepository) FindAll() ([]domain.Candidate, error) {
	var candidates []domain.Candidate
	err := r.db.Order("created_at ASC").Find(&candidates).Error
	if err != nil {
		return nil, err
	}
	return candidates, nil
}

func (r *gormCandidateRepository) Update(candidate *domain.Candidate) error {
	return r.db.Save(candidate).Error
}
