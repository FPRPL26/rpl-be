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
	ReviewService interface {
		SubmitReview(ctx context.Context, userID string, req dto.SubmitReviewRequest) (dto.SubmitReviewResponse, error)
	}

	reviewService struct {
		reviewRepo            repository.ReviewRepository
		classTransRepo        repository.ClassTransactionRepository
		barterTransRepo       repository.BarterSkillTransactionRepository
		classRequestTransRepo repository.ClassRequestTransactionRepository
	}
)

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	classTransRepo repository.ClassTransactionRepository,
	barterTransRepo repository.BarterSkillTransactionRepository,
	classRequestTransRepo repository.ClassRequestTransactionRepository,
) ReviewService {
	return &reviewService{
		reviewRepo:            reviewRepo,
		classTransRepo:        classTransRepo,
		barterTransRepo:       barterTransRepo,
		classRequestTransRepo: classRequestTransRepo,
	}
}

func (s *reviewService) SubmitReview(ctx context.Context, userID string, req dto.SubmitReviewRequest) (dto.SubmitReviewResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.SubmitReviewResponse{}, err
	}

	transUUID, err := uuid.Parse(req.TransactionID)
	if err != nil {
		return dto.SubmitReviewResponse{}, err
	}

	// 1. Get Transaction and Validate
	switch entity.ReviewType(req.TransactionType) {
	case entity.ReviewTypeClass:
		transaction, err := s.classTransRepo.GetById(ctx, nil, req.TransactionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.SubmitReviewResponse{}, myerror.New("class transaction not found", http.StatusNotFound)
			}
			return dto.SubmitReviewResponse{}, err
		}

		if transaction.UserID != userUUID {
			return dto.SubmitReviewResponse{}, myerror.New("you are not authorized to review this transaction", http.StatusForbidden)
		}

		if transaction.Status != entity.ClassTransactionStatusSuccess {
			return dto.SubmitReviewResponse{}, myerror.New("you can only review completed classes", http.StatusBadRequest)
		}

	case entity.ReviewTypeBarter:
		transaction, err := s.barterTransRepo.GetById(ctx, nil, req.TransactionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.SubmitReviewResponse{}, myerror.New("barter transaction not found", http.StatusNotFound)
			}
			return dto.SubmitReviewResponse{}, err
		}

		if transaction.MentorProfileID1 != userUUID && transaction.MentorProfileID2 != userUUID {
			return dto.SubmitReviewResponse{}, myerror.New("you are not authorized to review this transaction", http.StatusForbidden)
		}

		if transaction.Status != entity.BarterSkillTransactionStatusFinished {
			return dto.SubmitReviewResponse{}, myerror.New("you can only review finished barter transactions", http.StatusBadRequest)
		}

	case entity.ReviewTypeClassRequest:
		transaction, err := s.classRequestTransRepo.GetById(ctx, nil, req.TransactionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.SubmitReviewResponse{}, myerror.New("class request transaction not found", http.StatusNotFound)
			}
			return dto.SubmitReviewResponse{}, err
		}

		if transaction.UserID != userUUID {
			return dto.SubmitReviewResponse{}, myerror.New("you are not authorized to review this transaction", http.StatusForbidden)
		}

		if transaction.Status != entity.ClassRequestTransactionStatusFinished {
			return dto.SubmitReviewResponse{}, myerror.New("you can only review finished class requests", http.StatusBadRequest)
		}

	default:
		return dto.SubmitReviewResponse{}, myerror.New("invalid transaction type", http.StatusBadRequest)
	}

	// 2. Check if already reviewed
	_, err = s.reviewRepo.GetByTransactionId(ctx, nil, req.TransactionID)
	if err == nil {
		return dto.SubmitReviewResponse{}, myerror.New("you have already reviewed this transaction", http.StatusBadRequest)
	}

	// 3. Create Review
	review := entity.Review{
		TransactionID:   transUUID,
		TransactionType: entity.ReviewType(req.TransactionType),
		UserID:          userUUID,
		Rating:          req.Rating,
		Comment:         req.Comment,
		CreatedAt:       time.Now(),
	}

	created, err := s.reviewRepo.Create(ctx, nil, review)
	if err != nil {
		return dto.SubmitReviewResponse{}, err
	}

	return dto.SubmitReviewResponse{
		ReviewID: created.ID,
	}, nil
}
