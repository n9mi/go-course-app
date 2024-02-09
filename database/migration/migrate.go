package migration

import (
	"github.com/n9mi/go-course-app/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Role{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Category{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Course{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.PaymentMethod{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Purchase{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.CourseMember{}); err != nil {
		return err
	}

	return nil
}

func Drop(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&entity.CourseMember{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Purchase{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.PaymentMethod{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Course{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Category{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable("user_roles"); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Role{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.User{}); err != nil {
		return err
	}

	return nil
}
