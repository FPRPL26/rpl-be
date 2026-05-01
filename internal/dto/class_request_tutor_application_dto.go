package dto

type (
	CreateClassRequestTutorApplicationRequest struct {
		RequestID      string `json:"request_id" binding:"required"`
		TutorProfileID string `json:"tutor_profile_id" binding:"required"`
	}

	UpdateClassRequestTutorApplicationStatusRequest struct {
		Status string `json:"status" binding:"required"`
	}

	ClassRequestTutorApplicationResponse struct {
		ID             string `json:"id"`
		RequestID      string `json:"request_id"`
		RequestName    string `json:"request_name,omitempty"`
		TutorProfileID string `json:"tutor_profile_id"`
		TutorName      string `json:"tutor_name,omitempty"`
		Status         string `json:"status"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at,omitempty"`
	}
)
