package entity

import (
	"time"

	"github.com/google/uuid"
)

type BarterSkill struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RequestSkills  int64     `json:"request_skills" gorm:"not null"`
	OfferedSkills  int64     `json:"offered_skills" gorm:"not null"`
	Name           string    `json:"name" gorm:"not null"`
	Description    string    `json:"description" gorm:"not null"`
	Accepted       bool      `json:"accepted" gorm:"not null"`
	TutorProfileID uuid.UUID `json:"tutor_profile_id" gorm:"type:uuid;not null"`
	ChatWA         *string   `json:"chat_wa" gorm:""`

	TutorProfile TutorProfile `json:"tutor_profile" gorm:"foreignKey:TutorProfileID"`
	RequestSkill Skill        `json:"request_skill" gorm:"foreignKey:RequestSkills"`
	OfferedSkill Skill        `json:"offered_skill" gorm:"foreignKey:OfferedSkills"`

	CreatedAt time.Time `gorm:"type:timestamp without time zone" json:"created_at"`
}

func (b *BarterSkill) TableName() string {
	return "barter_skills"
}
