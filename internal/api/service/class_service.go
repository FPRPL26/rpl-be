package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/FPRPL26/rpl-be/internal/api/repository"
	"github.com/FPRPL26/rpl-be/internal/dto"
	"github.com/FPRPL26/rpl-be/internal/entity"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ClassService interface {
		Create(ctx context.Context, tutorId string, req dto.CreateClassRequest) (entity.Class, error)
		GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.ClassResponse, meta.Meta, error)
		GetAllEnrolled(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.ClassResponse, meta.Meta, error)
		GetAllByTutorId(ctx context.Context, tutorId string, metaReq meta.Meta) ([]dto.ClassResponse, meta.Meta, error)
		GetById(ctx context.Context, classId string) (dto.ClassDetailResponse, error)
		GetSchedules(ctx context.Context, metaReq meta.Meta, classId string) ([]dto.ScheduleResponse, meta.Meta, error)
		Update(ctx context.Context, tutorId string, classId string, req dto.UpdateClassRequest) (dto.ClassResponse, error)
		Delete(ctx context.Context, tutorId string, classId string) error
		AddSchedules(ctx context.Context, classId string, req dto.AddSchedulesRequest) error
		UpdateSchedule(ctx context.Context, tutorId string, scheduleId string, req dto.UpdateScheduleRequest) (dto.ScheduleResponse, error)
		DeleteSchedule(ctx context.Context, tutorId string, scheduleId string) error
	}

	classService struct {
		classRepo       repository.ClassRepository
		scheduleRepo    repository.ScheduleRepository
		transactionRepo repository.ClassTransactionRepository
		db              *gorm.DB
	}
)

func NewClass(classRepo repository.ClassRepository, scheduleRepo repository.ScheduleRepository, transactionRepo repository.ClassTransactionRepository, db *gorm.DB) ClassService {
	return &classService{classRepo, scheduleRepo, transactionRepo, db}
}

func (s *classService) Create(ctx context.Context, tutorId string, req dto.CreateClassRequest) (entity.Class, error) {
	tutorUUID, err := uuid.Parse(tutorId)
	if err != nil {
		return entity.Class{}, err
	}

	class := entity.Class{
		TutorID:      tutorUUID,
		Name:         req.Name,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		Price:        req.Price,
	}

	if req.ChatWA != "" {
		class.ChatWA = &req.ChatWA
	}

	return s.classRepo.Create(ctx, nil, class)
}

func (s *classService) GetAll(ctx context.Context, metaReq meta.Meta) ([]dto.ClassResponse, meta.Meta, error) {
	classes, meta, err := s.classRepo.GetAll(ctx, nil, metaReq, "TutorProfile")
	if err != nil {
		return nil, meta, err
	}

	classResponses := make([]dto.ClassResponse, 0, len(classes))
	for _, class := range classes {
		classResponses = append(classResponses, dto.ClassResponse{
			ID:           class.ID.String(),
			Name:         class.Name,
			ThumbnailURL: class.ThumbnailURL,
			MentorID:     class.TutorID.String(),
			MentorName:   class.TutorProfile.Name,
			Price:        class.Price,
		})
	}

	return classResponses, meta, nil
}

func (s *classService) GetAllEnrolled(ctx context.Context, userId string, metaReq meta.Meta) ([]dto.ClassResponse, meta.Meta, error) {
	transactions, meta, err := s.transactionRepo.GetAllByUserId(ctx, nil, userId, metaReq, "Class", "Class.TutorProfile")
	if err != nil {
		return nil, meta, err
	}

	classResponses := make([]dto.ClassResponse, 0, len(transactions))
	for _, tx := range transactions {
		classResponses = append(classResponses, dto.ClassResponse{
			ID:           tx.Class.ID.String(),
			Name:         tx.Class.Name,
			ThumbnailURL: tx.Class.ThumbnailURL,
			MentorID:     tx.Class.TutorID.String(),
			MentorName:   tx.Class.TutorProfile.Name,
			Price:        tx.Class.Price,
		})
	}

	return classResponses, meta, nil
}

