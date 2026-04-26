package cron

import (
	"strings"
	"time"

	"github.com/FPRPL26/rpl-be/internal/entity"
	mylog "github.com/FPRPL26/rpl-be/internal/pkg/logger"
	"github.com/FPRPL26/rpl-be/internal/utils"
	"gorm.io/gorm"
)

func StartMediaCron(db *gorm.DB) {
	ticker := time.NewTicker(24 * time.Hour) // Run every 24 hours
	go func() {
		for {
			select {
			case <-ticker.C:
				CleanUnusedMedia(db)
			}
		}
	}()

	// Also run once at startup
	go CleanUnusedMedia(db)
}

func CleanUnusedMedia(db *gorm.DB) {
	mylog.Infoln("Starting cleanup of unused media assets...")

	var unusedMedia []entity.MediaAsset
	// Find media assets that are not used and were created more than 1 hour ago (to avoid deleting files currently being uploaded/processed)
	if err := db.Where("is_used = ? AND created_at < ?", false, time.Now().Add(-1*time.Hour)).Find(&unusedMedia).Error; err != nil {
		mylog.Errorf("Failed to fetch unused media: %v", err)
		return
	}

	if len(unusedMedia) == 0 {
		mylog.Infoln("No unused media assets found for cleanup.")
		return
	}

	for _, media := range unusedMedia {
		// Extract filename from URL (assuming format: http://host/api/static/filename)
		parts := strings.Split(media.URL, "/")
		if len(parts) > 0 {
			filename := parts[len(parts)-1]

			// Delete file from storage
			if err := utils.DeleteFile(filename); err != nil {
				mylog.Errorf("Failed to delete file %s: %v", filename, err)
				// Even if file deletion fails (e.g. file already gone), we might want to continue or remove from DB
			} else {
				mylog.Infof("Deleted unused file: %s", filename)
			}
		}

		// Delete record from database
		if err := db.Unscoped().Delete(&media).Error; err != nil {
			mylog.Errorf("Failed to delete media record %s: %v", media.ID, err)
		}
	}

	mylog.Infof("Cleanup completed. Processed %d assets.", len(unusedMedia))
}
