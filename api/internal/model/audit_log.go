package model

import (
	"time"

	"github.com/google/uuid"
)

// AuditLog is an append-only record of a security-relevant action.
type AuditLog struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OccurredAt   time.Time  `gorm:"not null;index:idx_audit_logs_occurred_at;index:idx_audit_logs_actor_occurred,priority:2,sort:desc"`
	Action       string     `gorm:"type:varchar(64);not null;index:idx_audit_logs_action"`
	ActorUserID  *uuid.UUID `gorm:"type:uuid;index:idx_audit_logs_actor_occurred,priority:1"`
	ActorEmail   string     `gorm:"type:varchar(320);not null;default:''"`
	ActorRole    string     `gorm:"type:varchar(32);not null;default:''"`
	TargetUserID *uuid.UUID `gorm:"type:uuid"`
	ProjectID    *uuid.UUID `gorm:"type:uuid;index:idx_audit_logs_project_id"`
	ProjectTitle string     `gorm:"type:varchar(80);not null;default:''"`
	ProjectSlug  string     `gorm:"type:varchar(80);not null;default:''"`
	AuthChannel  string     `gorm:"type:varchar(32);not null"`
	IP           string     `gorm:"type:varchar(64);not null;default:''"`
	UserAgent    string     `gorm:"type:varchar(512);not null;default:''"`
	Origin       string     `gorm:"type:varchar(512);not null;default:''"`
	CreatedAt    time.Time
}
