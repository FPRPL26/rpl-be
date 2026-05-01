package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"gorm.io/gorm"
)

type (
	PortofolioRepository interface {
		Create(ctx context.Context, tx *gorm.DB, portofolio entity.Portofolio) (entity.Portofolio, error)
		GetAll(ctx context.Context, tx *gorm.DB, tutorProfileID string) ([]entity.Portofolio, error)
		GetById(ctx context.Context, tx *gorm.DB, id string) (entity.Portofolio, error)
		Update(ctx context.Context, tx *gorm.DB, portofolio entity.Portofolio) (entity.Portofolio, error)
		Delete(ctx context.Context, tx *gorm.DB, portofolio entity.Portofolio) error
	}

	portofolioRepository struct {
		db *gorm.DB
	}
)

func NewPortofolioRepository(db *gorm.DB) PortofolioRepository {
	return &portofolioRepository{db}
}

func (r *portofolioRepository) Create(ctx context.Context, tx *gorm.DB, portofolio entity.Portofolio) (entity.Portofolio, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&portofolio).Error; err != nil {
		return portofolio, err
	}

	return portofolio, nil
}

func (r *portofolioRepository) GetAll(ctx context.Context, tx *gorm.DB, tutorProfileID string) ([]entity.Portofolio, error) {
	if tx == nil {
		tx = r.db
	}

	var portofolios []entity.Portofolio
	query := tx.WithContext(ctx).Model(&entity.Portofolio{})
	if tutorProfileID != "" {
		query = query.Where("tutor_profile_id = ?", tutorProfileID)
	}

	if err := query.Find(&portofolios).Error; err != nil {
		return nil, err
	}

	return portofolios, nil
}

func (r *portofolioRepository) GetById(ctx context.Context, tx *gorm.DB, id string) (entity.Portofolio, error) {
	if tx == nil {
		tx = r.db
	}

	var portofolio entity.Portofolio
	if err := tx.WithContext(ctx).Take(&portofolio, "id = ?", id).Error; err != nil {
		return entity.Portofolio{}, err
	}

	return portofolio, nil
}

func (r *portofolioRepository) Update(ctx context.Context, tx *gorm.DB, portofolio entity.Portofolio) (entity.Portofolio, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&portofolio).Error; err != nil {
		return entity.Portofolio{}, err
	}

	return portofolio, nil
}

func (r *portofolioRepository) Delete(ctx context.Context, tx *gorm.DB, portofolio entity.Portofolio) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&portofolio).Error; err != nil {
		return err
	}

	return nil
}
