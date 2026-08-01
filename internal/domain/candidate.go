package domain

import "time"

type Candidate struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	UserID                *uint     `gorm:"index" json:"user_id"`
	Name                  string    `gorm:"size:120;not null" json:"name"`
	CPF                   string    `gorm:"size:14;not null;uniqueIndex" json:"cpf"`
	Email                 string    `gorm:"size:100;not null" json:"email"`
	Password              string    `gorm:"size:255" json:"-"`
	Avatar                string    `gorm:"size:500" json:"avatar"`
	Headline              string    `gorm:"size:255" json:"headline"`
	Seniority             string    `gorm:"size:40;not null" json:"seniority"`
	Area                  string    `gorm:"size:120;not null" json:"area"`
	CareerGoal            string    `gorm:"size:500" json:"career_goal"`
	ProfessionalMoment    string    `gorm:"size:60;not null" json:"professional_moment"`
	RequestHunterContact  bool      `gorm:"default:true" json:"request_hunter_contact"`
	LGPDConsent           bool      `json:"lgpd_consent"`
	IsApproved            bool      `gorm:"default:false" json:"is_approved"`
	LinkedInURL           string    `gorm:"size:500" json:"linkedin_url"`
	WhatsAppNumber        string    `gorm:"size:20" json:"whatsapp_number"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type SaveCandidateDTO struct {
	Name                 string `json:"name"`
	CPF                  string `json:"cpf"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	Avatar               string `json:"avatar"`
	Headline             string `json:"headline"`
	Seniority            string `json:"seniority"`
	Area                 string `json:"area"`
	CareerGoal           string `json:"career_goal"`
	ProfessionalMoment   string `json:"professional_moment"`
	RequestHunterContact bool   `json:"request_hunter_contact"`
	LGPDConsent          bool   `json:"lgpd_consent"`
	LinkedInURL          string `json:"linkedin_url"`
	WhatsAppNumber       string `json:"whatsapp_number"`
}

type CandidateResponseDTO struct {
	ID                   uint      `json:"id"`
	UserID               *uint     `json:"user_id"`
	Name                 string    `json:"name"`
	CPF                  string    `json:"cpf"`
	Email                string    `json:"email"`
	Avatar               string    `json:"avatar"`
	Headline             string    `json:"headline"`
	Seniority            string    `json:"seniority"`
	Area                 string    `json:"area"`
	CareerGoal           string    `json:"career_goal"`
	ProfessionalMoment   string    `json:"professional_moment"`
	RequestHunterContact bool      `json:"request_hunter_contact"`
	LGPDConsent          bool      `json:"lgpd_consent"`
	IsApproved           bool      `json:"is_approved"`
	LinkedInURL          string    `json:"linkedin_url"`
	WhatsAppNumber       string    `json:"whatsapp_number"`
	CreatedAt            time.Time `json:"created_at"`
}

func (c *Candidate) ToResponse() CandidateResponseDTO {
	return CandidateResponseDTO{
		ID:                   c.ID,
		UserID:               c.UserID,
		Name:                 c.Name,
		CPF:                  c.CPF,
		Email:                c.Email,
		Avatar:               c.Avatar,
		Headline:             c.Headline,
		Seniority:            c.Seniority,
		Area:                 c.Area,
		CareerGoal:           c.CareerGoal,
		ProfessionalMoment:   c.ProfessionalMoment,
		RequestHunterContact: c.RequestHunterContact,
		LGPDConsent:          c.LGPDConsent,
		IsApproved:           c.IsApproved,
		LinkedInURL:          c.LinkedInURL,
		WhatsAppNumber:       c.WhatsAppNumber,
		CreatedAt:            c.CreatedAt,
	}
}
