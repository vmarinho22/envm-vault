package database

import (
	"fmt"

	"github.com/vmarinho22/envm-vault/api/internal/model"
	"gorm.io/gorm"
)

// DefaultTags are seeded on migrate/startup when missing (case-insensitive title).
var DefaultTags = []model.Tag{
	{Title: "backend", Color: "#000000"},
	{Title: "frontend", Color: "#FFFFFF"},
	{Title: "infra", Color: "#F97316"},
}

// EnsureDefaultTags inserts the catalog defaults if a tag with the same title
// (case-insensitive) does not already exist. Existing rows are left unchanged.
func EnsureDefaultTags(db *gorm.DB) error {
	for _, def := range DefaultTags {
		var count int64
		if err := db.Model(&model.Tag{}).Where("LOWER(title) = LOWER(?)", def.Title).Count(&count).Error; err != nil {
			return fmt.Errorf("check default tag %q: %w", def.Title, err)
		}
		if count > 0 {
			continue
		}
		tag := model.Tag{Title: def.Title, Color: def.Color}
		if err := db.Create(&tag).Error; err != nil {
			return fmt.Errorf("seed default tag %q: %w", def.Title, err)
		}
	}
	return nil
}
