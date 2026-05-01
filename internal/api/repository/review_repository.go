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
		GetAverageRatingByClassId(ctx context.Context, tx *gorm.DB, classId string) (float64, error)
		GetLatestReviewsByClassId(ctx context.Context, tx *gorm.DB, classId string, limit int) ([]entity.Review, error)
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

func (r *reviewRepository) GetAverageRatingByClassId(ctx context.Context, tx *gorm.DB, classId string) (float64, error) {
	if tx == nil {
		tx = r.db
	}

	var result struct {
		Average float64
	}

	err := tx.WithContext(ctx).
		Table("reviews").
		Select("AVG(reviews.rating) as average").
		Joins("JOIN class_transactions ON class_transactions.id = reviews.transaction_id").
		Where("class_transactions.class_id = ? AND reviews.transaction_type = ?", classId, entity.ReviewTypeClass).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Average, nil
}

func (r *reviewRepository) GetLatestReviewsByClassId(ctx context.Context, tx *gorm.DB, classId string, limit int) ([]entity.Review, error) {
	if tx == nil {
		tx = r.db
	}

	var reviews []entity.Review
	err := tx.WithContext(ctx).
		Preload("User").
		Joins("JOIN class_transactions ON class_transactions.id = reviews.transaction_id").
		Where("class_transactions.class_id = ? AND reviews.transaction_type = ?", classId, entity.ReviewTypeClass).
		Order("reviews.created_at DESC").
		Limit(limit).
		Find(&reviews).Error

	if err != nil {
		return nil, err
	}

	return reviews, nil
}
