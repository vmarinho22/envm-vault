package user

import (
	"context"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
)

type Service struct {
	users *repository.UserRepository
}

func New(users *repository.UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) Register(ctx context.Context, actorRole, rawEmail, role string) (*service.UserResponse, error) {
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

	resp := service.ToUserResponse(*user)
	return &resp, nil
}
