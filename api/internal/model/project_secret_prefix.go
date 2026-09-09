package model

import (
	"time"

	"github.com/google/uuid"
)

// ProjectSecretPrefix is a download-time name prefix for a secret in one
// project environment (e.g. VITE_). Catalog secret names stay unchanged.
type ProjectSecretPrefix struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID     uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_psp_project_secret_env,priority:1;index"`
	Project       Project     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SecretID      uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_psp_project_secret_env,priority:2;index"`
	Secret        Secret      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EnvironmentID uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_psp_project_secret_env,priority:3;index"`
	Environment   Environment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Prefix        string      `gorm:"type:varchar(64);not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
