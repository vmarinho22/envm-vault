package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	users  *repository.UserRepository
	tokens *repository.AccessTokenRepository
	audit  audit.Logger
}

func New(users *repository.UserRepository, tokens *repository.AccessTokenRepository, auditLog audit.Logger) *Service {
	return &Service{users: users, tokens: tokens, audit: audit.OrNop(auditLog)}
}

func (s *Service) Register(ctx context.Context, actorID, actorEmail, actorRole, rawEmail, role string) (*service.UserResponse, error) {
	if actorRole != service.RoleSuperAdmin && actorRole != service.RoleAdmin {
		return nil, service.ErrForbidden
	}

	emailAddr, err := service.NormalizeEmail(rawEmail)
	if err != nil {
		return nil, service.ErrInvalidEmail
	}

	if role != service.RoleAdmin && role != service.RoleUser {
		return nil, service.ErrInvalidRole
	}

	exists, err := s.users.ExistsByEmail(ctx, emailAddr)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, service.ErrEmailAlreadyExists
	}

	user, err := s.users.Create(ctx, &model.User{
		Email: emailAddr,
		Role:  role,
	})
	if err != nil {
		return nil, err
	}

	var actorPtr *uuid.UUID
	if actorUUID, parseErr := uuid.Parse(actorID); parseErr == nil {
		actorPtr = &actorUUID
	}

	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionUserRegistered,
		ActorUserID:  actorPtr,
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		TargetUserID: &user.ID,
	})

	resp := service.ToUserResponse(*user)
	return &resp, nil
}

func (s *Service) List(ctx context.Context, actorRole string) ([]service.UserWithTokenResponse, error) {
	if actorRole != service.RoleSuperAdmin && actorRole != service.RoleAdmin {
		return nil, service.ErrForbidden
	}

	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		ids = append(ids, user.ID)
	}

	tokens, err := s.tokens.FindByUserIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	byUser := make(map[uuid.UUID]*model.AccessToken, len(tokens))
	for i := range tokens {
		token := &tokens[i]
		byUser[token.UserID] = token
	}

	out := make([]service.UserWithTokenResponse, 0, len(users))
	for _, user := range users {
		item := service.UserWithTokenResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			Role:      user.Role,
			BlockedAt: user.BlockedAt,
		}
		if token := byUser[user.ID]; token != nil {
			item.Token = &service.TokenStatusResponse{
				Prefix:    token.Prefix,
				CreatedAt: token.CreatedAt,
				RevokedAt: token.RevokedAt,
			}
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *Service) Block(ctx context.Context, actorID, actorEmail, actorRole, targetUserID string) (*service.UserResponse, error) {
	return s.setBlocked(ctx, actorID, actorEmail, actorRole, targetUserID, true)
}

func (s *Service) Unblock(ctx context.Context, actorID, actorEmail, actorRole, targetUserID string) (*service.UserResponse, error) {
	return s.setBlocked(ctx, actorID, actorEmail, actorRole, targetUserID, false)
}

func (s *Service) setBlocked(ctx context.Context, actorID, actorEmail, actorRole, targetUserID string, block bool) (*service.UserResponse, error) {
	if actorRole != service.RoleSuperAdmin && actorRole != service.RoleAdmin {
		return nil, service.ErrForbidden
	}

	actorUUID, err := uuid.Parse(actorID)
	if err != nil {
		return nil, service.ErrUnauthorized
	}

	targetUUID, err := uuid.Parse(strings.TrimSpace(targetUserID))
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	if actorUUID == targetUUID {
		return nil, service.ErrForbidden
	}

	target, err := s.users.FindByID(ctx, targetUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}

	if !service.CanManageTarget(actorRole, target.Role) {
		return nil, service.ErrForbidden
	}

	alreadyBlocked := service.IsBlocked(target)
	if block == alreadyBlocked {
		resp := service.ToUserResponse(*target)
		return &resp, nil
	}

	now := time.Now().UTC()
	if block {
		target.BlockedAt = &now
	} else {
		target.BlockedAt = nil
	}

	if err := s.users.Update(ctx, target); err != nil {
		return nil, err
	}

	if block {
		_ = s.revokeActiveToken(ctx, targetUUID, actorUUID)
	}

	action := audit.ActionUserUnblocked
	if block {
		action = audit.ActionUserBlocked
	}
	s.audit.Record(ctx, audit.Event{
		Action:       action,
		ActorUserID:  &actorUUID,
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		TargetUserID: &targetUUID,
	})

	resp := service.ToUserResponse(*target)
	return &resp, nil
}

func (s *Service) revokeActiveToken(ctx context.Context, userID, revokedBy uuid.UUID) error {
	if s.tokens == nil {
		return nil
	}
	token, err := s.tokens.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if token.RevokedAt != nil {
		return nil
	}
	return s.tokens.Revoke(ctx, token.ID, revokedBy, time.Now().UTC())
}
