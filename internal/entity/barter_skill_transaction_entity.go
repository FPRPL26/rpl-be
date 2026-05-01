package entity

import (
	"time"

	"github.com/google/uuid"
)

type BarterSkillTransactionStatus string

const (
	BarterSkillTransactionStatusPending    BarterSkillTransactionStatus = "PENDING"
	BarterSkillTransactionStatusAccepted   BarterSkillTransactionStatus = "ACCEPTED"
	BarterSkillTransactionStatusRejected   BarterSkillTransactionStatus = "REJECTED"
	BarterSkillTransactionStatusOnProgress BarterSkillTransactionStatus = "ON_PROGRESS"
	BarterSkillTransactionStatusFinished   BarterSkillTransactionStatus = "FINISHED"
)

type BarterSkillTransaction struct {
	ID               uuid.UUID                    `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	BarterSkillID    uuid.UUID                    `json:"barter_skill_id" gorm:"type:uuid;not null"`
	MentorProfileID1 uuid.UUID                    `json:"mentor_profile_id1" gorm:"type:uuid;not null"`
	MentorProfileID2 uuid.UUID                    `json:"mentor_profile_id2" gorm:"type:uuid;not null"`
	Status           BarterSkillTransactionStatus `json:"status" gorm:"not null"`
	CreatedAt        time.Time                    `json:"created_at" gorm:"type:timestamp without time zone;not null"`

	BarterSkill    BarterSkill  `json:"barter_skill" gorm:"foreignKey:BarterSkillID"`
	MentorProfile1 TutorProfile `json:"mentor_profile_1" gorm:"foreignKey:MentorProfileID1"`
	MentorProfile2 TutorProfile `json:"mentor_profile_2" gorm:"foreignKey:MentorProfileID2"`
}

func (b *BarterSkillTransaction) TableName() string {
	return "barter_skill_transactions"
}
