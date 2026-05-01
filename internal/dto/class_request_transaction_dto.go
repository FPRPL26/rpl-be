package dto

type (
	CreateClassRequestTransactionRequest struct {
		RequestID     string `json:"request_id" binding:"required"`
		ApplicationID string `json:"application_id" binding:"required"`
		Price         int64  `json:"price"`
	}

	UpdateClassRequestTransactionRequest struct {
		Status string `json:"status" binding:"required"`
	}

	ClassRequestTransactionResponse struct {
		ID             string `json:"id"`
		UserID         string `json:"user_id"`
		RequestID      string `json:"request_id"`
		RequestName    string `json:"request_name,omitempty"`
		TutorProfileID string `json:"tutor_profile_id"`
		TutorName      string `json:"tutor_name,omitempty"`
		Status         string `json:"status"`
		PaymentURL     string `json:"payment_url,omitempty"`
		Price          int64  `json:"price"`
		CreatedAt      string `json:"created_at"`
	}

	ClassRequestTransactionListResponse struct {
		ID             string `json:"id"`
		UserID         string `json:"user_id"`
		RequestID      string `json:"request_id"`
		RequestName    string `json:"request_name,omitempty"`
		TutorProfileID string `json:"tutor_profile_id"`
		TutorName      string `json:"tutor_name,omitempty"`
		Status         string `json:"status"`
		Price          int64  `json:"price"`
		CreatedAt      string `json:"created_at"`
	}
)
