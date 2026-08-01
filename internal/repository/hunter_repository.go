package repository

import (
	"errors"

	"github.com/wouerner/runter-backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrHunterNotFound      = errors.New("job hunter não encontrado")
	ErrHunterAlreadyExists = errors.New("job hunter com este CPF já cadastrado")
)

type HunterRepository interface {
	Create(hunter *domain.Hunter) error
	FindByID(id uint) (*domain.Hunter, error)
	FindByUserID(userID uint) (*domain.Hunter, error)
	FindByCPF(cpf string) (*domain.Hunter, error)
	FindAll() ([]domain.Hunter, error)
	Update(hunter *domain.Hunter) error
}

type gormHunterRepository struct {
	db *gorm.DB
}

func NewHunterRepository(db *gorm.DB) HunterRepository {
	return &gormHunterRepository{db: db}
}

func (r *gormHunterRepository) Create(hunter *domain.Hunter) error {
	var existing domain.Hunter
	if err := r.db.Where("cpf = ?", hunter.CPF).First(&existing).Error; err == nil {
		return ErrHunterAlreadyExists
	}

	return r.db.Create(hunter).Error
}

func (r *gormHunterRepository) FindByID(id uint) (*domain.Hunter, error) {
	var hunter domain.Hunter
	err := r.db.First(&hunter, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHunterNotFound
		}
		return nil, err
	}
	return &hunter, nil
}

func (r *gormHunterRepository) FindByUserID(userID uint) (*domain.Hunter, error) {
	var hunter domain.Hunter
	err := r.db.Where("user_id = ?", userID).First(&hunter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHunterNotFound
		}
		return nil, err
	}
	return &hunter, nil
}

func (r *gormHunterRepository) FindByCPF(cpf string) (*domain.Hunter, error) {
	var hunter domain.Hunter
	err := r.db.Where("cpf = ?", cpf).First(&hunter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrHunterNotFound
		}
		return nil, err
	}
	return &hunter, nil
}

func (r *gormHunterRepository) FindAll() ([]domain.Hunter, error) {
	var hunters []domain.Hunter
	err := r.db.Order("created_at ASC").Find(&hunters).Error
	if err != nil {
		return nil, err
	}
	return hunters, nil
}

func (r *gormHunterRepository) Update(hunter *domain.Hunter) error {
	return r.db.Save(hunter).Error
}
