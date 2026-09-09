package model

import (
	"github.com/google/uuid"
)

// ProjectSecret is the M2M junction linking a project to a catalog secret.
type ProjectSecret struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey;index"`
	SecretID  uuid.UUID `gorm:"type:uuid;primaryKey;index"`
}
