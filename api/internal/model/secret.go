package model

import (
	"time"

	"github.com/google/uuid"
)

// Secret is a named env-var key scoped to a project (values live in SecretValue).
type Secret struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_secrets_project_name,priority:1;index"`
	Project   Project   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Name      string    `gorm:"type:varchar(128);not null;uniqueIndex:idx_secrets_project_name,priority:2"`
	Values    []SecretValue
	CreatedAt time.Time
	UpdatedAt time.Time
}
