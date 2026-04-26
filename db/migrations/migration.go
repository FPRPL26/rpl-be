package migrations

import (
	"fmt"

	"github.com/FPRPL26/rpl-be/internal/entity"
	mylog "github.com/FPRPL26/rpl-be/internal/pkg/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	fmt.Println(mylog.ColorizeInfo("\n=========== Start Migrate ==========="))
	mylog.Infof("Migrating Tables...")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		return err
	}

	//migrate table
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Skill{},
		&entity.TutorProfile{},
		&entity.Task{},
		&entity.MediaAsset{},
		&entity.RefreshToken{},
		&entity.Class{},
		&entity.Schedule{},
		&entity.Portofolio{},
		&entity.ClassRequest{},
		&entity.BarterSkill{},
		&entity.BarterSkillTransaction{},
		&entity.ClassTransaction{},
		&entity.ClassRequestTransaction{},
		&entity.Review{},
	); err != nil {
		return err
	}

	// Add manual foreign key for TutorProfile where ID is also the UserID
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'fk_tutor_profile_user') THEN
				ALTER TABLE tutor_profile ADD CONSTRAINT fk_tutor_profile_user FOREIGN KEY (id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		return err
	}

	// if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL;`).Error; err != nil {
	// 	return err
	// }

	return nil
}