func (s *classService) GetAllByTutorId(ctx context.Context, tutorId string, metaReq meta.Meta) ([]dto.ClassResponse, meta.Meta, error) {
	tutorUUID, err := uuid.Parse(tutorId)
	if err != nil {
		return nil, metaReq, err
	}

	tx := s.db.Where("tutor_id = ?", tutorUUID)

	classes, meta, err := s.classRepo.GetAll(ctx, tx, metaReq, "TutorProfile")
	if err != nil {
		return nil, meta, err
	}

	classResponses := make([]dto.ClassResponse, 0, len(classes))
	for _, class := range classes {
		classResponses = append(classResponses, dto.ClassResponse{
			ID:           class.ID.String(),
			Name:         class.Name,
			ThumbnailURL: class.ThumbnailURL,
			MentorID:     class.TutorID.String(),
			MentorName:   class.TutorProfile.Name,
			Price:        class.Price,
		})
	}

	return classResponses, meta, nil
}

func (s *classService) GetById(ctx context.Context, classId string) (dto.ClassDetailResponse, error) {
	class, err := s.classRepo.GetById(ctx, nil, classId, "TutorProfile")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassDetailResponse{}, myerror.New("class not found", http.StatusNotFound)
		}
		return dto.ClassDetailResponse{}, err
	}

	chatWA := ""
	if class.ChatWA != nil {
		chatWA = *class.ChatWA
	}

	schedules, _, err := s.scheduleRepo.GetAllByClassId(ctx, nil, meta.Default(), classId)
	if err != nil {
		return dto.ClassDetailResponse{}, err
	}

	scheduleResponses := make([]dto.ScheduleResponse, 0, len(schedules))
	for _, sch := range schedules {
		scheduleResponses = append(scheduleResponses, dto.ScheduleResponse{
			ID:         sch.ID.String(),
			ClassID:    sch.ClassID.String(),
			StartTime:  sch.StartTime.Format(time.RFC3339),
			EndTime:    sch.EndTime.Format(time.RFC3339),
			Date:       sch.Date.Format("02-01-2006"),
			MaxStudent: sch.MaxStudent,
			Remaining:  sch.Remaining,
			Repeted:    sch.Repeted,
		})
	}

	return dto.ClassDetailResponse{
		ID:           class.ID.String(),
		Name:         class.Name,
		Description:  class.Description,
		ThumbnailURL: class.ThumbnailURL,
		ChatWA:       chatWA,
		Price:        class.Price,
		MentorID:     class.TutorID.String(),
		MentorName:   class.TutorProfile.Name,
		Schedules:    scheduleResponses,
	}, nil
}

func (s *classService) GetSchedules(ctx context.Context, metaReq meta.Meta, classId string) ([]dto.ScheduleResponse, meta.Meta, error) {
	// Verify class exists
	_, err := s.classRepo.GetById(ctx, nil, classId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, metaReq, myerror.New("class not found", http.StatusNotFound)
		}
		return nil, metaReq, err
	}

	schedules, metaRes, err := s.scheduleRepo.GetAllByClassId(ctx, nil, metaReq, classId)
	if err != nil {
		return nil, metaRes, err
	}

	scheduleResponses := make([]dto.ScheduleResponse, 0, len(schedules))
	for _, sch := range schedules {
		scheduleResponses = append(scheduleResponses, dto.ScheduleResponse{
			ID:         sch.ID.String(),
			ClassID:    sch.ClassID.String(),
			StartTime:  sch.StartTime.Format(time.RFC3339),
			EndTime:    sch.EndTime.Format(time.RFC3339),
			Date:       sch.Date.Format("02-01-2006"),
			MaxStudent: sch.MaxStudent,
			Remaining:  sch.Remaining,
			Repeted:    sch.Repeted,
		})
	}

	return scheduleResponses, metaRes, nil
}

func (s *classService) Update(ctx context.Context, tutorId string, classId string, req dto.UpdateClassRequest) (dto.ClassResponse, error) {
	class, err := s.classRepo.GetById(ctx, nil, classId)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	if class.TutorID.String() != tutorId {
		return dto.ClassResponse{}, myerror.New("unauthorized", 403)
	}

	if req.Name != "" {
		class.Name = req.Name
	}
	if req.Description != "" {
		class.Description = req.Description
	}
	if req.ThumbnailURL != "" {
		class.ThumbnailURL = req.ThumbnailURL
	}
	if req.ChatWA != "" {
		class.ChatWA = &req.ChatWA
	}
	if req.Price != 0 {
		class.Price = req.Price
	}

	updatedClass, err := s.classRepo.Update(ctx, nil, class)
	if err != nil {
		return dto.ClassResponse{}, err
	}

	return dto.ClassResponse{
		ID:           updatedClass.ID.String(),
		Name:         updatedClass.Name,
		ThumbnailURL: updatedClass.ThumbnailURL,
		MentorID:     updatedClass.TutorID.String(),
		MentorName:   updatedClass.TutorProfile.Name,
		Price:        updatedClass.Price,
	}, nil
}

