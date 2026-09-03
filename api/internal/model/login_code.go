package model

import (
	"time"

	"github.com/google/uuid"
)

type LoginCode struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string     `gorm:"index;not null"`
	CodeHash  string     `gorm:"not null"`
	ExpiresAt time.Time  `gorm:"index;not null"`
	UsedAt    *time.Time `gorm:"index"`
	CreatedAt time.Time
}
