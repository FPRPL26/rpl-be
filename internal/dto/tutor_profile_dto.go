package dto

import "github.com/google/uuid"

type TutorRequest struct {
	Name              string `json:"name" validate:"required"`
	Semester          int    `json:"semester" validate:"required,min=1"`
	Jurusan           int64  `json:"jurusan" validate:"required"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

type TutorUpdateRequest struct {
	Name              string `json:"name"`
	Semester          int    `json:"semester" validate:"omitempty,min=1"`
	Jurusan           int64  `json:"jurusan"`
	IsVerified        *bool  `json:"is_verified"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

type TutorResponse struct {
	ID                uuid.UUID    `json:"id"`
	Name              string       `json:"name"`
	ProfilePictureURL string       `json:"profile_picture_url"`
	Semester          int          `json:"semester"`
	Jurusan           int64        `json:"jurusan"`
	Rating            float64      `json:"rating"`
	IsVerified        bool         `json:"is_verified"`
	User              UserResponse `json:"user,omitempty"`
}

type TutorListResponse struct {
	Data   []TutorResponse `json:"data"`
	Total  int64           `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}
