package service

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/api/repository"
	"github.com/FPRPL26/rpl-be/internal/dto"
	"github.com/FPRPL26/rpl-be/internal/entity"
	myerror "github.com/FPRPL26/rpl-be/internal/pkg/error"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	PortofolioService interface {
		Create(ctx context.Context, tutorProfileID string, req dto.CreatePortofolioRequest) (dto.PortofolioResponse, error)
		GetAll(ctx context.Context, tutorProfileID string) ([]dto.PortofolioResponse, error)
		GetAllByTutorProfile(ctx context.Context, tutorProfileID string) ([]dto.PortofolioResponse, error)
		GetById(ctx context.Context, id string) (dto.PortofolioResponse, error)
		Update(ctx context.Context, tutorProfileID string, id string, req dto.UpdatePortofolioRequest) (dto.PortofolioResponse, error)
		Delete(ctx context.Context, tutorProfileID string, id string) error
	}

	portofolioService struct {
		repo repository.PortofolioRepository
		db   *gorm.DB
	}
)

func NewPortofolio(repo repository.PortofolioRepository, db *gorm.DB) PortofolioService {
	return &portofolioService{repo: repo, db: db}
}

func (s *portofolioService) Create(ctx context.Context, tutorProfileID string, req dto.CreatePortofolioRequest) (dto.PortofolioResponse, error) {
	tutorUUID, err := uuid.Parse(tutorProfileID)
	if err != nil {
		return dto.PortofolioResponse{}, err
	}

	portofolio := entity.Portofolio{
		ID:             uuid.New(),
		Name:           req.Name,
		Description:    req.Description,
		FileURL:        req.FileURL,
		TutorProfileID: tutorUUID,
	}

	created, err := s.repo.Create(ctx, nil, portofolio)
	if err != nil {
		return dto.PortofolioResponse{}, err
	}

	return s.mapToResponse(created), nil
}

func (s *portofolioService) GetAll(ctx context.Context, tutorProfileID string) ([]dto.PortofolioResponse, error) {
	portofolios, err := s.repo.GetAll(ctx, nil, tutorProfileID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PortofolioResponse, 0, len(portofolios))
	for _, p := range portofolios {
		responses = append(responses, s.mapToResponse(p))
	}

	return responses, nil
}

func (s *portofolioService) GetAllByTutorProfile(ctx context.Context, tutorProfileID string) ([]dto.PortofolioResponse, error) {
	portofolios, err := s.repo.GetAllByTutorProfileId(ctx, nil, tutorProfileID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PortofolioResponse, 0, len(portofolios))
	for _, p := range portofolios {
		responses = append(responses, s.mapToResponse(p))
	}

	return responses, nil
}

func (s *portofolioService) GetById(ctx context.Context, id string) (dto.PortofolioResponse, error) {
	portofolio, err := s.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.PortofolioResponse{}, err
	}

	return s.mapToResponse(portofolio), nil
}

func (s *portofolioService) Update(ctx context.Context, tutorProfileID string, id string, req dto.UpdatePortofolioRequest) (dto.PortofolioResponse, error) {
	portofolio, err := s.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.PortofolioResponse{}, err
	}

	if portofolio.TutorProfileID.String() != tutorProfileID {
		return dto.PortofolioResponse{}, myerror.New("unauthorized to update portofolio", 403)
	}

	if req.Name != "" {
		portofolio.Name = req.Name
	}
	if req.Description != "" {
		portofolio.Description = req.Description
	}
	if req.FileURL != "" {
		portofolio.FileURL = req.FileURL
	}

	updated, err := s.repo.Update(ctx, nil, portofolio)
	if err != nil {
		return dto.PortofolioResponse{}, err
	}

	return s.mapToResponse(updated), nil
}

func (s *portofolioService) Delete(ctx context.Context, tutorProfileID string, id string) error {
	portofolio, err := s.repo.GetById(ctx, nil, id)
	if err != nil {
		return err
	}

	if portofolio.TutorProfileID.String() != tutorProfileID {
		return myerror.New("unauthorized to delete portofolio", 403)
	}

	return s.repo.Delete(ctx, nil, portofolio)
}

func (s *portofolioService) mapToResponse(portofolio entity.Portofolio) dto.PortofolioResponse {
	return dto.PortofolioResponse{
		ID:             portofolio.ID.String(),
		Name:           portofolio.Name,
		Description:    portofolio.Description,
		FileURL:        portofolio.FileURL,
		TutorProfileID: portofolio.TutorProfileID.String(),
	}
}
