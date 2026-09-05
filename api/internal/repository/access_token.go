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

type AccessTokenRepository struct {
	db *gorm.DB
}

func NewAccessTokenRepository(db *gorm.DB) *AccessTokenRepository {
	return &AccessTokenRepository{db: db}
}

func (r *AccessTokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*model.AccessToken, error) {
	var token model.AccessToken
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find access token by user id: %w", err)
	}
	return &token, nil
}

func (r *AccessTokenRepository) FindByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]model.AccessToken, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	var tokens []model.AccessToken
	if err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("find access tokens by user ids: %w", err)
	}
	return tokens, nil
}

func (r *AccessTokenRepository) Create(ctx context.Context, token *model.AccessToken) (*model.AccessToken, error) {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}
	return token, nil
}

func (r *AccessTokenRepository) Rotate(ctx context.Context, id uuid.UUID, tokenHash, prefix string, createdAt time.Time) error {
	err := r.db.WithContext(ctx).
		Model(&model.AccessToken{}).
		Where("id = ?", id).
		Select("token_hash", "prefix", "revoked_at", "revoked_by", "created_at").
		Updates(map[string]any{
			"token_hash": tokenHash,
			"prefix":     prefix,
			"revoked_at": nil,
			"revoked_by": nil,
			"created_at": createdAt,
		}).Error
	if err != nil {
		return fmt.Errorf("rotate access token: %w", err)
	}
	return nil
}

func (r *AccessTokenRepository) Revoke(ctx context.Context, id, revokedBy uuid.UUID, revokedAt time.Time) error {
	err := r.db.WithContext(ctx).
		Model(&model.AccessToken{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"revoked_at": revokedAt,
			"revoked_by": revokedBy,
		}).Error
	if err != nil {
		return fmt.Errorf("revoke access token: %w", err)
	}
	return nil
}
