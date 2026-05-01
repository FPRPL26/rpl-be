package dto

import "github.com/google/uuid"

type SubmitReviewRequest struct {
	TransactionID   string `json:"transaction_id" binding:"required"`
	Rating          int64  `json:"rating" binding:"required,min=1,max=5"`
	Comment         string `json:"comment" binding:"required"`
	TransactionType string `json:"transaction_type" binding:"required"`
}

type SubmitReviewResponse struct {
	ReviewID int64 `json:"review_id"`
}

type ReviewResponse struct {
	ID        int64     `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"user_name"`
	Rating    int64     `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt string    `json:"created_at"`
}
