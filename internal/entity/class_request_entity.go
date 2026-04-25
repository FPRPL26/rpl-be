package entity

import (
	"time"

	"github.com/google/uuid"
)

type ClassRequestStatus string

const (
	ClassRequestStatusWaiting  ClassRequestStatus = "WAITING"
	ClassRequestStatusFinished ClassRequestStatus = "FINISHED"
)

type ClassRequest struct {
	ID          uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID      uuid.UUID          `json:"user_id" gorm:"type:uuid;not null"`
	Name        string             `json:"name" gorm:"not null"`
	Description string             `json:"description" gorm:"not null"`
	Start       time.Time          `json:"start" gorm:"type:timestamp without time zone;not null"`
	End         time.Time          `json:"end" gorm:"type:timestamp without time zone;not null"`
	Date        time.Time          `json:"date" gorm:"type:date;not null"`
	Status      ClassRequestStatus `json:"status" gorm:"not null"`
	Price       int64              `json:"price" gorm:"not null"`
	ChatWA      string             `json:"chat_wa" gorm:"not null"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}

func (c *ClassRequest) TableName() string {
	return "class_requests"
}
