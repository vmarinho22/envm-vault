package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title        string        `gorm:"type:varchar(80);not null;index"`
	Slug         string        `gorm:"type:varchar(80);uniqueIndex;not null"`
	Teams        []Team        `gorm:"many2many:project_teams;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Environments []Environment `gorm:"many2many:project_environments;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Tags         []Tag         `gorm:"many2many:project_tags;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
