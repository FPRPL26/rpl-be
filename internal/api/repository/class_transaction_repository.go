package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"gorm.io/gorm"
)

type (
	ClassTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, transaction entity.ClassTransaction) (entity.ClassTransaction, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassTransaction, error)
		GetByUserAndSchedule(ctx context.Context, tx *gorm.DB, userID, scheduleID string) (entity.ClassTransaction, error)
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

func (r *classTransactionRepository) Update(ctx context.Context, tx *gorm.DB, transaction entity.ClassTransaction) (entity.ClassTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&transaction).Error; err != nil {
		return entity.ClassTransaction{}, err
	}

	return transaction, nil
}
