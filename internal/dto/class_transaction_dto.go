package dto

import "github.com/google/uuid"

type CheckoutClassRequest struct {
	ScheduleID uuid.UUID `json:"schedule_id" binding:"required"`
}

type ClassTransactionResponse struct {
	TransactionID uuid.UUID `json:"transaction_id"`
	Status        string    `json:"status"`
	TotalPrice    int64     `json:"total_price"`
	PaymentURL    string    `json:"payment_url"`
}

type ClassTransactionListResponse struct {
	ID         uuid.UUID `json:"id"`
	ClassID    uuid.UUID `json:"class_id"`
	ClassName  string    `json:"class_name"`
	ScheduleID uuid.UUID `json:"schedule_id"`
	TotalPrice int64     `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  string    `json:"created_at"`
}
