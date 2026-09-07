package model

import (
	"time"

	"github.com/google/uuid"
)

// SecretValue holds the encrypted value of a secret for one environment.
type SecretValue struct {
	ID            uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SecretID      uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_secret_values_secret_env,priority:1;index"`
	Secret        Secret      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	EnvironmentID uuid.UUID   `gorm:"type:uuid;not null;uniqueIndex:idx_secret_values_secret_env,priority:2;index"`
	Environment   Environment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Ciphertext    []byte      `gorm:"type:bytea;not null"`
	KeyVersion    int         `gorm:"not null;default:1"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
