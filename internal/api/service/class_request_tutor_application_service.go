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
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ClassRequestTutorApplicationService interface {
		Apply(ctx context.Context, userID string, req dto.CreateClassRequestTutorApplicationRequest) (dto.ClassRequestTutorApplicationResponse, error)
		GetAllByRequest(ctx context.Context, userID, requestID string) ([]dto.ClassRequestTutorApplicationResponse, error)
		GetAllByTutorProfile(ctx context.Context, tutorProfileID string) ([]dto.ClassRequestTutorApplicationResponse, error)
		GetById(ctx context.Context, userID, id string) (dto.ClassRequestTutorApplicationResponse, error)
		UpdateStatus(ctx context.Context, userID, id string, req dto.UpdateClassRequestTutorApplicationStatusRequest) (dto.ClassRequestTutorApplicationResponse, error)
	}

	classRequestTutorApplicationService struct {
		db               *gorm.DB
		applicationRepo  repository.ClassRequestTutorApplicationRepository
		classRequestRepo repository.ClassRequestRepository
		tutorProfileRepo repository.TutorProfileRepository
		userRepo         repository.UserRepository
	}
)

func NewClassRequestTutorApplicationService(
	db *gorm.DB,
	applicationRepo repository.ClassRequestTutorApplicationRepository,
	classRequestRepo repository.ClassRequestRepository,
	tutorProfileRepo repository.TutorProfileRepository,
	userRepo repository.UserRepository,
) ClassRequestTutorApplicationService {
	return &classRequestTutorApplicationService{
		db:               db,
		applicationRepo:  applicationRepo,
		classRequestRepo: classRequestRepo,
		tutorProfileRepo: tutorProfileRepo,
		userRepo:         userRepo,
	}
}

func (s *classRequestTutorApplicationService) Apply(ctx context.Context, userID string, req dto.CreateClassRequestTutorApplicationRequest) (dto.ClassRequestTutorApplicationResponse, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("invalid user ID", http.StatusBadRequest)
	}

	tutorProfileID, err := uuid.Parse(req.TutorProfileID)
	if err != nil {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("invalid tutor profile ID", http.StatusBadRequest)
	}

	if tutorProfileID.String() != userID {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("unauthorized tutor profile", http.StatusForbidden)
	}

	if _, err := s.tutorProfileRepo.GetByID(ctx, tutorProfileID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTutorApplicationResponse{}, myerror.New("tutor profile not found", http.StatusNotFound)
		}
		return dto.ClassRequestTutorApplicationResponse{}, err
	}

	requestID, err := uuid.Parse(req.RequestID)
	if err != nil {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("invalid request ID", http.StatusBadRequest)
	}

	classRequest, err := s.classRequestRepo.GetById(ctx, nil, req.RequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTutorApplicationResponse{}, myerror.New("class request not found", http.StatusNotFound)
		}
		return dto.ClassRequestTutorApplicationResponse{}, err
	}

	if classRequest.UserID.String() == tutorProfileID.String() {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("cannot apply to your own request", http.StatusForbidden)
	}

	if _, err := s.applicationRepo.GetByTutorAndRequest(ctx, nil, req.TutorProfileID, req.RequestID); err == nil {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("application already exists for this tutor and request", http.StatusBadRequest)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.ClassRequestTutorApplicationResponse{}, err
	}

	application := entity.ClassRequestTutorApplication{
		ID:             uuid.New(),
		RequestID:      requestID,
		TutorProfileID: tutorProfileID,
		Status:         entity.ClassRequestTutorApplicationStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	created, err := s.applicationRepo.Create(ctx, nil, application)
	if err != nil {
		return dto.ClassRequestTutorApplicationResponse{}, err
	}

	return dto.ClassRequestTutorApplicationResponse{
		ID:             created.ID.String(),
		RequestID:      created.RequestID.String(),
		TutorProfileID: created.TutorProfileID.String(),
		Status:         string(created.Status),
		CreatedAt:      created.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      created.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *classRequestTutorApplicationService) GetAllByRequest(ctx context.Context, userID, requestID string) ([]dto.ClassRequestTutorApplicationResponse, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, myerror.New("invalid user ID", http.StatusBadRequest)
	}

	classRequest, err := s.classRequestRepo.GetById(ctx, nil, requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerror.New("class request not found", http.StatusNotFound)
		}
		return nil, err
	}

	if classRequest.UserID.String() != userID {
		return nil, myerror.New("unauthorized", http.StatusForbidden)
	}

	applications, err := s.applicationRepo.GetAllByRequestId(ctx, nil, requestID, "ClassRequest", "TutorProfile")
	if err != nil {
		return nil, err
	}

	res := make([]dto.ClassRequestTutorApplicationResponse, 0, len(applications))
	for _, application := range applications {
		res = append(res, dto.ClassRequestTutorApplicationResponse{
			ID:             application.ID.String(),
			RequestID:      application.RequestID.String(),
			RequestName:    application.ClassRequest.Name,
			TutorProfileID: application.TutorProfileID.String(),
			TutorName:      application.TutorProfile.Name,
			Status:         string(application.Status),
			CreatedAt:      application.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      application.UpdatedAt.Format(time.RFC3339),
		})
	}

	return res, nil
}

