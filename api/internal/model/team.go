package model

import (
	"time"

	"github.com/google/uuid"
)

type Team struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name          string    `gorm:"type:varchar(80);not null"`
	Slug          string    `gorm:"type:varchar(80);uniqueIndex;not null"`
	AccountableID uuid.UUID `gorm:"type:uuid;not null;index"`
	Accountable   User      `gorm:"foreignKey:AccountableID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
