package db

import (
	"startup_back/internal/platform/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(cfg *config.AppConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DBConnectionString()), &gorm.Config{})
}
