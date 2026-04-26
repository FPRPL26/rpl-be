package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClassRequestTransactionStatus string

const (
	ClassRequestTransactionStatusPending    ClassRequestTransactionStatus = "PENDING"
	ClassRequestTransactionStatusAccepted   ClassRequestTransactionStatus = "ACCEPTED"
	ClassRequestTransactionStatusOnProgress ClassRequestTransactionStatus = "ON_PROGRESS"
	ClassRequestTransactionStatusFinished   ClassRequestTransactionStatus = "FINISHED"
)

type ClassRequestTransaction struct {
	ID            uuid.UUID                     `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID        uuid.UUID                     `json:"user_id" gorm:"type:uuid;not null"`
	RequestID     uuid.UUID                     `json:"request_id" gorm:"type:uuid;not null"`
	TutorProfileID uuid.UUID                    `json:"tutor_profile" gorm:"type:uuid;not null"`
	Status        ClassRequestTransactionStatus `json:"status" gorm:"not null"`
	CreatedAt     time.Time                     `json:"created_at" gorm:"type:timestamp without time zone;not null"`
	Price         int64                         `json:"price" gorm:"not null"`

	User         User         `json:"user" gorm:"foreignKey:UserID"`
	ClassRequest ClassRequest `json:"class_request" gorm:"foreignKey:RequestID"`
	TutorProfile TutorProfile `json:"tutor_profile_data" gorm:"foreignKey:TutorProfileID"`
}

func (c *ClassRequestTransaction) TableName() string {
	return "class_request_transaction"
}
