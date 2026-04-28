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
	ClassTransactionService interface {
		Checkout(ctx context.Context, userID string, req dto.CheckoutClassRequest) (dto.ClassTransactionResponse, error)
		Complete(ctx context.Context, transactionID string) error
	}

	classTransactionService struct {
		db              *gorm.DB
		transactionRepo repository.ClassTransactionRepository
		scheduleRepo    repository.ScheduleRepository
		classRepo       repository.ClassRepository
	}
)

func NewClassTransactionService(
	db *gorm.DB,
	transactionRepo repository.ClassTransactionRepository,
	scheduleRepo repository.ScheduleRepository,
	classRepo repository.ClassRepository,
) ClassTransactionService {
	return &classTransactionService{
		db:              db,
		transactionRepo: transactionRepo,
		scheduleRepo:    scheduleRepo,
		classRepo:       classRepo,
	}
}

func (s *classTransactionService) Checkout(ctx context.Context, userID string, req dto.CheckoutClassRequest) (dto.ClassTransactionResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.ClassTransactionResponse{}, err
	}

	var transactionResult dto.ClassTransactionResponse

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Check if user already registered for this schedule
		existing, err := s.transactionRepo.GetByUserAndSchedule(ctx, tx, userID, req.ScheduleID.String())
		if err == nil && existing.ID != uuid.Nil {
			return myerror.New("you have already registered for this schedule", http.StatusBadRequest)
		}

		// 2. Get Schedule with Lock for capacity check
		schedule, err := s.scheduleRepo.GetById(ctx, tx.Set("gorm:query_option", "FOR UPDATE"), req.ScheduleID.String(), "Class")
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return myerror.New("schedule not found", http.StatusNotFound)
			}
			return err
		}

		// 2. Check Capacity
		if schedule.Remaining <= 0 {
			return myerror.New("class is full", http.StatusBadRequest)
		}

		// 3. Create Transaction
		transaction := entity.ClassTransaction{
			ID:         uuid.New(),
			ClassID:    schedule.ClassID,
			ScheduleID: schedule.ID,
			UserID:     userUUID,
			TotalPrice: schedule.Class.Price,
			Status:     entity.ClassTransactionStatusPending,
			CreatedAt:  time.Now(),
		}

		created, err := s.transactionRepo.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		// 4. Decrement Remaining Capacity
		schedule.Remaining -= 1
		_, err = s.scheduleRepo.Update(ctx, tx, schedule)
		if err != nil {
			return err
		}

		transactionResult = dto.ClassTransactionResponse{
			TransactionID: created.ID,
			Status:        string(created.Status),
			TotalPrice:    created.TotalPrice,
		}

		return nil
	})

	if err != nil {
		return dto.ClassTransactionResponse{}, err
	}

	return transactionResult, nil
}

func (s *classTransactionService) Complete(ctx context.Context, transactionID string) error {
	transaction, err := s.transactionRepo.GetById(ctx, nil, transactionID)
	if err != nil {
		return err
	}

	transaction.Status = entity.ClassTransactionStatusSuccess
	_, err = s.transactionRepo.Update(ctx, nil, transaction)
	return err
}
