package model

import (
	"time"

	"github.com/google/uuid"
)

// Secret is a named env-var key in the global catalog (values live in SecretValue).
// Projects link to secrets via the project_secrets M2M table.
type Secret struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"type:varchar(128);uniqueIndex;not null"`
	Values    []SecretValue
	Projects  []Project `gorm:"many2many:project_secrets;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
