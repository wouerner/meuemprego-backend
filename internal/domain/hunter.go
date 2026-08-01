package domain

import (
	"encoding/json"
	"time"
)

type Hunter struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UserID              *uint     `gorm:"index" json:"user_id"`
	Name                string    `gorm:"size:120;not null" json:"name"`
	CPF                 string    `gorm:"size:14;not null;uniqueIndex" json:"cpf"`
	Email               string    `gorm:"size:100;not null" json:"email"`
	Password            string    `gorm:"size:255" json:"-"`
	Avatar              string    `gorm:"size:500" json:"avatar"`
	Headline            string    `gorm:"size:255" json:"headline"`
	Bio                 string    `gorm:"size:1000" json:"bio"`
	Specialties         string    `gorm:"type:text" json:"-"`
	SenioritiesServed   string    `gorm:"type:text" json:"-"`
	ServiceModel        string    `gorm:"size:60;not null" json:"service_model"`
	Status              string    `gorm:"size:20;not null;default:Pendente" json:"status"`
	Rating              float64   `gorm:"default:0" json:"rating"`
	TotalContactsCount  int       `gorm:"default:0" json:"total_contacts_count"`
	LinkedInURL         string    `gorm:"size:500" json:"linkedin_url"`
	WhatsAppNumber      string    `gorm:"size:20" json:"whatsapp_number"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

const (
	HunterStatusPendente  = "Pendente"
	HunterStatusAprovado  = "Aprovado"
	HunterStatusRejeitado = "Rejeitado"
)

type SaveHunterDTO struct {
	Name               string   `json:"name"`
	CPF                string   `json:"cpf"`
	Email              string   `json:"email"`
	Password           string   `json:"password"`
	Avatar             string   `json:"avatar"`
	Headline           string   `json:"headline"`
	Bio                string   `json:"bio"`
	Specialties        []string `json:"specialties"`
	SenioritiesServed  []string `json:"seniorities_served"`
	ServiceModel       string   `json:"service_model"`
	LinkedInURL        string   `json:"linkedin_url"`
	WhatsAppNumber     string   `json:"whatsapp_number"`
}

type HunterResponseDTO struct {
	ID                 uint      `json:"id"`
	UserID             *uint     `json:"user_id"`
	Name               string    `json:"name"`
	CPF                string    `json:"cpf"`
	Email              string    `json:"email"`
	Avatar             string    `json:"avatar"`
	Headline           string    `json:"headline"`
	Bio                string    `json:"bio"`
	Specialties        []string  `json:"specialties"`
	SenioritiesServed  []string  `json:"seniorities_served"`
	ServiceModel       string    `json:"service_model"`
	Status             string    `json:"status"`
	Rating             float64   `json:"rating"`
	TotalContactsCount int       `json:"total_contacts_count"`
	LinkedInURL        string    `json:"linkedin_url"`
	WhatsAppNumber     string    `json:"whatsapp_number"`
	CreatedAt          time.Time `json:"created_at"`
}

func marshalStrings(values []string) string {
	data, _ := json.Marshal(values)
	return string(data)
}

func unmarshalStrings(data string) []string {
	if data == "" {
		return []string{}
	}
	var values []string
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		return []string{}
	}
	return values
}

func (h *Hunter) SpecialtiesSlice() []string {
	return unmarshalStrings(h.Specialties)
}

func (h *Hunter) ToResponse() HunterResponseDTO {
	return HunterResponseDTO{
		ID:                 h.ID,
		UserID:             h.UserID,
		Name:               h.Name,
		CPF:                h.CPF,
		Email:              h.Email,
		Avatar:             h.Avatar,
		Headline:           h.Headline,
		Bio:                h.Bio,
		Specialties:        unmarshalStrings(h.Specialties),
		SenioritiesServed:  unmarshalStrings(h.SenioritiesServed),
		ServiceModel:       h.ServiceModel,
		Status:             h.Status,
		Rating:             h.Rating,
		TotalContactsCount: h.TotalContactsCount,
		LinkedInURL:        h.LinkedInURL,
		WhatsAppNumber:     h.WhatsAppNumber,
		CreatedAt:          h.CreatedAt,
	}
}

func (h *Hunter) ApplyDTO(dto SaveHunterDTO) {
	h.Name = dto.Name
	h.CPF = dto.CPF
	h.Email = dto.Email
	if dto.Password != "" {
		h.Password = dto.Password
	}
	h.Avatar = dto.Avatar
	h.Headline = dto.Headline
	h.Bio = dto.Bio
	h.Specialties = marshalStrings(dto.Specialties)
	h.SenioritiesServed = marshalStrings(dto.SenioritiesServed)
	h.ServiceModel = dto.ServiceModel
	h.LinkedInURL = dto.LinkedInURL
	h.WhatsAppNumber = dto.WhatsAppNumber
}
