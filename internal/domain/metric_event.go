package domain

import "time"

type MetricEvent struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TargetType       string    `gorm:"size:20;not null" json:"target_type"`
	TargetID         string    `gorm:"size:50;not null" json:"target_id"`
	TargetName       string    `gorm:"size:150;not null" json:"target_name"`
	Channel          string    `gorm:"size:20;not null" json:"channel"`
	InitiatedByRole  string    `gorm:"size:30;not null" json:"initiated_by_role"`
	Timestamp        time.Time `json:"timestamp"`
	CreatedAt        time.Time `json:"created_at"`
}

type TrackMetricDTO struct {
	TargetType      string `json:"target_type"`
	TargetID        string `json:"target_id"`
	TargetName      string `json:"target_name"`
	Channel         string `json:"channel"`
	InitiatedByRole string `json:"initiated_by_role"`
}

type MetricEventResponseDTO struct {
	ID              uint      `json:"id"`
	TargetType      string    `json:"target_type"`
	TargetID        string    `json:"target_id"`
	TargetName      string    `json:"target_name"`
	Channel         string    `json:"channel"`
	InitiatedByRole string    `json:"initiated_by_role"`
	Timestamp       time.Time `json:"timestamp"`
}

func (e *MetricEvent) ToResponse() MetricEventResponseDTO {
	return MetricEventResponseDTO{
		ID:              e.ID,
		TargetType:      e.TargetType,
		TargetID:        e.TargetID,
		TargetName:      e.TargetName,
		Channel:         e.Channel,
		InitiatedByRole: e.InitiatedByRole,
		Timestamp:       e.Timestamp,
	}
}
