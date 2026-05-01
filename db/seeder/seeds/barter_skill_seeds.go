package seeds

import (
	"encoding/json"
	"io"
	"os"

	"github.com/FPRPL26/rpl-be/internal/entity"
	mylog "github.com/FPRPL26/rpl-be/internal/pkg/logger"
	"gorm.io/gorm"
)

func BarterSkillSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("db/seeder/data/barter_skill_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var barters []entity.BarterSkill
	json.Unmarshal(byteValue, &barters)

	for _, barter := range barters {
		if err := db.FirstOrCreate(&barter, entity.BarterSkill{ID: barter.ID}).Error; err != nil {
			mylog.Errorln(err)
		}
	}
	mylog.Infof("[COMPLETE] Seeding barter skill transactions completed")
	return nil
}
