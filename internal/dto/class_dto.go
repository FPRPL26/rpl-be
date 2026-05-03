package dto

import "github.com/google/uuid"

type Mentor struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	ProfilePictureURL string `json:"profile_picture_url"`
	IsVerified        bool   `json:"is_verified"`
}

type CreateClassRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url" binding:"required"`
	ChatWA       string `json:"chat_wa"`
	Price        int64  `json:"price" binding:"required"`
}

type CreateClassResponse struct {
	ID uuid.UUID `json:"id"`
}

type ClassResponse struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Price        int64    `json:"price"`
	Rating       *float64 `json:"rating"`
	Mentor       Mentor   `json:"mentor"`
}

type ClassDetailResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	ThumbnailURL string             `json:"thumbnail_url"`
	ChatWA       string             `json:"chat_wa"`
	Price        int64              `json:"price"`
	Rating       *float64           `json:"rating"`
	Reviews      []ReviewResponse   `json:"reviews"`
	Mentor       Mentor             `json:"mentor"`
	Schedules    []ScheduleResponse `json:"schedules"`
}

type UpdateClassRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumbnail_url"`
	ChatWA       string `json:"chat_wa"`
	Price        int64  `json:"price"`
}
