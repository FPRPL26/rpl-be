package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	ClassRequestRepository interface {
		Create(ctx context.Context, tx *gorm.DB, cr entity.ClassRequest, preloads ...string) (entity.ClassRequest, error)
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.ClassRequest, meta.Meta, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequest, error)
		Update(ctx context.Context, tx *gorm.DB, cr entity.ClassRequest, preloads ...string) (entity.ClassRequest, error)
		Delete(ctx context.Context, tx *gorm.DB, cr entity.ClassRequest) error
	}

	classRequestRepository struct {
		db *gorm.DB
	}
)

func NewClassRequestRepository(db *gorm.DB) ClassRequestRepository {
	return &classRequestRepository{db}
}

func (r *classRequestRepository) Create(ctx context.Context, tx *gorm.DB, cr entity.ClassRequest, preloads ...string) (entity.ClassRequest, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Create(&cr).Error; err != nil {
		return cr, err
	}

	return cr, nil
}

func (r *classRequestRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.ClassRequest, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var items []entity.ClassRequest

	tx = tx.WithContext(ctx).Model(entity.ClassRequest{})

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.ClassRequest{})).Find(&items).Error; err != nil {
		return nil, metaReq, err
	}

	return items, metaReq, nil
}

func (r *classRequestRepository) GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequest, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var cr entity.ClassRequest
	if err := tx.WithContext(ctx).Take(&cr, "id = ?", id).Error; err != nil {
		return entity.ClassRequest{}, err
	}

	return cr, nil
}

func (r *classRequestRepository) Update(ctx context.Context, tx *gorm.DB, cr entity.ClassRequest, preloads ...string) (entity.ClassRequest, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	if err := tx.WithContext(ctx).Save(&cr).Error; err != nil {
		return entity.ClassRequest{}, err
	}

	return cr, nil
}

func (r *classRequestRepository) Delete(ctx context.Context, tx *gorm.DB, cr entity.ClassRequest) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&cr).Error; err != nil {
		return err
	}

	return nil
}
