package entity

import "github.com/google/uuid"

type Review struct {
	ID                 int64     `json:"id" gorm:"primaryKey"`
	ClassTransactionID uuid.UUID `json:"class_transaction_id" gorm:"type:uuid;not null"`
	UserID             uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Rating             int64     `json:"rating" gorm:"not null"`
	Comment            string    `json:"comment" gorm:"not null"`
	ScheduleID         uuid.UUID `json:"schedule_id" gorm:"type:uuid;not null"`

	ClassTransaction ClassTransaction `json:"class_transaction" gorm:"foreignKey:ClassTransactionID"`
	User             User             `json:"user" gorm:"foreignKey:UserID"`
	Schedule         Schedule         `json:"schedule" gorm:"foreignKey:ScheduleID"`
}

func (r *Review) TableName() string {
	return "reviews"
}
