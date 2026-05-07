package dto

type (
	CreateClassRequestRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		Start       string `json:"start" binding:"required"`
		End         string `json:"end" binding:"required"`
		Date        string `json:"date" binding:"required"`
		Price       int64  `json:"price" binding:"required"`
		ChatWA      string `json:"chat_wa" binding:"required"`
	}

	UpdateClassRequestRequest struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Start       string `json:"start"`
		End         string `json:"end"`
		Date        string `json:"date"`
		Status      string `json:"status"`
		Price       int64  `json:"price"`
		ChatWA      string `json:"chat_wa"`
	}

	ClassRequestResponse struct {
		ID          string `json:"id"`
		UserID      string `json:"user_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Start       string `json:"start"`
		End         string `json:"end"`
		Date        string `json:"date"`
		Status      string `json:"status"`
		Price       int64  `json:"price"`
		ChatWA      string `json:"chat_wa"`
	}
)
