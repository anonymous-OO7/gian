package migrations

import (
	"gian/db"
	"gian/models"
)

func Migrate() error {
	return db.DB.AutoMigrate(
		&models.User{},
		&models.Jobs{},
		&models.Saved{},
		&models.Applications{},
	)
}
