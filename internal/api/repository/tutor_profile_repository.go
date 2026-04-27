package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	TutorProfileRepository interface {
		Create(ctx context.Context, profile *entity.TutorProfile) error
		GetByID(ctx context.Context, id uuid.UUID) (*entity.TutorProfile, error)
		Update(ctx context.Context, profile *entity.TutorProfile) error
		Delete(ctx context.Context, id uuid.UUID) error
		List(ctx context.Context, limit, offset int) ([]entity.TutorProfile, error)
	}

	tutorProfileRepository struct {
		db *gorm.DB
	}
)

func NewTutorProfileRepository(db *gorm.DB) TutorProfileRepository {
	return &tutorProfileRepository{db: db}
}

func (r *tutorProfileRepository) Create(ctx context.Context, profile *entity.TutorProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *tutorProfileRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TutorProfile, error) {
	var profile entity.TutorProfile
	err := r.db.WithContext(ctx).Preload("User").First(&profile, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *tutorProfileRepository) Update(ctx context.Context, profile *entity.TutorProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *tutorProfileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.TutorProfile{}, "id = ?", id).Error
}

func (r *tutorProfileRepository) List(ctx context.Context, limit, offset int) ([]entity.TutorProfile, error) {
	var profiles []entity.TutorProfile
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&profiles).Error
	return profiles, err
}
