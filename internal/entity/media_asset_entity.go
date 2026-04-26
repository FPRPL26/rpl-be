package entity

import (
	"github.com/google/uuid"
)

type MediaAsset struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	URL    string    `gorm:"type:varchar(255);not null" json:"url"`
	IsUsed bool      `gorm:"type:boolean;default:false" json:"is_used"`
	Timestamp
}
