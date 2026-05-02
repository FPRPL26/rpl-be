package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	ClassRequestTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, crt entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequestTransaction, error)
		GetByUserAndRequest(ctx context.Context, tx *gorm.DB, userID, requestID string) (entity.ClassRequestTransaction, error)
		GetAllByUserId(ctx context.Context, tx *gorm.DB, userID string, metaReq meta.Meta, preloads ...string) ([]entity.ClassRequestTransaction, meta.Meta, error)
		Update(ctx context.Context, tx *gorm.DB, crt entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error)
	}

	classRequestTransactionRepository struct {
		db *gorm.DB
	}
)

func NewClassRequestTransactionRepository(db *gorm.DB) ClassRequestTransactionRepository {
	return &classRequestTransactionRepository{db}
}

func (r *classRequestTransactionRepository) Create(ctx context.Context, tx *gorm.DB, crt entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&crt).Error; err != nil {
		return crt, err
	}

	return crt, nil
}

func (r *classRequestTransactionRepository) GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequestTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var crt entity.ClassRequestTransaction
	if err := tx.WithContext(ctx).Take(&crt, "id = ?", id).Error; err != nil {
		return entity.ClassRequestTransaction{}, err
	}

	return crt, nil
}

func (r *classRequestTransactionRepository) GetByUserAndRequest(ctx context.Context, tx *gorm.DB, userID, requestID string) (entity.ClassRequestTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	var crt entity.ClassRequestTransaction
	if err := tx.WithContext(ctx).Where("user_id = ? AND request_id = ?", userID, requestID).First(&crt).Error; err != nil {
		return entity.ClassRequestTransaction{}, err
	}

	return crt, nil
}

func (r *classRequestTransactionRepository) GetAllByUserId(ctx context.Context, tx *gorm.DB, userID string, metaReq meta.Meta, preloads ...string) ([]entity.ClassRequestTransaction, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var crts []entity.ClassRequestTransaction
	tx = tx.WithContext(ctx).Model(entity.ClassRequestTransaction{}).Where("user_id = ?", userID)

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.ClassRequestTransaction{})).Find(&crts).Error; err != nil {
		return nil, metaReq, err
	}

	return crts, metaReq, nil
}

func (r *classRequestTransactionRepository) Update(ctx context.Context, tx *gorm.DB, crt entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&crt).Error; err != nil {
		return entity.ClassRequestTransaction{}, err
	}

	return crt, nil
}
