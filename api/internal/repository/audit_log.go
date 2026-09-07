package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

// AuditLogListFilter controls listing of audit_logs.
type AuditLogListFilter struct {
	Page          int
	PageSize      int
	SortAsc       bool
	ActorUserIDs  []uuid.UUID
	TargetUserIDs []uuid.UUID
	ProjectIDs    []uuid.UUID
	OccurredFrom  *time.Time
	OccurredTo    *time.Time
}

func (r *AuditLogRepository) List(ctx context.Context, filter AuditLogListFilter) ([]model.AuditLog, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.AuditLog{})

	if len(filter.ActorUserIDs) > 0 {
		q = q.Where("actor_user_id IN ?", filter.ActorUserIDs)
	}
	if len(filter.TargetUserIDs) > 0 {
		q = q.Where("target_user_id IN ?", filter.TargetUserIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		q = q.Where("project_id IN ?", filter.ProjectIDs)
	}
	if filter.OccurredFrom != nil {
		q = q.Where("occurred_at >= ?", *filter.OccurredFrom)
	}
	if filter.OccurredTo != nil {
		q = q.Where("occurred_at <= ?", *filter.OccurredTo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count audit logs: %w", err)
	}

	order := "occurred_at DESC"
	if filter.SortAsc {
		order = "occurred_at ASC"
	}

	offset := (filter.Page - 1) * filter.PageSize
	var rows []model.AuditLog
	if err := q.Order(order).Offset(offset).Limit(filter.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list audit logs: %w", err)
	}
	return rows, total, nil
}
