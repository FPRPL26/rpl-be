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
	ClassRequestService interface {
		Create(ctx context.Context, userId string, req dto.CreateClassRequestRequest) (entity.ClassRequest, error)
		GetAll(ctx context.Context) ([]dto.ClassRequestResponse, error)
		GetById(ctx context.Context, id string) (dto.ClassRequestResponse, error)
		Update(ctx context.Context, userId string, id string, req dto.UpdateClassRequestRequest) (dto.ClassRequestResponse, error)
		Delete(ctx context.Context, userId string, id string) error
	}

	classRequestService struct {
		repo repository.ClassRequestRepository
	}
)

func NewClassRequestService(repo repository.ClassRequestRepository) ClassRequestService {
	return &classRequestService{repo}
}

func (s *classRequestService) Create(ctx context.Context, userId string, req dto.CreateClassRequestRequest) (entity.ClassRequest, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return entity.ClassRequest{}, err
	}

	start, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		return entity.ClassRequest{}, err
	}
	end, err := time.Parse(time.RFC3339, req.End)
	if err != nil {
		return entity.ClassRequest{}, err
	}
	date, err := time.Parse("02-01-2006", req.Date)
	if err != nil {
		return entity.ClassRequest{}, err
	}

	cr := entity.ClassRequest{
		UserID:      userUUID,
		Name:        req.Name,
		Description: req.Description,
		Start:       start,
		End:         end,
		Date:        date,
		Status:      entity.ClassRequestStatusWaiting,
		Price:       req.Price,
		ChatWA:      req.ChatWA,
	}

	return s.repo.Create(ctx, nil, cr)
}

func (s *classRequestService) GetAll(ctx context.Context) ([]dto.ClassRequestResponse, error) {
	items, _, err := s.repo.GetAll(ctx, nil, meta.Default())
	if err != nil {
		return nil, err
	}

	res := make([]dto.ClassRequestResponse, 0, len(items))
	for _, it := range items {
		res = append(res, dto.ClassRequestResponse{
			ID:          it.ID.String(),
			UserID:      it.UserID.String(),
			Name:        it.Name,
			Description: it.Description,
			Start:       it.Start.Format(time.RFC3339),
			End:         it.End.Format(time.RFC3339),
			Date:        it.Date.Format("02-01-2006"),
			Status:      string(it.Status),
			Price:       it.Price,
			ChatWA:      it.ChatWA,
		})
	}

	return res, nil
}

func (s *classRequestService) GetById(ctx context.Context, id string) (dto.ClassRequestResponse, error) {
	it, err := s.repo.GetById(ctx, nil, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestResponse{}, myerror.New("class request not found", http.StatusNotFound)
		}
		return dto.ClassRequestResponse{}, err
	}

	return dto.ClassRequestResponse{
		ID:          it.ID.String(),
		UserID:      it.UserID.String(),
		Name:        it.Name,
		Description: it.Description,
		Start:       it.Start.Format(time.RFC3339),
		End:         it.End.Format(time.RFC3339),
		Date:        it.Date.Format("02-01-2006"),
		Status:      string(it.Status),
		Price:       it.Price,
		ChatWA:      it.ChatWA,
	}, nil
}

func (s *classRequestService) Update(ctx context.Context, userId string, id string, req dto.UpdateClassRequestRequest) (dto.ClassRequestResponse, error) {
	it, err := s.repo.GetById(ctx, nil, id)
	if err != nil {
		return dto.ClassRequestResponse{}, err
	}

	if it.UserID.String() != userId {
		return dto.ClassRequestResponse{}, myerror.New("unauthorized", 403)
	}

	if req.Name != "" {
		it.Name = req.Name
	}
	if req.Description != "" {
		it.Description = req.Description
	}
	if req.Start != "" {
		if t, err := time.Parse(time.RFC3339, req.Start); err == nil {
			it.Start = t
		}
	}
	if req.End != "" {
		if t, err := time.Parse(time.RFC3339, req.End); err == nil {
			it.End = t
		}
	}
	if req.Date != "" {
		if d, err := time.Parse("02-01-2006", req.Date); err == nil {
			it.Date = d
		}
	}
	if req.Status != "" {
		it.Status = entity.ClassRequestStatus(req.Status)
	}
	if req.Price != 0 {
		it.Price = req.Price
	}
	if req.ChatWA != "" {
		it.ChatWA = req.ChatWA
	}

	updated, err := s.repo.Update(ctx, nil, it)
	if err != nil {
		return dto.ClassRequestResponse{}, err
	}

	return dto.ClassRequestResponse{
		ID:          updated.ID.String(),
		UserID:      updated.UserID.String(),
		Name:        updated.Name,
		Description: updated.Description,
		Start:       updated.Start.Format(time.RFC3339),
		End:         updated.End.Format(time.RFC3339),
		Date:        updated.Date.Format("02-01-2006"),
		Status:      string(updated.Status),
		Price:       updated.Price,
		ChatWA:      updated.ChatWA,
	}, nil
}

func (s *classRequestService) Delete(ctx context.Context, userId string, id string) error {
	it, err := s.repo.GetById(ctx, nil, id)
	if err != nil {
		return err
	}

	if it.UserID.String() != userId {
		return myerror.New("unauthorized", 403)
	}

	return s.repo.Delete(ctx, nil, it)
}
