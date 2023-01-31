package database

import (
	"github.com/MicBun/go-100-coverage-docker-crud/core"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(core.User{}); err != nil {
		return err
	}
	return nil
}
