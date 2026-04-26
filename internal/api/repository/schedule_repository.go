package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	ScheduleRepository interface {
		Create(ctx context.Context, tx *gorm.DB, schedule entity.Schedule) (entity.Schedule, error)
		GetAllByClassId(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, classId string, preloads ...string) ([]entity.Schedule, meta.Meta, error)
		GetById(ctx context.Context, tx *gorm.DB, scheduleId string, preloads ...string) (entity.Schedule, error)
		Update(ctx context.Context, tx *gorm.DB, schedule entity.Schedule) (entity.Schedule, error)
		Delete(ctx context.Context, tx *gorm.DB, schedule entity.Schedule) error
		AddSchedules(ctx context.Context, tx *gorm.DB, schedules []entity.Schedule) error
	}

	scheduleRepository struct {
		db *gorm.DB
	}
)

func NewScheduleRepository(db *gorm.DB) ScheduleRepository {
	return &scheduleRepository{db}
}

func (r *scheduleRepository) Create(ctx context.Context, tx *gorm.DB, schedule entity.Schedule) (entity.Schedule, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&schedule).Error; err != nil {
		return schedule, err
	}

	return schedule, nil
}

func (r *scheduleRepository) GetAllByClassId(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, classId string, preloads ...string) ([]entity.Schedule, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var schedules []entity.Schedule

	tx = tx.WithContext(ctx).Model(entity.Schedule{}).Where("class_id = ?", classId)

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.Schedule{})).Find(&schedules).Error; err != nil {
		return nil, metaReq, err
	}

	return schedules, metaReq, nil
}

func (r *scheduleRepository) GetById(ctx context.Context, tx *gorm.DB, scheduleId string, preloads ...string) (entity.Schedule, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var schedule entity.Schedule
	if err := tx.WithContext(ctx).Take(&schedule, "id = ?", scheduleId).Error; err != nil {
		return entity.Schedule{}, err
	}

	return schedule, nil
}

func (r *scheduleRepository) Update(ctx context.Context, tx *gorm.DB, schedule entity.Schedule) (entity.Schedule, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&schedule).Error; err != nil {
		return entity.Schedule{}, err
	}

	return schedule, nil
}

func (r *scheduleRepository) Delete(ctx context.Context, tx *gorm.DB, schedule entity.Schedule) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&schedule).Error; err != nil {
		return err
	}

	return nil
}

func (r *scheduleRepository) AddSchedules(ctx context.Context, tx *gorm.DB, schedules []entity.Schedule) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&schedules).Error; err != nil {
		return err
	}

	return nil
}
