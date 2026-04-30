package seeders

import (
	"fmt"

	"github.com/FPRPL26/rpl-be/db/seeder/seeds"
	mylog "github.com/FPRPL26/rpl-be/internal/pkg/logger"
	"gorm.io/gorm"
)

func Seeding(db *gorm.DB) error {
	seeders := []func(*gorm.DB) error{
		seeds.SeederUser,
		seeds.SeederTutorProfile,
		seeds.SeederClass,
		seeds.SeederSchedule,
		seeds.SeederClassTransaction,
		seeds.SeederSkills,
	}

	fmt.Println(mylog.ColorizeInfo("\n=========== Start Seeding ==========="))
	for _, seeder := range seeders {
		if err := seeder(db); err != nil {
			return err
		}
	}

	return nil
}
