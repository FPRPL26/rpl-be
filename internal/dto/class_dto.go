package dto

import "github.com/google/uuid"

type CreateClassRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description" binding:"required"`
	ThumbnailURL string `json:"thumbnail_url" binding:"required"`
	ChatWA       string `json:"chat_wa"`
}

type CreateClassResponse struct {
	ID uuid.UUID `json:"id"`
}

type ClassResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	MentorID     string `json:"mentor_id"`
	MentorName   string `json:"mentor_name"`
}

type ClassDetailResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	ThumbnailURL string             `json:"thumbnail_url"`
	ChatWA       string             `json:"chat_wa"`
	MentorID     string             `json:"mentor_id"`
	MentorName   string             `json:"mentor_name"`
	Schedules    []ScheduleResponse `json:"schedules"`
}

type UpdateClassRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ThumbnailURL string `json:"thumnail_url"`
	ChatWA       string `json:"chat_wa"`
}
