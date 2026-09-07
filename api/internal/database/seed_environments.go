package database

import (
	"fmt"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"gorm.io/gorm"
)

// DefaultEnvironments are seeded on migrate/startup when missing (case-insensitive title).
var DefaultEnvironments = []model.Environment{
	{Title: "dev"},
	{Title: "test"},
	{Title: "staging"},
	{Title: "production"},
}

// EnsureDefaultEnvironments inserts the catalog defaults if an environment with the
// same title (case-insensitive) does not already exist. Existing rows are left unchanged.
func EnsureDefaultEnvironments(db *gorm.DB) error {
	for _, def := range DefaultEnvironments {
		var count int64
		if err := db.Model(&model.Environment{}).Where("LOWER(title) = LOWER(?)", def.Title).Count(&count).Error; err != nil {
			return fmt.Errorf("check default environment %q: %w", def.Title, err)
		}
		if count > 0 {
			continue
		}
		env := model.Environment{Title: def.Title}
		if err := db.Create(&env).Error; err != nil {
			return fmt.Errorf("seed default environment %q: %w", def.Title, err)
		}
	}
	return nil
}
