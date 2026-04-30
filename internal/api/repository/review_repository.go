package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"gorm.io/gorm"
)

type (
	ReviewRepository interface {
		Create(ctx context.Context, tx *gorm.DB, review entity.Review) (entity.Review, error)
		GetByTransactionId(ctx context.Context, tx *gorm.DB, transactionID string) (entity.Review, error)
	}

	reviewRepository struct {
		db *gorm.DB
	}
)

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db}
}

func (r *reviewRepository) Create(ctx context.Context, tx *gorm.DB, review entity.Review) (entity.Review, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&review).Error; err != nil {
		return review, err
	}

	return review, nil
}

func (r *reviewRepository) GetByTransactionId(ctx context.Context, tx *gorm.DB, transactionID string) (entity.Review, error) {
	if tx == nil {
		tx = r.db
	}

	var review entity.Review
	if err := tx.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&review).Error; err != nil {
		return entity.Review{}, err
	}

	return review, nil
}
