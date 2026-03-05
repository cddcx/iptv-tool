package model

import (
	"log/slog"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto-migrate the schema
	err = DB.AutoMigrate(
		&User{},
		&LiveSource{},
		&EPGSource{},
		&ChannelLogo{},
		&PublishInterface{},
		&AggregationRule{},
		&ParsedChannel{},
		&ParsedEPG{},
	)
	if err != nil {
		return err
	}

	// Defensive reset: set is_syncing to false for all sources on startup
	DB.Model(&LiveSource{}).Where("is_syncing = ?", true).Update("is_syncing", false)
	DB.Model(&EPGSource{}).Where("is_syncing = ?", true).Update("is_syncing", false)

	slog.Info("Database initialized and migrated successfully.")
	return nil
}
