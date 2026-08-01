package domain

import "time"

const (
	AccessRequestPending  = "pending"
	AccessRequestAccepted = "accepted"
	AccessRequestRejected = "rejected"
)

type AccessRequest struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	HunterID    uint      `gorm:"index;not null" json:"hunter_id"`
	CandidateID uint      `gorm:"index;not null" json:"candidate_id"`
	Message     string    `gorm:"size:1000;not null" json:"message"`
	Status      string    `gorm:"size:20;not null;default:pending" json:"status"`
	CreatedAt   time.Time `json:"requested_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateAccessRequestDTO struct {
	CandidateID uint   `json:"candidate_id"`
	Message     string `json:"message"`
}

type RespondAccessRequestDTO struct {
	Status string `json:"status"`
}

type AccessRequestResponseDTO struct {
	ID                uint      `json:"id"`
	HunterID          uint      `json:"hunter_id"`
	CandidateID       uint      `json:"candidate_id"`
	HunterName        string    `json:"hunter_name"`
	HunterAvatar      string    `json:"hunter_avatar"`
	HunterHeadline    string    `json:"hunter_headline"`
	HunterSpecialties []string  `json:"hunter_specialties"`
	Message           string    `json:"message"`
	Status            string    `json:"status"`
	RequestedAt       time.Time `json:"requested_at"`
}
