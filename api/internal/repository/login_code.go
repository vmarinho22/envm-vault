package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoginCodeRepository struct {
	db *gorm.DB
}

func NewLoginCodeRepository(db *gorm.DB) *LoginCodeRepository {
	return &LoginCodeRepository{db: db}
}

func (r *LoginCodeRepository) WithDB(db *gorm.DB) *LoginCodeRepository {
	return &LoginCodeRepository{db: db}
}

func (r *LoginCodeRepository) InvalidatePrevious(ctx context.Context, email string) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&model.LoginCode{}).
		Where("email = ? AND used_at IS NULL", email).
		Update("used_at", now).Error
	if err != nil {
		return fmt.Errorf("invalidate previous codes: %w", err)
	}
	return nil
}

func (r *LoginCodeRepository) Create(ctx context.Context, code *model.LoginCode) (*model.LoginCode, error) {
	if err := r.db.WithContext(ctx).Create(code).Error; err != nil {
		return nil, fmt.Errorf("create login code: %w", err)
	}
	return code, nil
}

func (r *LoginCodeRepository) FindLatestValid(ctx context.Context, email string) (*model.LoginCode, error) {
	var code model.LoginCode
	err := r.db.WithContext(ctx).
		Where("email = ? AND used_at IS NULL AND expires_at > ?", email, time.Now()).
		Order("created_at DESC").
		First(&code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find latest valid login code: %w", err)
	}
	return &code, nil
}

func (r *LoginCodeRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&model.LoginCode{}).
		Where("id = ?", id).
		Update("used_at", now).Error
	if err != nil {
		return fmt.Errorf("mark login code used: %w", err)
	}
	return nil
}

func (r *LoginCodeRepository) IncrementAttempts(ctx context.Context, id uuid.UUID) (int, error) {
	var attemptCount int
	err := r.db.WithContext(ctx).Raw(
		`UPDATE login_codes SET attempt_count = attempt_count + 1 WHERE id = ? RETURNING attempt_count`,
		id,
	).Scan(&attemptCount).Error
	if err != nil {
		return 0, fmt.Errorf("increment login code attempts: %w", err)
	}
	return attemptCount, nil
}
