package seeds

import (
	"encoding/json"
	"os"

	"github.com/FPRPL26/rpl-be/internal/entity"
	mylog "github.com/FPRPL26/rpl-be/internal/pkg/logger"
	"gorm.io/gorm"
)

func SeederTutorProfile(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding tutor profiles...")
	jsonFile, err := os.Open("./db/seeder/data/tutor_profile_data.json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	var listEntity []entity.TutorProfile
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding tutor profiles completed")
	return nil
}
