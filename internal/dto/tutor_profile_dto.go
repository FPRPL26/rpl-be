package dto

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type TutorRequest struct {
	Name     string `json:"name" validate:"required"`
	Semester int    `json:"semester" validate:"required,min=1"`
	Jurusan  int64  `json:"jurusan" validate:"required"`
}

type TutorUpdateRequest struct {
	Name           string                `form:"name"`
	Semester       int                   `form:"semester" validate:"omitempty,min=1"`
	Jurusan        int64                 `form:"jurusan"`
	IsVerified     *bool                 `form:"is_verified"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture"`
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
