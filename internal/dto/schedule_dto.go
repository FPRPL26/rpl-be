package dto

type ScheduleRequest struct {
	StartTime  string `json:"start_time" validate:"required"`
	EndTime    string `json:"end_time" validate:"required"`
	Date       string `json:"date" validate:"required"`
	MaxStudent int64  `json:"max_student" validate:"required"`
	Repeted    int    `json:"repeted"`
}

type AddSchedulesRequest struct {
	Schedules []ScheduleRequest `json:"schedules" validate:"required,dive"`
}

type UpdateScheduleRequest struct {
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Date       string `json:"date"`
	MaxStudent int64  `json:"max_student"`
	Repeted    int    `json:"repeted"`
}

type ScheduleResponse struct {
	ID         string `json:"id"`
	ClassID    string `json:"class_id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Date       string `json:"date"`
	MaxStudent int64  `json:"max_student"`
	Repeted    int    `json:"repeted"`
}
