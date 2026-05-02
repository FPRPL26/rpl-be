package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/FPRPL26/rpl-be/internal/api/repository"
	"github.com/FPRPL26/rpl-be/internal/dto"
	"github.com/FPRPL26/rpl-be/internal/entity"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TutorService interface {
		CreateTutor(ctx context.Context, userID uuid.UUID, req dto.TutorRequest) (dto.TutorResponse, error)
		GetTutorByID(ctx context.Context, id uuid.UUID) (dto.TutorResponse, error)
		UpdateTutor(ctx context.Context, id uuid.UUID, req dto.TutorUpdateRequest) (dto.TutorResponse, error)
		// DeleteTutor(ctx context.Context, id uuid.UUID) error
		// ListTutors(ctx context.Context, limit, offset int) (dto.TutorListResponse, error)
	}

	tutorService struct {
		repo repository.TutorProfileRepository
	}
)

func NewTutorService(repo repository.TutorProfileRepository) TutorService {
	return &tutorService{repo: repo}
}

func (s *tutorService) mapToResponse(tutor *entity.TutorProfile) dto.TutorResponse {
	portofolios := make([]dto.PortofolioResponse, len(tutor.Portofolios))
	for i, p := range tutor.Portofolios {
		portofolios[i] = dto.PortofolioResponse{
			ID:             p.ID.String(),
			Name:           p.Name,
			Description:    p.Description,
			FileURL:        p.FileURL,
			TutorProfileID: p.TutorProfileID.String(),
		}
	}

	return dto.TutorResponse{
		ID:                tutor.ID,
		Name:              tutor.Name,
		ProfilePictureURL: tutor.ProfilePictureURL,
		Semester:          tutor.Semester,
		Jurusan:           tutor.Jurusan,
		Rating:            tutor.Rating,
		IsVerified:        tutor.IsVerified,
		Portofolios:       portofolios,
	}
}

func (s *tutorService) CreateTutor(ctx context.Context, userID uuid.UUID, req dto.TutorRequest) (dto.TutorResponse, error) {
	tutorExisted, err := s.repo.GetByID(ctx, userID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.TutorResponse{}, err
	}
	if tutorExisted != nil {
		return dto.TutorResponse{}, myerror.New("tutor profile already exists", http.StatusConflict)
	}

	tutor := &entity.TutorProfile{
		ID:                userID,
		Name:              req.Name,
		Semester:          req.Semester,
		Jurusan:           req.Jurusan,
		ProfilePictureURL: req.ProfilePictureURL,
		Rating:            0.0,
	}

	if err := s.repo.Create(ctx, tutor); err != nil {
		return dto.TutorResponse{}, err
	}

	return s.mapToResponse(tutor), nil
}

func (s *tutorService) UpdateTutor(ctx context.Context, id uuid.UUID, req dto.TutorUpdateRequest) (dto.TutorResponse, error) {
	tutor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TutorResponse{}, myerror.New("tutor not found", http.StatusNotFound)
		}
		return dto.TutorResponse{}, err
	}

	if req.Name != "" {
		tutor.Name = req.Name
	}
	if req.Semester > 0 {
		tutor.Semester = req.Semester
	}
	if req.Jurusan != "" {
		tutor.Jurusan = req.Jurusan
	}
	if req.IsVerified != nil {
		tutor.IsVerified = *req.IsVerified
	}

	if req.ProfilePictureURL != "" {
		tutor.ProfilePictureURL = req.ProfilePictureURL
	}

	if err := s.repo.Update(ctx, tutor); err != nil {
		return dto.TutorResponse{}, err
	}

	return s.mapToResponse(tutor), nil
}

func (s *tutorService) GetTutorByID(ctx context.Context, id uuid.UUID) (dto.TutorResponse, error) {
	tutor, err := s.repo.GetByID(ctx, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.TutorResponse{}, myerror.New("tutor not found", http.StatusNotFound)
		}
		return dto.TutorResponse{}, err
	}

	return s.mapToResponse(tutor), nil
}

// func (s *tutorService) DeleteTutor(ctx context.Context, id uuid.UUID) error {
// 	_, err := s.repo.GetByID(ctx, id)
// 	if err != nil {
// 		return err
// 	}
// 	return s.repo.Delete(ctx, id)
// }

// func (s *tutorService) ListTutors(ctx context.Context, limit, offset int) (dto.TutorListResponse, error) {
// 	tutors, err := s.repo.List(ctx, limit, offset)
// 	if err != nil {
// 		return dto.TutorListResponse{}, err
// 	}

// 	data := make([]dto.TutorResponse, len(tutors))
// 	for i, t := range tutors {
// 		data[i] = s.mapToResponse(&t)
// 	}

// 	return dto.TutorListResponse{
// 		Data:   data,
// 		Limit:  limit,
// 		Offset: offset,
// 		Total:  int64(len(data)),
// 	}, nil
// }
