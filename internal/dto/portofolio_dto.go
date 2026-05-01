package dto

type CreatePortofolioRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	FileURL     string `json:"file_url" validate:"required,url"`
}

type UpdatePortofolioRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FileURL     string `json:"file_url"`
}

type PortofolioResponse struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	FileURL        string `json:"file_url"`
	TutorProfileID string `json:"tutor_profile_id"`
}

type PortofolioListResponse struct {
	Data []PortofolioResponse `json:"data"`
}
