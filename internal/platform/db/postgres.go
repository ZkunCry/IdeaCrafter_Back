package db

import (
	"context"
	"fmt"
	"time"

	"startup_back/internal/platform/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(cfg *config.AppConfig) (*gorm.DB, error) {
	gormLogLevel := logger.Warn
	if !cfg.IsProduction() {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DBConnectionString()), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to access sql.DB: %w", err)
	}

	poolSize := cfg.Database.PoolSize
	if poolSize <= 0 {
		poolSize = 10
	}
	idleConn := cfg.Database.MaxIdleConn
	if idleConn <= 0 || idleConn > poolSize {
		idleConn = poolSize / 2
	}

	sqlDB.SetMaxOpenConns(poolSize)
	sqlDB.SetMaxIdleConns(idleConn)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
