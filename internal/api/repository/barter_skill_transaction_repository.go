package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"gorm.io/gorm"
)

type (
	BarterSkillTransactionRepository interface {
		Create(ctx context.Context, tx *gorm.DB, transaction entity.BarterSkillTransaction) (entity.BarterSkillTransaction, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.BarterSkillTransaction, error)
		Update(ctx context.Context, tx *gorm.DB, transaction entity.BarterSkillTransaction) (entity.BarterSkillTransaction, error)
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

func (r *barterSkillTransactionRepository) Update(ctx context.Context, tx *gorm.DB, transaction entity.BarterSkillTransaction) (entity.BarterSkillTransaction, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&transaction).Error; err != nil {
		return entity.BarterSkillTransaction{}, err
	}

	return transaction, nil
}
