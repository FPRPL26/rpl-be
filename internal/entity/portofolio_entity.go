package entity

import "github.com/google/uuid"

type Portofolio struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description" gorm:"not null"`
	FileURL        string    `json:"file_url" gorm:"not null"`
	TutorProfileID uuid.UUID `json:"tutor_profile_id" gorm:"type:uuid;not null"`

	TutorProfile TutorProfile `json:"tutor_profile" gorm:"foreignKey:TutorProfileID;references:ID"`
}

func (p *Portofolio) TableName() string {
	return "portofolios"
}
