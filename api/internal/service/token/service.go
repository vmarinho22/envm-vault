package token

import (
	"context"
	"errors"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	tokens *repository.AccessTokenRepository
	users  *repository.UserRepository
}

func New(tokens *repository.AccessTokenRepository, users *repository.UserRepository) *Service {
	return &Service{tokens: tokens, users: users}
}

type MeResponse struct {
	Active    bool       `json:"active"`
	Prefix    string     `json:"prefix,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type CreateResponse struct {
	Token     string    `json:"token"`
	Prefix    string    `json:"prefix"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) GetMine(ctx context.Context, actorID string) (*MeResponse, error) {
	userID, err := parseActorID(actorID)
	if err != nil {
		return nil, err
	}

	token, err := s.tokens.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &MeResponse{Active: false}, nil
		}
		return nil, err
	}

	createdAt := token.CreatedAt
	return &MeResponse{
		Active:    token.RevokedAt == nil,
		Prefix:    token.Prefix,
		CreatedAt: &createdAt,
	}, nil
}

func (s *Service) Create(ctx context.Context, actorID string) (*CreateResponse, error) {
	userID, err := parseActorID(actorID)
	if err != nil {
		return nil, err
	}

	plaintext, err := GeneratePlaintext()
	if err != nil {
		return nil, err
	}

	prefix := DisplayPrefix(plaintext)
	hash := HashToken(plaintext)
	now := time.Now().UTC()

	existing, err := s.tokens.FindByUserID(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existing == nil {
		created, err := s.tokens.Create(ctx, &model.AccessToken{
			UserID:    userID,
			TokenHash: hash,
			Prefix:    prefix,
			CreatedAt: now,
		})
		if err != nil {
			return nil, err
		}
		return &CreateResponse{
			Token:     plaintext,
			Prefix:    created.Prefix,
			CreatedAt: created.CreatedAt,
		}, nil
	}

	if err := s.tokens.Rotate(ctx, existing.ID, hash, prefix, now); err != nil {
		return nil, err
	}

	return &CreateResponse{
		Token:     plaintext,
		Prefix:    prefix,
		CreatedAt: now,
	}, nil
}

func (s *Service) RevokeOwn(ctx context.Context, actorID string) error {
	userID, err := parseActorID(actorID)
	if err != nil {
		return err
	}
	return s.revokeActive(ctx, userID, userID)
}

func (s *Service) RevokeForUser(ctx context.Context, actorID, actorRole, targetUserID string) error {
	if actorRole != service.RoleSuperAdmin && actorRole != service.RoleAdmin {
		return service.ErrForbidden
	}

	actorUUID, err := parseActorID(actorID)
	if err != nil {
		return err
	}

	targetUUID, err := uuid.Parse(targetUserID)
	if err != nil {
		return service.ErrUserNotFound
	}

	target, err := s.users.FindByID(ctx, targetUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrUserNotFound
		}
		return err
	}

	if !CanRevokeTarget(actorRole, target.Role) {
		return service.ErrForbidden
	}

	return s.revokeActive(ctx, targetUUID, actorUUID)
}

func (s *Service) revokeActive(ctx context.Context, userID, revokedBy uuid.UUID) error {
	token, err := s.tokens.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrTokenNotFound
		}
		return err
	}
	if token.RevokedAt != nil {
		return service.ErrTokenNotFound
	}

	return s.tokens.Revoke(ctx, token.ID, revokedBy, time.Now().UTC())
}

func parseActorID(actorID string) (uuid.UUID, error) {
	id, err := uuid.Parse(actorID)
	if err != nil {
		return uuid.Nil, service.ErrUnauthorized
	}
	return id, nil
}

func CanRevokeTarget(actorRole, targetRole string) bool {
	if actorRole == service.RoleSuperAdmin {
		return true
	}
	if actorRole == service.RoleAdmin {
		return targetRole == service.RoleAdmin || targetRole == service.RoleUser
	}
	return false
}