func (s *classService) Delete(ctx context.Context, tutorId string, classId string) error {
	class, err := s.classRepo.GetById(ctx, nil, classId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.New("class not found", http.StatusNotFound)
		}
		return err
	}

	if class.TutorID.String() != tutorId {
		return myerror.New("unauthorized", 403)
	}

	return s.classRepo.Delete(ctx, nil, class)
}

func (s *classService) AddSchedules(ctx context.Context, classId string, req dto.AddSchedulesRequest) error {
	classUUID, err := uuid.Parse(classId)
	if err != nil {
		return err
	}

	// Verify class exists
	_, err = s.classRepo.GetById(ctx, nil, classId)
	if err != nil {
		return err
	}

	var schedules []entity.Schedule
	for _, schReq := range req.Schedules {
		startTime, err := time.Parse(time.RFC3339, schReq.StartTime)
		if err != nil {
			return err
		}

		endTime, err := time.Parse(time.RFC3339, schReq.EndTime)
		if err != nil {
			return err
		}

		date, err := time.Parse("02-01-2006", schReq.Date)
		if err != nil {
			return err
		}

		schedules = append(schedules, entity.Schedule{
			ClassID:    classUUID,
			StartTime:  startTime,
			EndTime:    endTime,
			Date:       date,
			MaxStudent: schReq.MaxStudent,
			Remaining:  schReq.MaxStudent,
			Repeted:    schReq.Repeted,
		})
	}

	return s.scheduleRepo.AddSchedules(ctx, nil, schedules)
}

func (s *classService) UpdateSchedule(ctx context.Context, tutorId string, scheduleId string, req dto.UpdateScheduleRequest) (dto.ScheduleResponse, error) {
	schedule, err := s.scheduleRepo.GetById(ctx, nil, scheduleId, "Class")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ScheduleResponse{}, myerror.New("class not found", http.StatusNotFound)
		}
		return dto.ScheduleResponse{}, err
	}

	if schedule.Class.TutorID.String() != tutorId {
		return dto.ScheduleResponse{}, myerror.New("unauthorized", 403)
	}

	if req.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, req.StartTime)
		if err != nil {
			return dto.ScheduleResponse{}, err
		}
		schedule.StartTime = startTime
	}

	if req.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, req.EndTime)
		if err != nil {
			return dto.ScheduleResponse{}, err
		}
		schedule.EndTime = endTime
	}

	if req.Date != "" {
		date, err := time.Parse("02-01-2006", req.Date)
		if err != nil {
			return dto.ScheduleResponse{}, err
		}
		schedule.Date = date
	}

	if req.MaxStudent != 0 {
		schedule.MaxStudent = req.MaxStudent
	}

	if req.Repeted != 0 {
		schedule.Repeted = req.Repeted
	}

	updatedSchedule, err := s.scheduleRepo.Update(ctx, nil, schedule)
	if err != nil {
		return dto.ScheduleResponse{}, err
	}

	return dto.ScheduleResponse{
		ID:         updatedSchedule.ID.String(),
		ClassID:    updatedSchedule.ClassID.String(),
		StartTime:  updatedSchedule.StartTime.Format(time.RFC3339),
		EndTime:    updatedSchedule.EndTime.Format(time.RFC3339),
		Date:       updatedSchedule.Date.Format("02-01-2006"),
		MaxStudent: updatedSchedule.MaxStudent,
		Remaining:  updatedSchedule.Remaining,
		Repeted:    updatedSchedule.Repeted,
	}, nil
}

func (s *classService) DeleteSchedule(ctx context.Context, tutorId string, scheduleId string) error {
	schedule, err := s.scheduleRepo.GetById(ctx, nil, scheduleId, "Class")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerror.New("schedule not found", http.StatusNotFound)
		}
		return err
	}

	if schedule.Class.TutorID.String() != tutorId {
		return myerror.New("unauthorized", 403)
	}

	return s.scheduleRepo.Delete(ctx, nil, schedule)
}
