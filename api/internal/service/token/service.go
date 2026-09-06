package token

import (
	"context"
	"errors"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	tokens *repository.AccessTokenRepository
	users  *repository.UserRepository
	audit  audit.Logger
}

func New(tokens *repository.AccessTokenRepository, users *repository.UserRepository, auditLog audit.Logger) *Service {
	return &Service{tokens: tokens, users: users, audit: audit.OrNop(auditLog)}
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

	actorEmail, actorRole := s.actorIdentity(ctx, userID)

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
		s.audit.Record(ctx, audit.Event{
			Action:      audit.ActionPersonalTokenCreated,
			Summary:     audit.SummaryPersonalTokenCreated,
			ActorUserID: &userID,
			ActorEmail:  actorEmail,
			ActorRole:   actorRole,
		})
		return &CreateResponse{
			Token:     plaintext,
			Prefix:    created.Prefix,
			CreatedAt: created.CreatedAt,
		}, nil
	}

	if err := s.tokens.Rotate(ctx, existing.ID, hash, prefix, now); err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionPersonalTokenRotated,
		Summary:     audit.SummaryPersonalTokenRotated,
		ActorUserID: &userID,
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})

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
	if err := s.revokeActive(ctx, userID, userID); err != nil {
		return err
	}
	actorEmail, actorRole := s.actorIdentity(ctx, userID)
	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionPersonalTokenRevoked,
		Summary:     audit.SummaryPersonalTokenRevoked,
		ActorUserID: &userID,
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})
	return nil
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

	if err := s.revokeActive(ctx, targetUUID, actorUUID); err != nil {
		return err
	}

	actorEmail, _ := s.actorIdentity(ctx, actorUUID)
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionPersonalTokenRevokedAdmin,
		Summary:      audit.SummaryPersonalTokenRevokedAdmin,
		ActorUserID:  &actorUUID,
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		TargetUserID: &targetUUID,
	})
	return nil
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

func (s *Service) actorIdentity(ctx context.Context, userID uuid.UUID) (email, role string) {
	if s.users == nil {
		return "", ""
	}
	user, err := s.users.FindByID(ctx, userID)
	if err != nil || user == nil {
		return "", ""
	}
	return user.Email, user.Role
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
