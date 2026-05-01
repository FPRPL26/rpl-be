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
	BarterSkillService interface {
		CreateOffer(ctx context.Context, tutorID string, req dto.CreateBarterOfferRequest) (dto.CreateBarterOfferResponse, error)
		RequestOffer(ctx context.Context, tutorID string, barterID string) (dto.RequestBarterOfferResponse, error)
		ApproveRequest(ctx context.Context, ownerID string, barterID string, req dto.ApproveBarterRequestRequest) error
		GetAllOffers(ctx context.Context, metaReq meta.Meta) ([]dto.BarterOfferResponse, meta.Meta, error)
		GetOfferById(ctx context.Context, barterID string) (dto.BarterOfferResponse, error)
		GetMyOffers(ctx context.Context, tutorID string, metaReq meta.Meta) ([]dto.BarterOfferResponse, meta.Meta, error)
		GetMyRequests(ctx context.Context, tutorID string, metaReq meta.Meta) ([]dto.BarterRequestResponse, meta.Meta, error)
		GetIncomingRequests(ctx context.Context, tutorID string, metaReq meta.Meta) ([]dto.BarterIncomingRequestResponse, meta.Meta, error)
	}

	barterSkillService struct {
		db              *gorm.DB
		barterRepo      repository.BarterSkillRepository
		transactionRepo repository.BarterSkillTransactionRepository
		tutorRepo       repository.TutorProfileRepository
	}
)

func NewBarterSkillService(
	db *gorm.DB,
	barterRepo repository.BarterSkillRepository,
	transactionRepo repository.BarterSkillTransactionRepository,
	tutorRepo repository.TutorProfileRepository,
) BarterSkillService {
	return &barterSkillService{
		db:              db,
		barterRepo:      barterRepo,
		transactionRepo: transactionRepo,
		tutorRepo:       tutorRepo,
	}
}

func (s *barterSkillService) CreateOffer(ctx context.Context, tutorID string, req dto.CreateBarterOfferRequest) (dto.CreateBarterOfferResponse, error) {
	tutorUUID, err := uuid.Parse(tutorID)
	if err != nil {
		return dto.CreateBarterOfferResponse{}, err
	}

	// Verify tutor exists
	_, err = s.tutorRepo.GetByID(ctx, tutorUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.CreateBarterOfferResponse{}, myerror.New("only verified tutors can create barter offers", http.StatusForbidden)
		}
		return dto.CreateBarterOfferResponse{}, err
	}

	barter := entity.BarterSkill{
		ID:             uuid.New(),
		TutorProfileID: tutorUUID,
		RequestSkills:  req.RequestSkillID,
		OfferedSkills:  req.OfferedSkillID,
		Name:           req.Name,
		Description:    req.Description,
		ChatWA:         &req.ChatWA,
		Accepted:       false,
	}

	created, err := s.barterRepo.Create(ctx, nil, barter)
	if err != nil {
		return dto.CreateBarterOfferResponse{}, err
	}

	return dto.CreateBarterOfferResponse{
		BarterID: created.ID,
	}, nil
}

func (s *barterSkillService) RequestOffer(ctx context.Context, tutorID string, barterID string) (dto.RequestBarterOfferResponse, error) {
	tutorUUID, err := uuid.Parse(tutorID)
	if err != nil {
		return dto.RequestBarterOfferResponse{}, err
	}

	// Verify requesting tutor exists
	_, err = s.tutorRepo.GetByID(ctx, tutorUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.RequestBarterOfferResponse{}, myerror.New("only verified tutors can request barter offers", http.StatusForbidden)
		}
		return dto.RequestBarterOfferResponse{}, err
	}

	// 1. Get Barter Offer
	barter, err := s.barterRepo.GetById(ctx, nil, barterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.RequestBarterOfferResponse{}, myerror.New("barter offer not found", http.StatusNotFound)
		}
		return dto.RequestBarterOfferResponse{}, err
	}

	// 2. Validate
	if barter.Accepted {
		return dto.RequestBarterOfferResponse{}, myerror.New("barter offer already closed", http.StatusBadRequest)
	}
	if barter.TutorProfileID == tutorUUID {
		return dto.RequestBarterOfferResponse{}, myerror.New("you cannot request your own barter offer", http.StatusBadRequest)
	}

	// 3. Check if already requested
	existing, err := s.transactionRepo.GetByBarterIdAndMentor2(ctx, nil, barterID, tutorID)
	if err == nil && existing.ID != uuid.Nil {
		return dto.RequestBarterOfferResponse{}, myerror.New("you have already requested this barter offer", http.StatusBadRequest)
	}

	// 4. Create PENDING Transaction
	transaction := entity.BarterSkillTransaction{
		ID:               uuid.New(),
		BarterSkillID:    barter.ID,
		MentorProfileID1: barter.TutorProfileID,
		MentorProfileID2: tutorUUID,
		Status:           entity.BarterSkillTransactionStatusPending,
		CreatedAt:        time.Now(),
	}

	createdTrans, err := s.transactionRepo.Create(ctx, nil, transaction)
	if err != nil {
		return dto.RequestBarterOfferResponse{}, err
	}

	return dto.RequestBarterOfferResponse{
		TransactionID: createdTrans.ID,
	}, nil
}

