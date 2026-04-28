package dto

import "github.com/google/uuid"

type CheckoutClassRequest struct {
	ScheduleID uuid.UUID `json:"schedule_id" binding:"required"`
}

type ClassTransactionResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	TotalPrice    int64     `json:"total_price"`
}
