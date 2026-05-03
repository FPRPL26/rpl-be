package repository

import (
	"github.com/FPRPL26/rpl-be/internal/entity"
	"gorm.io/gorm"
)

type MediaAssetRepository interface {
	MarkAsUsed(url string) error
	MarkAsUnused(url string) error
}

type mediaAssetRepository struct {
	db *gorm.DB
}

func NewMediaAsset(db *gorm.DB) MediaAssetRepository {
	return &mediaAssetRepository{db}
}

func (r *mediaAssetRepository) MarkAsUsed(url string) error {
	if url == "" {
		return nil
	}
	return r.db.Model(&entity.MediaAsset{}).Where("url = ?", url).Update("is_used", true).Error
}

func (r *mediaAssetRepository) MarkAsUnused(url string) error {
	if url == "" {
		return nil
	}
	return r.db.Model(&entity.MediaAsset{}).Where("url = ?", url).Update("is_used", false).Error
}