func (s *barterSkillService) ApproveRequest(ctx context.Context, ownerID string, barterID string, req dto.ApproveBarterRequestRequest) error {
	ownerUUID, err := uuid.Parse(ownerID)
	if err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Get Barter Offer and lock
		barter, err := s.barterRepo.GetById(ctx, tx.Set("gorm:query_option", "FOR UPDATE"), barterID)
		if err != nil {
			return err
		}

		// 2. Validate Ownership
		if barter.TutorProfileID != ownerUUID {
			return myerror.New("you are not authorized to approve requests for this offer", http.StatusForbidden)
		}
		if barter.Accepted {
			return myerror.New("barter offer already closed", http.StatusBadRequest)
		}

		// 3. Get the specific request
		targetTrans, err := s.transactionRepo.GetById(ctx, tx.Set("gorm:query_option", "FOR UPDATE"), req.TransactionID)
		if err != nil {
			return err
		}
		if targetTrans.BarterSkillID != barter.ID {
			return myerror.New("invalid transaction for this barter offer", http.StatusBadRequest)
		}
		if targetTrans.Status != entity.BarterSkillTransactionStatusPending {
			return myerror.New("transaction is no longer pending", http.StatusBadRequest)
		}

		// 4. Update Barter Offer to Accepted (Closed)
		barter.Accepted = true
		if _, err := s.barterRepo.Update(ctx, tx, barter); err != nil {
			return err
		}

		// 5. Update Target Transaction to ACCEPTED
		targetTrans.Status = entity.BarterSkillTransactionStatusAccepted
		if _, err := s.transactionRepo.Update(ctx, tx, targetTrans); err != nil {
			return err
		}

		// 6. Reject all other PENDING requests
		if err := s.transactionRepo.RejectAllOtherRequests(ctx, tx, barterID, targetTrans.ID.String()); err != nil {
			return err
		}

		return nil
	})
}

func (s *barterSkillService) GetAllOffers(ctx context.Context, metaReq meta.Meta) ([]dto.BarterOfferResponse, meta.Meta, error) {
	barters, metaRes, err := s.barterRepo.GetAll(ctx, nil, metaReq, "TutorProfile", "RequestSkill", "OfferedSkill")
	if err != nil {
		return nil, metaRes, err
	}

	res := make([]dto.BarterOfferResponse, 0, len(barters))
	for _, b := range barters {
		chatWA := ""
		if b.ChatWA != nil {
			chatWA = *b.ChatWA
		}
		res = append(res, dto.BarterOfferResponse{
			ID:               b.ID,
			TutorID:          b.TutorProfileID,
			TutorName:        b.TutorProfile.Name,
			RequestSkillID:   b.RequestSkills,
			RequestSkillName: b.RequestSkill.Name,
			OfferedSkillID:   b.OfferedSkills,
			OfferedSkillName: b.OfferedSkill.Name,
			Name:             b.Name,
			Description:      b.Description,
			ChatWA:           chatWA,
			Accepted:         b.Accepted,
		})
	}

	return res, metaRes, nil
}

