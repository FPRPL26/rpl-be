package entity

import "github.com/google/uuid"

type TutorProfile struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name         string    `json:"name" gorm:"not null"`
	Semester     int       `json:"semester" gorm:"not null"`
	Jurusan      int64     `json:"jurusan" gorm:"not null"`
	Rating       float64   `json:"rating" gorm:"type:decimal(8,2);not null"`
	IsVerified   bool      `json:"is_verified" gorm:"not null"`
	PortofolioID string    `json:"portofolio_id" gorm:"not null"`
}

func (t *TutorProfile) TableName() string {
	return "tutor_profile"
}
