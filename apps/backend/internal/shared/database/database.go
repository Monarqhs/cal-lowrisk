// Package database opens the GORM connection used for runtime queries only.
// Schema (DDL) is owned by golang-migrate, NOT GORM AutoMigrate
// (see docs/02-system/architecture.md §5.1).
package database

import (
	"fmt"

	"github.com/Monarqhs/cal-lowrisk/apps/backend/internal/shared/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New opens a GORM DB using the pooled connection URL and applies the small
// connection-pool limits required by Neon's free tier (deployment.md §3).
func New(cfg *config.Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
		// Schema is managed by migrations; disable implicit behaviours.
		SkipDefaultTransaction: true,
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("database: open: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database: ping: %w", err)
	}
	return db, nil
}
