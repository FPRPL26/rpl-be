package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClassTransactionStatus string

const (
	ClassTransactionStatusSuccess   ClassTransactionStatus = "SUCCESS"
	ClassTransactionStatusPending   ClassTransactionStatus = "PENDING"
	ClassTransactionStatusCancelled ClassTransactionStatus = "CANCELLED"
)

type ClassTransaction struct {
	ID         uuid.UUID              `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ClassID    uuid.UUID              `json:"class_id" gorm:"type:uuid;not null"`
	ScheduleID uuid.UUID              `json:"schedule_id" gorm:"type:uuid;not null"`
	UserID     uuid.UUID              `json:"user_id" gorm:"type:uuid;not null"`
	TotalPrice int64                  `json:"total_price" gorm:"not null"`
	Status     ClassTransactionStatus `json:"status" gorm:"not null"`
	CreatedAt  time.Time              `json:"created_at" gorm:"type:timestamp without time zone;not null"`

	Class    Class    `json:"class" gorm:"foreignKey:ClassID"`
	Schedule Schedule `json:"schedule" gorm:"foreignKey:ScheduleID"`
	User     User     `json:"user" gorm:"foreignKey:UserID"`
}

func (c *ClassTransaction) TableName() string {
	return "class_transactions"
}
