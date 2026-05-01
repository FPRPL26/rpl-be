package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"gorm.io/gorm"
)

type (
	ClassRequestTutorApplicationRepository interface {
		Create(ctx context.Context, tx *gorm.DB, application entity.ClassRequestTutorApplication) (entity.ClassRequestTutorApplication, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequestTutorApplication, error)
		GetByTutorAndRequest(ctx context.Context, tx *gorm.DB, tutorProfileID, requestID string) (entity.ClassRequestTutorApplication, error)
		GetAllByRequestId(ctx context.Context, tx *gorm.DB, requestID string, preloads ...string) ([]entity.ClassRequestTutorApplication, error)
		GetAllByTutorProfileId(ctx context.Context, tx *gorm.DB, tutorProfileID string, preloads ...string) ([]entity.ClassRequestTutorApplication, error)
		Update(ctx context.Context, tx *gorm.DB, application entity.ClassRequestTutorApplication) (entity.ClassRequestTutorApplication, error)
	}

	classRequestTutorApplicationRepository struct {
		db *gorm.DB
	}
)

func NewClassRequestTutorApplicationRepository(db *gorm.DB) ClassRequestTutorApplicationRepository {
	return &classRequestTutorApplicationRepository{db}
}

func (r *classRequestTutorApplicationRepository) Create(ctx context.Context, tx *gorm.DB, application entity.ClassRequestTutorApplication) (entity.ClassRequestTutorApplication, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&application).Error; err != nil {
		return application, err
	}

	return application, nil
}

func (r *classRequestTutorApplicationRepository) GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.ClassRequestTutorApplication, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var application entity.ClassRequestTutorApplication
	if err := tx.WithContext(ctx).Take(&application, "id = ?", id).Error; err != nil {
		return entity.ClassRequestTutorApplication{}, err
	}

	return application, nil
}

func (r *classRequestTutorApplicationRepository) GetByTutorAndRequest(ctx context.Context, tx *gorm.DB, tutorProfileID, requestID string) (entity.ClassRequestTutorApplication, error) {
	if tx == nil {
		tx = r.db
	}

	var application entity.ClassRequestTutorApplication
	if err := tx.WithContext(ctx).Where("tutor_profile_id = ? AND request_id = ?", tutorProfileID, requestID).First(&application).Error; err != nil {
		return entity.ClassRequestTutorApplication{}, err
	}

	return application, nil
}

func (r *classRequestTutorApplicationRepository) GetAllByRequestId(ctx context.Context, tx *gorm.DB, requestID string, preloads ...string) ([]entity.ClassRequestTutorApplication, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var applications []entity.ClassRequestTutorApplication
	if err := tx.WithContext(ctx).Where("request_id = ?", requestID).Find(&applications).Error; err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *classRequestTutorApplicationRepository) GetAllByTutorProfileId(ctx context.Context, tx *gorm.DB, tutorProfileID string, preloads ...string) ([]entity.ClassRequestTutorApplication, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var applications []entity.ClassRequestTutorApplication
	if err := tx.WithContext(ctx).Where("tutor_profile_id = ?", tutorProfileID).Find(&applications).Error; err != nil {
		return nil, err
	}

	return applications, nil
}

func (r *classRequestTutorApplicationRepository) Update(ctx context.Context, tx *gorm.DB, application entity.ClassRequestTutorApplication) (entity.ClassRequestTutorApplication, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&application).Error; err != nil {
		return entity.ClassRequestTutorApplication{}, err
	}

	return application, nil
}
