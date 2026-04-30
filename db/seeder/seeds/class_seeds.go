package seeds

import (
	"encoding/json"
	"os"

	"github.com/FPRPL26/rpl-be/internal/entity"
	mylog "github.com/FPRPL26/rpl-be/internal/pkg/logger"
	"gorm.io/gorm"
)

func SeederClass(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding classes...")
	jsonFile, err := os.Open("./db/seeder/data/class_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var listEntity []entity.Class
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding classes completed")
	return nil
}

func SeederSchedule(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding schedules...")
	jsonFile, err := os.Open("./db/seeder/data/class_schedule_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var listEntity []entity.Schedule
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding schedules completed")
	return nil
}

func SeederClassTransaction(db *gorm.DB) error {
	mylog.Infof("[PROCESS] Seeding class transactions...")
	jsonFile, err := os.Open("./db/seeder/data/class_transaction_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var listEntity []entity.ClassTransaction
	if err := json.NewDecoder(jsonFile).Decode(&listEntity); err != nil {
		return err
	}

	for _, entity := range listEntity {
		if err := db.Save(&entity).Error; err != nil {
			return err
		}
	}

	mylog.Infof("[COMPLETE] Seeding class transactions completed")
	return nil
}
