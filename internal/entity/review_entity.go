package entity

import (
	"time"

	"github.com/google/uuid"
)

type ReviewType string

const (
	ReviewTypeClass        ReviewType = "CLASS"
	ReviewTypeBarter       ReviewType = "BARTER"
	ReviewTypeClassRequest ReviewType = "CLASS_REQUEST"
)

type Review struct {
	ID              int64      `json:"id" gorm:"primaryKey"`
	TransactionID   uuid.UUID  `json:"transaction_id" gorm:"type:uuid;not null"`
	TransactionType ReviewType `json:"transaction_type" gorm:"not null"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	Rating          int64      `json:"rating" gorm:"not null"`
	Comment         string     `json:"comment" gorm:"not null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:timestamp without time zone;not null"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (r *Review) TableName() string {
	return "reviews"
}
