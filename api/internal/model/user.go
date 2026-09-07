package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Role      string    `gorm:"type:varchar(20);not null"`
	BlockedAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
