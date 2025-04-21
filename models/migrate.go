// models/migrate.go
package models

import (
	"gorm.io/gorm"
)

// MigrateDB runs auto-migration for all models
func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		// Add other models here
	)
}
