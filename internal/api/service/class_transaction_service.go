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
	"github.com/FPRPL26/rpl-be/internal/pkg/midtrans"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ClassTransactionService interface {
		Checkout(ctx context.Context, userID string, req dto.CheckoutClassRequest) (dto.ClassTransactionResponse, error)
		GetAllByUserId(ctx context.Context, userID string, metaReq meta.Meta) ([]dto.ClassTransactionListResponse, meta.Meta, error)
		Complete(ctx context.Context, transactionID string) error
		HandleMidtransCallback(ctx context.Context, payload map[string]interface{}) error
	}

	classTransactionService struct {
		db              *gorm.DB
		transactionRepo repository.ClassTransactionRepository
		scheduleRepo    repository.ScheduleRepository
		classRepo       repository.ClassRepository
		userRepo        repository.UserRepository
		midtransService midtrans.MidtransService
	}
)

func NewClassTransactionService(
	db *gorm.DB,
	transactionRepo repository.ClassTransactionRepository,
	scheduleRepo repository.ScheduleRepository,
	classRepo repository.ClassRepository,
	userRepo repository.UserRepository,
	midtransService midtrans.MidtransService,
) ClassTransactionService {
	return &classTransactionService{
		db:              db,
		transactionRepo: transactionRepo,
		scheduleRepo:    scheduleRepo,
		classRepo:       classRepo,
		userRepo:        userRepo,
		midtransService: midtransService,
	}
}

func (s *classTransactionService) Checkout(ctx context.Context, userID string, req dto.CheckoutClassRequest) (dto.ClassTransactionResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.ClassTransactionResponse{}, err
	}

	user, err := s.userRepo.GetById(ctx, nil, userID)
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

		// 4. Generate Midtrans Payment URL
		paymentURL, err := s.midtransService.CreateSnapTransaction(transaction.ID.String(), transaction.TotalPrice, user.Name, user.Email)
		if err != nil {
			return err
		}
		transaction.PaymentURL = paymentURL

		created, err := s.transactionRepo.Create(ctx, tx, transaction)
		if err != nil {
			return err
		}

		// 5. Decrement Remaining Capacity
		schedule.Remaining -= 1
		_, err = s.scheduleRepo.Update(ctx, tx, schedule)
		if err != nil {
			return err
		}

		transactionResult = dto.ClassTransactionResponse{
			TransactionID: created.ID,
			Status:        string(created.Status),
			TotalPrice:    created.TotalPrice,
			PaymentURL:    created.PaymentURL,
		}

		return nil
	})

	if err != nil {
		return dto.ClassTransactionResponse{}, err
	}

	return transactionResult, nil
}

func (s *classTransactionService) GetAllByUserId(ctx context.Context, userID string, metaReq meta.Meta) ([]dto.ClassTransactionListResponse, meta.Meta, error) {
	transactions, metaRes, err := s.transactionRepo.GetAllByUserId(ctx, nil, userID, metaReq, "Class")
	if err != nil {
		return nil, metaRes, err
	}

	res := make([]dto.ClassTransactionListResponse, 0, len(transactions))
	for _, t := range transactions {
		res = append(res, dto.ClassTransactionListResponse{
			ID:         t.ID,
			ClassID:    t.ClassID,
			ClassName:  t.Class.Name,
			ScheduleID: t.ScheduleID,
			TotalPrice: t.TotalPrice,
			Status:     string(t.Status),
			CreatedAt:  t.CreatedAt.Format(time.RFC3339),
		})
	}

	return res, metaRes, nil
}

func (s *classTransactionService) Complete(ctx context.Context, transactionID string) error {
	transaction, err := s.transactionRepo.GetById(ctx, nil, transactionID)
	if err != nil {
		return err
	}

	if transaction.Status != entity.ClassTransactionStatusPaid {
		return myerror.New("transaction is not paid yet", http.StatusBadRequest)
	}

	transaction.Status = entity.ClassTransactionStatusSuccess
	_, err = s.transactionRepo.Update(ctx, nil, transaction)
	return err
}

func (s *classTransactionService) HandleMidtransCallback(ctx context.Context, payload map[string]interface{}) error {
	orderID := payload["order_id"].(string)
	statusCode := payload["status_code"].(string)
	grossAmount := payload["gross_amount"].(string)
	signatureKey := payload["signature_key"].(string)
	transactionStatus := payload["transaction_status"].(string)

	if !s.midtransService.VerifySignatureKey(orderID, statusCode, grossAmount, signatureKey) {
		return errors.New("invalid signature key")
	}

	transaction, err := s.transactionRepo.GetById(ctx, nil, orderID)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if transactionStatus == "capture" || transactionStatus == "settlement" {
			transaction.Status = entity.ClassTransactionStatusPaid
		} else if transactionStatus == "deny" || transactionStatus == "cancel" || transactionStatus == "expire" {
			// If cancelled or expired, restore capacity ONLY if it was previously PENDING
			if transaction.Status == entity.ClassTransactionStatusPending {
				transaction.Status = entity.ClassTransactionStatusCancelled

				schedule, err := s.scheduleRepo.GetById(ctx, tx.Set("gorm:query_option", "FOR UPDATE"), transaction.ScheduleID.String())
				if err == nil {
					schedule.Remaining += 1
					s.scheduleRepo.Update(ctx, tx, schedule)
				}
			} else {
				// Just update status if it wasn't pending (e.g. already cancelled)
				transaction.Status = entity.ClassTransactionStatusCancelled
			}
		}

		_, err = s.transactionRepo.Update(ctx, tx, transaction)
		return err
	})
}
