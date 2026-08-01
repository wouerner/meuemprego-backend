package service

import (
	"errors"
	"time"

	"github.com/wouerner/runter-backend/internal/domain"
	"github.com/wouerner/runter-backend/internal/repository"
)

var ErrInvalidMetric = errors.New("evento de métrica inválido")

type MetricService interface {
	Track(dto domain.TrackMetricDTO) (*domain.MetricEventResponseDTO, error)
	List() ([]domain.MetricEventResponseDTO, error)
}

type metricService struct {
	repo repository.MetricRepository
}

func NewMetricService(repo repository.MetricRepository) MetricService {
	return &metricService{repo: repo}
}

func isValidTargetType(targetType string) bool {
	return targetType == "hunter" || targetType == "candidate"
}

func isValidChannel(channel string) bool {
	return channel == "whatsapp" || channel == "linkedin"
}

func (s *metricService) Track(dto domain.TrackMetricDTO) (*domain.MetricEventResponseDTO, error) {
	if dto.TargetID == "" || dto.TargetName == "" || !isValidTargetType(dto.TargetType) || !isValidChannel(dto.Channel) {
		return nil, ErrInvalidMetric
	}

	event := &domain.MetricEvent{
		TargetType:      dto.TargetType,
		TargetID:        dto.TargetID,
		TargetName:      dto.TargetName,
		Channel:         dto.Channel,
		InitiatedByRole: dto.InitiatedByRole,
		Timestamp:       time.Now(),
	}

	if err := s.repo.Create(event); err != nil {
		return nil, err
	}
	res := event.ToResponse()
	return &res, nil
}

func (s *metricService) List() ([]domain.MetricEventResponseDTO, error) {
	events, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	dtos := make([]domain.MetricEventResponseDTO, 0, len(events))
	for _, e := range events {
		dtos = append(dtos, e.ToResponse())
	}
	return dtos, nil
}
