package model

import (
	"time"

	"github.com/google/uuid"
)

type Environment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title     string    `gorm:"type:varchar(80);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
