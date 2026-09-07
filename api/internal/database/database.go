package database

import (
	"context"
	"fmt"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.LoginCode{},
		&model.AccessToken{},
		&model.AuditLog{},
		&model.Team{},
		&model.Tag{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	// AutoMigrate does not drop removed columns; drop audit_logs.summary if present.
	if db.Migrator().HasColumn(&model.AuditLog{}, "summary") {
		if err := db.Migrator().DropColumn(&model.AuditLog{}, "summary"); err != nil {
			return nil, fmt.Errorf("drop audit_logs.summary: %w", err)
		}
	}

	if err := EnsureDefaultTags(db); err != nil {
		return nil, fmt.Errorf("seed default tags: %w", err)
	}

	return db, nil
}

func WithTx(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(fn)
}

func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql db: %w", err)
	}
	return sqlDB.Ping()
}