func (s *classRequestTutorApplicationService) GetAllByTutorProfile(ctx context.Context, tutorProfileID string) ([]dto.ClassRequestTutorApplicationResponse, error) {
	if _, err := uuid.Parse(tutorProfileID); err != nil {
		return nil, myerror.New("invalid tutor profile ID", http.StatusBadRequest)
	}

	applications, err := s.applicationRepo.GetAllByTutorProfileId(ctx, nil, tutorProfileID, "ClassRequest", "TutorProfile")
	if err != nil {
		return nil, err
	}

	res := make([]dto.ClassRequestTutorApplicationResponse, 0, len(applications))
	for _, application := range applications {
		res = append(res, dto.ClassRequestTutorApplicationResponse{
			ID:             application.ID.String(),
			RequestID:      application.RequestID.String(),
			RequestName:    application.ClassRequest.Name,
			TutorProfileID: application.TutorProfileID.String(),
			TutorName:      application.TutorProfile.Name,
			Status:         string(application.Status),
			CreatedAt:      application.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      application.UpdatedAt.Format(time.RFC3339),
		})
	}

	return res, nil
}

func (s *classRequestTutorApplicationService) GetById(ctx context.Context, userID, id string) (dto.ClassRequestTutorApplicationResponse, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("invalid user ID", http.StatusBadRequest)
	}

	application, err := s.applicationRepo.GetById(ctx, nil, id, "ClassRequest", "TutorProfile")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTutorApplicationResponse{}, myerror.New("application not found", http.StatusNotFound)
		}
		return dto.ClassRequestTutorApplicationResponse{}, err
	}

	if application.ClassRequest.UserID.String() != userID && application.TutorProfileID.String() != userID {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("unauthorized", http.StatusForbidden)
	}

	return dto.ClassRequestTutorApplicationResponse{
		ID:             application.ID.String(),
		RequestID:      application.RequestID.String(),
		RequestName:    application.ClassRequest.Name,
		TutorProfileID: application.TutorProfileID.String(),
		TutorName:      application.TutorProfile.Name,
		Status:         string(application.Status),
		CreatedAt:      application.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      application.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *classRequestTutorApplicationService) UpdateStatus(ctx context.Context, userID, id string, req dto.UpdateClassRequestTutorApplicationStatusRequest) (dto.ClassRequestTutorApplicationResponse, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("invalid user ID", http.StatusBadRequest)
	}

	status := entity.ClassRequestTutorApplicationStatus(req.Status)
	switch status {
	case entity.ClassRequestTutorApplicationStatusAccepted,
		entity.ClassRequestTutorApplicationStatusRejected:
		break
	default:
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("invalid status", http.StatusBadRequest)
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	application, err := s.applicationRepo.GetById(ctx, tx, id, "ClassRequest", "TutorProfile")
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTutorApplicationResponse{}, myerror.New("application not found", http.StatusNotFound)
		}
		return dto.ClassRequestTutorApplicationResponse{}, err
	}

	if application.ClassRequest.UserID.String() != userID {
		tx.Rollback()
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("unauthorized", http.StatusForbidden)
	}

	if application.Status != entity.ClassRequestTutorApplicationStatusPending {
		tx.Rollback()
		return dto.ClassRequestTutorApplicationResponse{}, myerror.New("only pending applications can be updated", http.StatusBadRequest)
	}

	if status == entity.ClassRequestTutorApplicationStatusAccepted {
		var existing entity.ClassRequestTutorApplication
		if err := tx.WithContext(ctx).Where("request_id = ? AND status = ?", application.RequestID, entity.ClassRequestTutorApplicationStatusAccepted).First(&existing).Error; err == nil {
			tx.Rollback()
			return dto.ClassRequestTutorApplicationResponse{}, myerror.New("request already has an accepted tutor", http.StatusBadRequest)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return dto.ClassRequestTutorApplicationResponse{}, err
		}

		application.Status = status
		application.UpdatedAt = time.Now()
		updated, err := s.applicationRepo.Update(ctx, tx, application)
		if err != nil {
			tx.Rollback()
			return dto.ClassRequestTutorApplicationResponse{}, err
		}

		application = updated

		application.ClassRequest.TutorProfileID = application.TutorProfileID
		application.ClassRequest.Status = entity.ClassRequestStatusAssigned
		if _, err := s.classRequestRepo.Update(ctx, tx, application.ClassRequest); err != nil {
			tx.Rollback()
			return dto.ClassRequestTutorApplicationResponse{}, err
		}

		if err := tx.WithContext(ctx).Model(&entity.ClassRequestTutorApplication{}).
			Where("request_id = ? AND id != ?", application.RequestID, application.ID).
			Updates(map[string]interface{}{
				"status":     entity.ClassRequestTutorApplicationStatusRejected,
				"updated_at": time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			return dto.ClassRequestTutorApplicationResponse{}, err
		}

		if err := tx.Commit().Error; err != nil {
			return dto.ClassRequestTutorApplicationResponse{}, err
		}
	} else {
		application.Status = status
		application.UpdatedAt = time.Now()
		updated, err := s.applicationRepo.Update(ctx, tx, application)
		if err != nil {
			tx.Rollback()
			return dto.ClassRequestTutorApplicationResponse{}, err
		}
		if err := tx.Commit().Error; err != nil {
			return dto.ClassRequestTutorApplicationResponse{}, err
		}
		application = updated
	}

	return dto.ClassRequestTutorApplicationResponse{
		ID:             application.ID.String(),
		RequestID:      application.RequestID.String(),
		RequestName:    application.ClassRequest.Name,
		TutorProfileID: application.TutorProfileID.String(),
		TutorName:      application.TutorProfile.Name,
		Status:         string(application.Status),
		CreatedAt:      application.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      application.UpdatedAt.Format(time.RFC3339),
	}, nil
}
