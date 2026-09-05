package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	Prefix    string    `gorm:"type:varchar(16);not null"`
	RevokedAt *time.Time
	RevokedBy *uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
