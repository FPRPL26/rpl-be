package repository

import (
	"context"

	"github.com/FPRPL26/rpl-be/internal/entity"
	"github.com/FPRPL26/rpl-be/internal/pkg/meta"
	"gorm.io/gorm"
)

type (
	SkillRepository interface {
		GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta) ([]entity.Skill, meta.Meta, error)
		GetById(ctx context.Context, tx *gorm.DB, id int64) (entity.Skill, error)
	}

	skillRepository struct {
		db *gorm.DB
	}
)

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepository{db}
}

func (r *skillRepository) GetAll(ctx context.Context, tx *gorm.DB, metaReq meta.Meta) ([]entity.Skill, meta.Meta, error) {
	if tx == nil {
		tx = r.db
	}

	var skills []entity.Skill
	tx = tx.WithContext(ctx).Model(entity.Skill{})

	if err := WithFilters(tx, &metaReq, AddModels(entity.Skill{})).Find(&skills).Error; err != nil {
		return nil, metaReq, err
	}

	return skills, metaReq, nil
}

func (r *skillRepository) GetById(ctx context.Context, tx *gorm.DB, id int64) (entity.Skill, error) {
	if tx == nil {
		tx = r.db
	}

	var skill entity.Skill
	if err := tx.WithContext(ctx).First(&skill, "id = ?", id).Error; err != nil {
		return entity.Skill{}, err
	}

	return skill, nil
}
