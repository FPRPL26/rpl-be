package seeds

import (
	"encoding/json"
	"os"

	"github.com/FPRPL26/rpl-be/internal/entity"
	mylog "github.com/FPRPL26/rpl-be/internal/pkg/logger"
	"gorm.io/gorm"
)

func SeederSkills(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding skills...")
	jsonFile, err := os.Open("./db/seeder/data/skills_data.json")
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	var listEntity []entity.Skill
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding skills completed")
	return nil
}
