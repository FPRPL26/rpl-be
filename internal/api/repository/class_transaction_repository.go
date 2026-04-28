package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	ClassTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, transaction entity.ClassTransaction) (entity.ClassTransaction, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassTransaction, error)
		GetByUserAndSchedule(ctx context.Context, tx *gorm.DB, userID, scheduleID string) (entity.ClassTransaction, error)
		GetAllByUserId(ctx context.Context, tx *gorm.DB, userID string, metaReq meta.Meta, preloads ...string) ([]entity.ClassTransaction, meta.Meta, error)
		Update(ctx context.Context, tx *gorm.DB, transaction entity.ClassTransaction) (entity.ClassTransaction, error)
	}

	classTransactionRepository struct {
		db *gorm.DB
	}
)

func NewClassTransactionRepository(db *gorm.DB) ClassTransactionRepository {
	return &classTransactionRepository{db}
}

func (r *classTransactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction entity.ClassTransaction) (entity.ClassTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *classTransactionRepository) GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var transaction entity.ClassTransaction
	if err := tx.WithContext(ctx).Take(&transaction, "id = ?", id).Error; err != nil {
		return entity.ClassTransaction{}, err
	}

	return transaction, nil
}

func (r *classTransactionRepository) GetByUserAndSchedule(ctx context.Context, tx *gorm.DB, userID, scheduleID string) (entity.ClassTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	var transaction entity.ClassTransaction
	if err := tx.WithContext(ctx).Where("user_id = ? AND schedule_id = ? AND status != ?", userID, scheduleID, entity.ClassTransactionStatusCancelled).First(&transaction).Error; err != nil {
		return entity.ClassTransaction{}, err
	}

	return transaction, nil
}

func (r *classTransactionRepository) GetAllByUserId(ctx context.Context, tx *gorm.DB, userID string, metaReq meta.Meta, preloads ...string) ([]entity.ClassTransaction, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var transactions []entity.ClassTransaction
	tx = tx.WithContext(ctx).Model(entity.ClassTransaction{}).Where("user_id = ?", userID)

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.ClassTransaction{})).Find(&transactions).Error; err != nil {
		return nil, metaReq, err
	}

	return transactions, metaReq, nil
}

func (r *classTransactionRepository) Update(ctx context.Context, tx *gorm.DB, transaction entity.ClassTransaction) (entity.ClassTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&transaction).Error; err != nil {
		return entity.ClassTransaction{}, err
	}

	return transaction, nil
}
