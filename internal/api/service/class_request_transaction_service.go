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
	ClassRequestTransactionService interface {
		Create(ctx context.Context, userID string, req dto.CreateClassRequestTransactionRequest) (dto.ClassRequestTransactionResponse, error)
		GetAllByUserId(ctx context.Context, userID string, metaReq meta.Meta) ([]dto.ClassRequestTransactionListResponse, meta.Meta, error)
		GetById(ctx context.Context, userID, id string) (dto.ClassRequestTransactionResponse, error)
		UpdateStatus(ctx context.Context, userID, id string, req dto.UpdateClassRequestTransactionRequest) (dto.ClassRequestTransactionResponse, error)
		Complete(ctx context.Context, userID, transactionID string) error
		HandleMidtransCallback(ctx context.Context, payload map[string]interface{}) error
	}

	classRequestTransactionService struct {
		db               *gorm.DB
		transactionRepo  repository.ClassRequestTransactionRepository
		classRequestRepo repository.ClassRequestRepository
		tutorProfileRepo repository.TutorProfileRepository
		applicationRepo  repository.ClassRequestTutorApplicationRepository
		userRepo         repository.UserRepository
		midtransService  midtrans.MidtransService
	}
)

func NewClassRequestTransactionService(
	db *gorm.DB,
	transactionRepo repository.ClassRequestTransactionRepository,
	classRequestRepo repository.ClassRequestRepository,
	tutorProfileRepo repository.TutorProfileRepository,
	applicationRepo repository.ClassRequestTutorApplicationRepository,
	userRepo repository.UserRepository,
	midtransService midtrans.MidtransService,
) ClassRequestTransactionService {
	return &classRequestTransactionService{
		db:               db,
		transactionRepo:  transactionRepo,
		classRequestRepo: classRequestRepo,
		tutorProfileRepo: tutorProfileRepo,
		applicationRepo:  applicationRepo,
		userRepo:         userRepo,
		midtransService:  midtransService,
	}
}

