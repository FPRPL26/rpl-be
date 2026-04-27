package service

import (
	"context"
	"fmt"

	"github.com/FPRPL26/rpl-be/internal/api/repository"
	"github.com/FPRPL26/rpl-be/internal/dto"
	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"github.com/google/uuid"
)

type (
	TutorService interface {
		CreateTutor(ctx context.Context, userID uuid.UUID, req dto.TutorRequest) (dto.TutorResponse, error)
		UpdateTutor(ctx context.Context, id uuid.UUID, req dto.TutorUpdateRequest) (dto.TutorResponse, error)
		GetTutorByID(ctx context.Context, id uuid.UUID) (dto.TutorResponse, error)
		DeleteTutor(ctx context.Context, id uuid.UUID) error
		ListTutors(ctx context.Context, limit, offset int) (dto.TutorListResponse, error)
	}

	tutorService struct {
		repo repository.TutorProfileRepository
	}
)

func NewTutorService(repo repository.TutorProfileRepository) TutorService {
	return &tutorService{repo: repo}
}

func (s *tutorService) CreateTutor(ctx context.Context, userID uuid.UUID, req dto.TutorRequest) (dto.TutorResponse, error) {
	tutor := &entity.TutorProfile{
		ID:       userID,
		Name:     req.Name,
		Semester: req.Semester,
		Jurusan:  req.Jurusan,
		Rating:   0.0,
	}

	if err := s.repo.Create(ctx, tutor); err != nil {
		return dto.TutorResponse{}, err
	}

	return s.mapToResponse(tutor), nil
}

func (s *tutorService) UpdateTutor(ctx context.Context, id uuid.UUID, req dto.TutorUpdateRequest) (dto.TutorResponse, error) {
	tutor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.TutorResponse{}, err
	}

	if req.Name != "" {
		tutor.Name = req.Name
	}
	if req.Semester > 0 {
		tutor.Semester = req.Semester
	}
	if req.Jurusan > 0 {
		tutor.Jurusan = req.Jurusan
	}
	if req.IsVerified != nil {
		tutor.IsVerified = *req.IsVerified
	}

	if req.ProfilePicture != nil {
		if tutor.ProfilePictureURL != "" {
			_ = utils.DeleteFile(tutor.ProfilePictureURL)
		}

		ext := utils.GetExtensions(req.ProfilePicture.Filename)
		relativeID := fmt.Sprintf("profiles/%s.%s", id.String(), ext)

		if err := utils.UploadFile(req.ProfilePicture, relativeID); err != nil {
			return dto.TutorResponse{}, err
		}

		tutor.ProfilePictureURL = relativeID
	}

	if err := s.repo.Update(ctx, tutor); err != nil {
		return dto.TutorResponse{}, err
	}

	return s.mapToResponse(tutor), nil
}

func (s *tutorService) GetTutorByID(ctx context.Context, id uuid.UUID) (dto.TutorResponse, error) {
	tutor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.TutorResponse{}, err
	}
	return s.mapToResponse(tutor), nil
}

func (s *tutorService) DeleteTutor(ctx context.Context, id uuid.UUID) error {
	tutor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if tutor.ProfilePictureURL != "" {
		_ = utils.DeleteFile(tutor.ProfilePictureURL)
	}

	return s.repo.Delete(ctx, id)
}

func (s *tutorService) ListTutors(ctx context.Context, limit, offset int) (dto.TutorListResponse, error) {
	tutors, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return dto.TutorListResponse{}, err
	}

	data := make([]dto.TutorResponse, len(tutors))
	for i, t := range tutors {
		data[i] = s.mapToResponse(&t)
	}

	return dto.TutorListResponse{
		Data:   data,
		Limit:  limit,
		Offset: offset,
		Total:  int64(len(data)),
	}, nil
}

func (s *tutorService) mapToResponse(t *entity.TutorProfile) dto.TutorResponse {
	return dto.TutorResponse{
		ID:                t.ID,
		Name:              t.Name,
		ProfilePictureURL: t.ProfilePictureURL,
		Semester:          t.Semester,
		Jurusan:           t.Jurusan,
		Rating:            t.Rating,
		IsVerified:        t.IsVerified,
		User: dto.UserResponse{
			Email: t.User.Email,
		},
	}
}
