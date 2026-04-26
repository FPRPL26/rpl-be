package entity

import "github.com/google/uuid"

type Class struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TutorID      uuid.UUID `json:"tutor_id" gorm:"type:uuid;not null"`
	Name         string    `json:"name" gorm:"not null"`
	Description  string    `json:"description" gorm:"not null"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"not null"`
	ChatWA       *string   `json:"chat_wa" gorm:""`

	TutorProfile TutorProfile `json:"tutor_profile" gorm:"foreignKey:TutorID"`
}

func (c *Class) TableName() string {
	return "class"
}