func (s *classRequestTransactionService) Create(ctx context.Context, userID string, req dto.CreateClassRequestTransactionRequest) (dto.ClassRequestTransactionResponse, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return dto.ClassRequestTransactionResponse{}, myerror.New("invalid user ID", http.StatusBadRequest)
	}

	user, err := s.userRepo.GetById(ctx, nil, userID)
	if err != nil {
		return dto.ClassRequestTransactionResponse{}, err
	}

	requestID, err := uuid.Parse(req.RequestID)
	if err != nil {
		return dto.ClassRequestTransactionResponse{}, myerror.New("invalid request ID", http.StatusBadRequest)
	}

	if _, err := uuid.Parse(req.ApplicationID); err != nil {
		return dto.ClassRequestTransactionResponse{}, myerror.New("invalid application ID", http.StatusBadRequest)
	}

	application, err := s.applicationRepo.GetById(ctx, nil, req.ApplicationID, "ClassRequest", "TutorProfile")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTransactionResponse{}, myerror.New("application not found", http.StatusNotFound)
		}
		return dto.ClassRequestTransactionResponse{}, err
	}

	if application.RequestID.String() != req.RequestID {
		return dto.ClassRequestTransactionResponse{}, myerror.New("application request mismatch", http.StatusBadRequest)
	}

	if application.Status != entity.ClassRequestTutorApplicationStatusAccepted {
		return dto.ClassRequestTransactionResponse{}, myerror.New("application is not accepted", http.StatusBadRequest)
	}

	classRequest, err := s.classRequestRepo.GetById(ctx, nil, req.RequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTransactionResponse{}, myerror.New("class request not found", http.StatusNotFound)
		}
		return dto.ClassRequestTransactionResponse{}, err
	}

	if classRequest.UserID.String() != userID {
		return dto.ClassRequestTransactionResponse{}, myerror.New("unauthorized", http.StatusForbidden)
	}

	if application.ClassRequest.UserID.String() != userID {
		return dto.ClassRequestTransactionResponse{}, myerror.New("unauthorized", http.StatusForbidden)
	}

	tutorProfileID := application.TutorProfileID

	if application.ClassRequest.TutorProfileID != uuid.Nil && application.ClassRequest.TutorProfileID != tutorProfileID {
		return dto.ClassRequestTransactionResponse{}, myerror.New("accepted tutor does not match request record", http.StatusBadRequest)
	}

	if _, err := s.tutorProfileRepo.GetByID(ctx, tutorProfileID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTransactionResponse{}, myerror.New("tutor profile not found", http.StatusNotFound)
		}
		return dto.ClassRequestTransactionResponse{}, err
	}

	if existing, err := s.transactionRepo.GetByUserAndRequest(ctx, nil, userID, req.RequestID); err == nil && existing.ID != uuid.Nil {
		return dto.ClassRequestTransactionResponse{}, myerror.New("transaction already created for this request", http.StatusBadRequest)
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.ClassRequestTransactionResponse{}, err
	}

	price := req.Price
	if price <= 0 {
		price = classRequest.Price
	}

	transaction := entity.ClassRequestTransaction{
		ID:             uuid.New(),
		UserID:         userUUID,
		RequestID:      requestID,
		TutorProfileID: tutorProfileID,
		Status:         entity.ClassRequestTransactionStatusPending,
		CreatedAt:      time.Now(),
		Price:          price,
	}

	paymentURL, err := s.midtransService.CreateSnapTransaction(transaction.ID.String(), transaction.Price, user.Name, user.Email)
	if err != nil {
		return dto.ClassRequestTransactionResponse{}, err
	}
	transaction.PaymentURL = paymentURL

	created, err := s.transactionRepo.Create(ctx, nil, transaction)
	if err != nil {
		return dto.ClassRequestTransactionResponse{}, err
	}

	return dto.ClassRequestTransactionResponse{
		ID:             created.ID.String(),
		UserID:         created.UserID.String(),
		RequestID:      created.RequestID.String(),
		TutorProfileID: created.TutorProfileID.String(),
		Status:         string(created.Status),
		PaymentURL:     created.PaymentURL,
		Price:          created.Price,
		CreatedAt:      created.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *classRequestTransactionService) GetAllByUserId(ctx context.Context, userID string, metaReq meta.Meta) ([]dto.ClassRequestTransactionListResponse, meta.Meta, error) {
	transactions, metaRes, err := s.transactionRepo.GetAllByUserId(ctx, nil, userID, metaReq, "ClassRequest", "TutorProfile")
	if err != nil {
		return nil, metaRes, err
	}

	res := make([]dto.ClassRequestTransactionListResponse, 0, len(transactions))
	for _, t := range transactions {
		res = append(res, dto.ClassRequestTransactionListResponse{
			ID:             t.ID.String(),
			UserID:         t.UserID.String(),
			RequestID:      t.RequestID.String(),
			RequestName:    t.ClassRequest.Name,
			TutorProfileID: t.TutorProfileID.String(),
			TutorName:      t.TutorProfile.Name,
			Status:         string(t.Status),
			Price:          t.Price,
			CreatedAt:      t.CreatedAt.Format(time.RFC3339),
		})
	}

	return res, metaRes, nil
}

