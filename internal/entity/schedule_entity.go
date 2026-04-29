package entity

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ClassID    uuid.UUID `json:"class_id" gorm:"type:uuid;not null"`
	StartTime  time.Time `json:"start_time" gorm:"type:timestamp without time zone;not null"`
	EndTime    time.Time `json:"end_time" gorm:"type:timestamp without time zone;not null"`
	Date       time.Time `json:"date" gorm:"type:date;not null"`
	Repeted    int       `json:"repeted" gorm:"not null"`
	MaxStudent int64     `json:"max_student" gorm:"not null"`
	Remaining  int64     `json:"remaining" gorm:"not null"`

	Class Class `json:"class" gorm:"foreignKey:ClassID"`
}

func (s *Schedule) TableName() string {
	return "schedules"
}