func (s *barterSkillService) GetMyOffers(ctx context.Context, tutorID string, metaReq meta.Meta) ([]dto.BarterOfferResponse, meta.Meta, error) {
	tutorUUID, err := uuid.Parse(tutorID)
	if err != nil {
		return nil, metaReq, err
	}

	tx := s.db.Where("tutor_profile_id = ?", tutorUUID)
	barters, metaRes, err := s.barterRepo.GetAll(ctx, tx, metaReq, "TutorProfile", "RequestSkill", "OfferedSkill")
	if err != nil {
		return nil, metaRes, err
	}

	res := make([]dto.BarterOfferResponse, 0, len(barters))
	for _, b := range barters {
		chatWA := ""
		if b.ChatWA != nil {
			chatWA = *b.ChatWA
		}
		res = append(res, dto.BarterOfferResponse{
			ID:               b.ID,
			TutorID:          b.TutorProfileID,
			TutorName:        b.TutorProfile.Name,
			RequestSkillID:   b.RequestSkills,
			RequestSkillName: b.RequestSkill.Name,
			OfferedSkillID:   b.OfferedSkills,
			OfferedSkillName: b.OfferedSkill.Name,
			Name:             b.Name,
			Description:      b.Description,
			ChatWA:           chatWA,
			Accepted:         b.Accepted,
		})
	}

	return res, metaRes, nil
}

func (s *barterSkillService) GetMyRequests(ctx context.Context, tutorID string, metaReq meta.Meta) ([]dto.BarterRequestResponse, meta.Meta, error) {
	_, err := uuid.Parse(tutorID)
	if err != nil {
		return nil, metaReq, err
	}

	transactions, metaRes, err := s.transactionRepo.GetAllByMentor2Id(ctx, nil, tutorID, metaReq, "BarterSkill", "BarterSkill.TutorProfile", "BarterSkill.RequestSkill", "BarterSkill.OfferedSkill")
	if err != nil {
		return nil, metaRes, err
	}

	res := make([]dto.BarterRequestResponse, 0, len(transactions))
	for _, t := range transactions {
		res = append(res, dto.BarterRequestResponse{
			TransactionID:    t.ID,
			Status:           string(t.Status),
			BarterID:         t.BarterSkillID,
			TutorID:          t.BarterSkill.TutorProfileID,
			TutorName:        t.BarterSkill.TutorProfile.Name,
			RequestSkillName: t.BarterSkill.RequestSkill.Name,
			OfferedSkillName: t.BarterSkill.OfferedSkill.Name,
			Name:             t.BarterSkill.Name,
			Description:      t.BarterSkill.Description,
		})
	}

	return res, metaRes, nil
}

func (s *barterSkillService) GetIncomingRequests(ctx context.Context, tutorID string, metaReq meta.Meta) ([]dto.BarterIncomingRequestResponse, meta.Meta, error) {
	_, err := uuid.Parse(tutorID)
	if err != nil {
		return nil, metaReq, err
	}

	transactions, metaRes, err := s.transactionRepo.GetAllByMentor1Id(ctx, nil, tutorID, metaReq, "MentorProfile2", "BarterSkill", "BarterSkill.RequestSkill", "BarterSkill.OfferedSkill")
	if err != nil {
		return nil, metaRes, err
	}

	res := make([]dto.BarterIncomingRequestResponse, 0, len(transactions))
	for _, t := range transactions {
		res = append(res, dto.BarterIncomingRequestResponse{
			TransactionID:    t.ID,
			Status:           string(t.Status),
			BarterID:         t.BarterSkillID,
			RequesterID:      t.MentorProfileID2,
			RequesterName:    t.MentorProfile2.Name,
			RequestSkillName: t.BarterSkill.RequestSkill.Name,
			OfferedSkillName: t.BarterSkill.OfferedSkill.Name,
			Name:             t.BarterSkill.Name,
			Description:      t.BarterSkill.Description,
		})
	}

	return res, metaRes, nil
}

func (s *barterSkillService) GetOfferById(ctx context.Context, barterID string) (dto.BarterOfferResponse, error) {
	b, err := s.barterRepo.GetById(ctx, nil, barterID, "TutorProfile", "RequestSkill", "OfferedSkill")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.BarterOfferResponse{}, myerror.New("barter offer not found", http.StatusNotFound)
		}
		return dto.BarterOfferResponse{}, err
	}

	chatWA := ""
	if b.ChatWA != nil {
		chatWA = *b.ChatWA
	}

	return dto.BarterOfferResponse{
		ID:               b.ID,
		TutorID:          b.TutorProfileID,
		TutorName:        b.TutorProfile.Name,
		RequestSkillID:   b.RequestSkills,
		RequestSkillName: b.RequestSkill.Name,
		OfferedSkillID:   b.OfferedSkills,
		OfferedSkillName: b.OfferedSkill.Name,
		Name:             b.Name,
		Description:      b.Description,
		ChatWA:           chatWA,
		Accepted:         b.Accepted,
	}, nil
}
