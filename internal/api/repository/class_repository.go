package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	ClassRepository interface {
		Create(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error)
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.Class, meta.Meta, error)
		GetById(ctx context.Context, tx *gorm.DB, classId string, preloads ...string) (entity.Class, error)
		Update(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error)
		Delete(ctx context.Context, tx *gorm.DB, class entity.Class) error
	}

	classRepository struct {
		db *gorm.DB
	}
)

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db}
}

func (r *classRepository) Create(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&class).Error; err != nil {
		return class, err
	}

	return class, nil
}

func (r *classRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.Class, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var classes []entity.Class

	tx = tx.WithContext(ctx).Model(entity.Class{})

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.Class{})).Find(&classes).Error; err != nil {
		return nil, metaReq, err
	}

	return classes, metaReq, nil
}

func (r *classRepository) GetById(ctx context.Context, tx *gorm.DB, classId string, preloads ...string) (entity.Class, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var class entity.Class
	if err := tx.WithContext(ctx).Take(&class, "id = ?", classId).Error; err != nil {
		return entity.Class{}, err
	}

	return class, nil
}

func (r *classRepository) Update(ctx context.Context, tx *gorm.DB, class entity.Class) (entity.Class, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&class).Error; err != nil {
		return entity.Class{}, err
	}

	return class, nil
}

func (r *classRepository) Delete(ctx context.Context, tx *gorm.DB, class entity.Class) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&class).Error; err != nil {
		return err
	}

	return nil
}