func (s *classRequestTransactionService) GetById(ctx context.Context, userID, id string) (dto.ClassRequestTransactionResponse, error) {
	transaction, err := s.transactionRepo.GetById(ctx, nil, id, "ClassRequest", "TutorProfile")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTransactionResponse{}, myerror.New("class request transaction not found", http.StatusNotFound)
		}
		return dto.ClassRequestTransactionResponse{}, err
	}

	if transaction.UserID.String() != userID {
		return dto.ClassRequestTransactionResponse{}, myerror.New("unauthorized", http.StatusForbidden)
	}

	return dto.ClassRequestTransactionResponse{
		ID:             transaction.ID.String(),
		UserID:         transaction.UserID.String(),
		RequestID:      transaction.RequestID.String(),
		RequestName:    transaction.ClassRequest.Name,
		TutorProfileID: transaction.TutorProfileID.String(),
		TutorName:      transaction.TutorProfile.Name,
		Status:         string(transaction.Status),
		PaymentURL:     transaction.PaymentURL,
		Price:          transaction.Price,
		CreatedAt:      transaction.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *classRequestTransactionService) UpdateStatus(ctx context.Context, userID, id string, req dto.UpdateClassRequestTransactionRequest) (dto.ClassRequestTransactionResponse, error) {
	transaction, err := s.transactionRepo.GetById(ctx, nil, id, "ClassRequest", "TutorProfile")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassRequestTransactionResponse{}, myerror.New("class request transaction not found", http.StatusNotFound)
		}
		return dto.ClassRequestTransactionResponse{}, err
	}

	if transaction.UserID.String() != userID {
		return dto.ClassRequestTransactionResponse{}, myerror.New("unauthorized", http.StatusForbidden)
	}

	status := entity.ClassRequestTransactionStatus(req.Status)
	switch status {
	case entity.ClassRequestTransactionStatusPending,
		entity.ClassRequestTransactionStatusAccepted,
		entity.ClassRequestTransactionStatusOnProgress,
		entity.ClassRequestTransactionStatusFinished,
		entity.ClassRequestTransactionStatusPaid,
		entity.ClassRequestTransactionStatusCancelled:
		transaction.Status = status
	default:
		return dto.ClassRequestTransactionResponse{}, myerror.New("invalid status", http.StatusBadRequest)
	}

	updated, err := s.transactionRepo.Update(ctx, nil, transaction)
	if err != nil {
		return dto.ClassRequestTransactionResponse{}, err
	}

	return dto.ClassRequestTransactionResponse{
		ID:             updated.ID.String(),
		UserID:         updated.UserID.String(),
		RequestID:      updated.RequestID.String(),
		RequestName:    updated.ClassRequest.Name,
		TutorProfileID: updated.TutorProfileID.String(),
		TutorName:      updated.TutorProfile.Name,
		Status:         string(updated.Status),
		PaymentURL:     updated.PaymentURL,
		Price:          updated.Price,
		CreatedAt:      updated.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *classRequestTransactionService) Complete(ctx context.Context, userID, transactionID string) error {
	transaction, err := s.transactionRepo.GetById(ctx, nil, transactionID)
	if err != nil {
		return err
	}

	if transaction.UserID.String() != userID {
		return myerror.New("unauthorized", http.StatusForbidden)
	}

	if transaction.Status != entity.ClassRequestTransactionStatusPaid {
		return myerror.New("transaction is not paid yet", http.StatusBadRequest)
	}

	transaction.Status = entity.ClassRequestTransactionStatusFinished
	_, err = s.transactionRepo.Update(ctx, nil, transaction)
	return err
}

func (s *classRequestTransactionService) HandleMidtransCallback(ctx context.Context, payload map[string]interface{}) error {
	orderID, ok := payload["order_id"].(string)
	if !ok {
		return myerror.New("missing order_id", http.StatusBadRequest)
	}

	statusCode, ok := payload["status_code"].(string)
	if !ok {
		return myerror.New("missing status_code", http.StatusBadRequest)
	}

	grossAmount, ok := payload["gross_amount"].(string)
	if !ok {
		return myerror.New("missing gross_amount", http.StatusBadRequest)
	}

	signatureKey, ok := payload["signature_key"].(string)
	if !ok {
		return myerror.New("missing signature_key", http.StatusBadRequest)
	}

	transactionStatus, ok := payload["transaction_status"].(string)
	if !ok {
		return myerror.New("missing transaction_status", http.StatusBadRequest)
	}

	if !s.midtransService.VerifySignatureKey(orderID, statusCode, grossAmount, signatureKey) {
		return myerror.New("invalid signature key", http.StatusBadRequest)
	}

	transaction, err := s.transactionRepo.GetById(ctx, nil, orderID, "ClassRequest", "TutorProfile")
	if err != nil {
		return err
	}

	if transactionStatus == "capture" || transactionStatus == "settlement" {
		transaction.Status = entity.ClassRequestTransactionStatusPaid
	} else if transactionStatus == "deny" || transactionStatus == "cancel" || transactionStatus == "expire" {
		transaction.Status = entity.ClassRequestTransactionStatusCancelled
	} else {
		return nil
	}

	_, err = s.transactionRepo.Update(ctx, nil, transaction)
	return err
}
