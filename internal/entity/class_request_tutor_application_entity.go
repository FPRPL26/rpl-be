package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClassRequestTutorApplicationStatus string

const (
	ClassRequestTutorApplicationStatusPending  ClassRequestTutorApplicationStatus = "PENDING"
	ClassRequestTutorApplicationStatusAccepted ClassRequestTutorApplicationStatus = "ACCEPTED"
	ClassRequestTutorApplicationStatusRejected ClassRequestTutorApplicationStatus = "REJECTED"
)

type ClassRequestTutorApplication struct {
	ID             uuid.UUID                          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	RequestID      uuid.UUID                          `json:"request_id" gorm:"type:uuid;not null"`
	TutorProfileID uuid.UUID                          `json:"tutor_profile_id" gorm:"type:uuid;not null"`
	Status         ClassRequestTutorApplicationStatus `json:"status" gorm:"not null"`
	CreatedAt      time.Time                          `json:"created_at" gorm:"type:timestamp without time zone;not null"`
	UpdatedAt      time.Time                          `json:"updated_at" gorm:"type:timestamp without time zone;not null"`

	ClassRequest ClassRequest `json:"class_request" gorm:"foreignKey:RequestID"`
	TutorProfile TutorProfile `json:"tutor_profile" gorm:"foreignKey:TutorProfileID"`
}

func (c *ClassRequestTutorApplication) TableName() string {
	return "class_request_tutor_application"
}
