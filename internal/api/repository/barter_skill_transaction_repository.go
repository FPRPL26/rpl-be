package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	BarterSkillTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, transaction entity.BarterSkillTransaction) (entity.BarterSkillTransaction, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.BarterSkillTransaction, error)
		GetByBarterIdAndMentor2(ctx context.Context, tx *gorm.DB, barterId string, mentor2Id string) (entity.BarterSkillTransaction, error)
		GetAllByMentor1Id(ctx context.Context, tx *gorm.DB, mentor1Id string, metaReq meta.Meta, preloads ...string) ([]entity.BarterSkillTransaction, meta.Meta, error)
		GetAllByMentor2Id(ctx context.Context, tx *gorm.DB, mentor2Id string, metaReq meta.Meta, preloads ...string) ([]entity.BarterSkillTransaction, meta.Meta, error)
		Update(ctx context.Context, tx *gorm.DB, transaction entity.BarterSkillTransaction) (entity.BarterSkillTransaction, error)
		RejectAllOtherRequests(ctx context.Context, tx *gorm.DB, barterId string, acceptedTransId string) error
		GetAllRequestOffer(ctx context.Context, tx *gorm.DB, barterId string) ([]entity.BarterSkillTransaction, error)
	}

	barterSkillTransactionRepository struct {
		db *gorm.DB
	}
)

func NewBarterSkillTransactionRepository(db *gorm.DB) BarterSkillTransactionRepository {
	return &barterSkillTransactionRepository{db}
}

func (r *barterSkillTransactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction entity.BarterSkillTransaction) (entity.BarterSkillTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *barterSkillTransactionRepository) GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.BarterSkillTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var transaction entity.BarterSkillTransaction
	if err := tx.WithContext(ctx).Take(&transaction, "id = ?", id).Error; err != nil {
		return entity.BarterSkillTransaction{}, err
	}

	return transaction, nil
}

func (r *barterSkillTransactionRepository) GetByBarterIdAndMentor2(ctx context.Context, tx *gorm.DB, barterId string, mentor2Id string) (entity.BarterSkillTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	var transaction entity.BarterSkillTransaction
	if err := tx.WithContext(ctx).Where("barter_skill_id = ? AND mentor_profile_id2 = ?", barterId, mentor2Id).First(&transaction).Error; err != nil {
		return entity.BarterSkillTransaction{}, err
	}

	return transaction, nil
}

func (r *barterSkillTransactionRepository) GetAllByMentor1Id(ctx context.Context, tx *gorm.DB, mentor1Id string, metaReq meta.Meta, preloads ...string) ([]entity.BarterSkillTransaction, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var transactions []entity.BarterSkillTransaction

	if metaReq.SortBy == "" {
		metaReq.SortBy = "created_at"
		metaReq.Sort = "desc"
	}

	tx = tx.WithContext(ctx).Model(entity.BarterSkillTransaction{}).Where("mentor_profile_id1 = ?", mentor1Id)

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.BarterSkillTransaction{}),
		AddCustomField("barter_id", "barter_skill_id = ?"),
	).Find(&transactions).Error; err != nil {
		return nil, metaReq, err
	}

	return transactions, metaReq, nil
}

func (r *barterSkillTransactionRepository) GetAllByMentor2Id(ctx context.Context, tx *gorm.DB, mentor2Id string, metaReq meta.Meta, preloads ...string) ([]entity.BarterSkillTransaction, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var transactions []entity.BarterSkillTransaction

	if metaReq.SortBy == "" {
		metaReq.SortBy = "created_at"
		metaReq.Sort = "desc"
	}

	tx = tx.WithContext(ctx).Model(entity.BarterSkillTransaction{}).Where("mentor_profile_id2 = ?", mentor2Id)

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.BarterSkillTransaction{}),
		AddCustomField("barter_id", "barter_skill_id = ?"),
	).Find(&transactions).Error; err != nil {
		return nil, metaReq, err
	}

	return transactions, metaReq, nil
}

func (r *barterSkillTransactionRepository) Update(ctx context.Context, tx *gorm.DB, transaction entity.BarterSkillTransaction) (entity.BarterSkillTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&transaction).Error; err != nil {
		return entity.BarterSkillTransaction{}, err
	}

	return transaction, nil
}

func (r *barterSkillTransactionRepository) RejectAllOtherRequests(ctx context.Context, tx *gorm.DB, barterId string, acceptedTransId string) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Model(&entity.BarterSkillTransaction{}).
		Where("barter_skill_id = ? AND id != ? AND status = ?", barterId, acceptedTransId, entity.BarterSkillTransactionStatusPending).
		Update("status", entity.BarterSkillTransactionStatusRejected).Error
}

func (r *barterSkillTransactionRepository) GetAllRequestOffer(ctx context.Context, tx *gorm.DB, barterId string) ([]entity.BarterSkillTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	var transactions []entity.BarterSkillTransaction
	if err := tx.WithContext(ctx).Where("barter_skill_id = ? AND status = ?", barterId, entity.BarterSkillTransactionStatusPending).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
