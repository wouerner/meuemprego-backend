package repository

import (
	"github.com/wouerner/runter-backend/internal/domain"
	"gorm.io/gorm"
)

type MetricRepository interface {
	Create(event *domain.MetricEvent) error
	FindAll() ([]domain.MetricEvent, error)
}

type gormMetricRepository struct {
	db *gorm.DB
}

func NewMetricRepository(db *gorm.DB) MetricRepository {
	return &gormMetricRepository{db: db}
}

func (r *gormMetricRepository) Create(event *domain.MetricEvent) error {
	return r.db.Create(event).Error
}

func (r *gormMetricRepository) FindAll() ([]domain.MetricEvent, error) {
	var events []domain.MetricEvent
	err := r.db.Order("timestamp DESC").Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}
