package user

import (
	"context"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/google/uuid"
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
		Summary:      audit.SummaryUserRegistered,
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
			ID:    user.ID.String(),
			Email: user.Email,
			Role:  user.Role,
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
