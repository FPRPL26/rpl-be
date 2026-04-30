package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"gorm.io/gorm"
)

type (
	ClassRequestTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, transaction entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequestTransaction, error)
		Update(ctx context.Context, tx *gorm.DB, transaction entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error)
	}

	classRequestTransactionRepository struct {
		db *gorm.DB
	}
)

func NewClassRequestTransactionRepository(db *gorm.DB) ClassRequestTransactionRepository {
	return &classRequestTransactionRepository{db}
}

func (r *classRequestTransactionRepository) Create(ctx context.Context, tx *gorm.DB, transaction entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *classRequestTransactionRepository) GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequestTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var transaction entity.ClassRequestTransaction
	if err := tx.WithContext(ctx).Take(&transaction, "id = ?", id).Error; err != nil {
		return entity.ClassRequestTransaction{}, err
	}

	return transaction, nil
}

func (r *classRequestTransactionRepository) Update(ctx context.Context, tx *gorm.DB, transaction entity.ClassRequestTransaction) (entity.ClassRequestTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&transaction).Error; err != nil {
		return entity.ClassRequestTransaction{}, err
	}

	return transaction, nil
}
