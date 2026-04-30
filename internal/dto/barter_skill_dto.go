package dto

import "github.com/google/uuid"

type CreateBarterOfferRequest struct {
	RequestSkillID int64  `json:"request_skill_id" binding:"required"`
	OfferedSkillID int64  `json:"offered_skill_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	ChatWA         string `json:"chat_wa" binding:"required"`
}

type CreateBarterOfferResponse struct {
	BarterID uuid.UUID `json:"barter_id"`
}

type RequestBarterOfferResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
}

type ApproveBarterRequestRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
}

type BarterRequestResponse struct {
	TransactionID    uuid.UUID `json:"transaction_id"`
	Status           string    `json:"status"`
	BarterID         uuid.UUID `json:"barter_id"`
	TutorID          uuid.UUID `json:"tutor_id"`
	TutorName        string    `json:"tutor_name"`
	RequestSkillName string    `json:"request_skill_name"`
	OfferedSkillName string    `json:"offered_skill_name"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
}

type BarterIncomingRequestResponse struct {
	TransactionID    uuid.UUID `json:"transaction_id"`
	Status           string    `json:"status"`
	BarterID         uuid.UUID `json:"barter_id"`
	RequesterID      uuid.UUID `json:"requester_id"`
	RequesterName    string    `json:"requester_name"`
	RequestSkillName string    `json:"request_skill_name"`
	OfferedSkillName string    `json:"offered_skill_name"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
}

type BarterOfferResponse struct {
	ID               uuid.UUID `json:"id"`
	TutorID          uuid.UUID `json:"tutor_id"`
	TutorName        string    `json:"tutor_name"`
	RequestSkillID   int64     `json:"request_skill_id"`
	RequestSkillName string    `json:"request_skill_name"`
	OfferedSkillID   int64     `json:"offered_skill_id"`
	OfferedSkillName string    `json:"offered_skill_name"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	ChatWA           string    `json:"chat_wa"`
	Accepted         bool      `json:"accepted"`
}

type BarterTransactionResponse struct {
	ID         uuid.UUID `json:"id"`
	MentorID   uuid.UUID `json:"mentor_id"`
	MentorName string    `json:"mentor_name"`
	Status     string    `json:"status"`
}
