package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	BarterSkillRepository interface {
		Create(ctx context.Context, tx *gorm.DB, barter entity.BarterSkill) (entity.BarterSkill, error)
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.BarterSkill, meta.Meta, error)
		GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.BarterSkill, error)
		GetAllByRequestSkill(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, skillID string, preloads ...string) ([]entity.BarterSkill, meta.Meta, error)
		GetAllByOfferedSkill(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, skillID string, preloads ...string) ([]entity.BarterSkill, meta.Meta, error)
		Update(ctx context.Context, tx *gorm.DB, barter entity.BarterSkill) (entity.BarterSkill, error)
		Delete(ctx context.Context, tx *gorm.DB, barter entity.BarterSkill) error
	}

	barterSkillRepository struct {
		db *gorm.DB
	}
)

func NewBarterSkillRepository(db *gorm.DB) BarterSkillRepository {
	return &barterSkillRepository{db}
}

func (r *barterSkillRepository) Create(ctx context.Context, tx *gorm.DB, barter entity.BarterSkill) (entity.BarterSkill, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&barter).Error; err != nil {
		return barter, err
	}

	return barter, nil
}

func (r *barterSkillRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, preloads ...string) ([]entity.BarterSkill, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var barters []entity.BarterSkill

	if metaReq.SortBy == "" {
		metaReq.SortBy = "created_at"
	}

	if metaReq.Sort == "" {
		metaReq.Sort = "desc"
	}

	tx = tx.WithContext(ctx).Model(entity.BarterSkill{}).
		Joins("LEFT JOIN skills AS rs ON rs.id = barter_skills.request_skills").
		Joins("LEFT JOIN skills AS os ON os.id = barter_skills.offered_skills")

	if err := WithFilters(tx, &metaReq,
		AddModels(entity.BarterSkill{}),
		AddCustomField("search", "barter_skills.name ILIKE ?"),
		AddCustomField("tutor_profile_id", "barter_skills.tutor_profile_id = ?"),
		AddCustomField("request_skill_id", "barter_skills.request_skills = ?"),
		AddCustomField("offered_skill_id", "barter_skills.offered_skills = ?"),
		AddCustomField("request_skill_name", "rs.name ILIKE ?"),
		AddCustomField("offered_skill_name", "os.name ILIKE ?"),
	).Find(&barters).Error; err != nil {
		return nil, metaReq, err
	}

	return barters, metaReq, nil
}

func (r *barterSkillRepository) GetById(ctx context.Context, tx *gorm.DB, id string, preloads ...string) (entity.BarterSkill, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var barter entity.BarterSkill
	if err := tx.WithContext(ctx).Take(&barter, "id = ?", id).Error; err != nil {
		return entity.BarterSkill{}, err
	}

	return barter, nil
}

func (r *barterSkillRepository) Update(ctx context.Context, tx *gorm.DB, barter entity.BarterSkill) (entity.BarterSkill, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Save(&barter).Error; err != nil {
		return barter, err
	}

	return barter, nil
}

func (r *barterSkillRepository) Delete(ctx context.Context, tx *gorm.DB, barter entity.BarterSkill) error {
	if tx == nil {
		tx = r.db
	}

	return tx.WithContext(ctx).Delete(&barter).Error
}

func (r *barterSkillRepository) GetAllByRequestSkill(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, skillID string, preloads ...string) ([]entity.BarterSkill, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var barters []entity.BarterSkill
	if err := WithFilters(tx, &metaReq,
		AddModels(entity.BarterSkill{}),
		AddCustomField("request_skill_id", "barter_skills.request_skills = ?"),
	).Find(&barters).Error; err != nil {
		return nil, metaReq, err
	}

	return barters, metaReq, nil
}

func (r *barterSkillRepository) GetAllByOfferedSkill(ctx context.Context, tx *gorm.DB, metaReq meta.Meta, skillID string, preloads ...string) ([]entity.BarterSkill, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	for _, preload := range preloads {
		tx = tx.Preload(preload)
	}

	var barters []entity.BarterSkill
	if err := WithFilters(tx, &metaReq,
		AddModels(entity.BarterSkill{}),
		AddCustomField("offered_skill_id", "barter_skills.offered_skills = ?"),
	).Find(&barters).Error; err != nil {
		return nil, metaReq, err
	}

	return barters, metaReq, nil
}
