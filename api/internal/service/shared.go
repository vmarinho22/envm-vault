package service

import (
	"net/mail"
	"strings"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
)

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleUser       = "user"

	// ProductionEnvironmentTitle is the protected catalog environment title.
	ProductionEnvironmentTitle = "production"
)

type UserResponse struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	BlockedAt *time.Time `json:"blocked_at,omitempty"`
}

type TokenStatusResponse struct {
	Prefix    string     `json:"prefix"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
}

type UserWithTokenResponse struct {
	ID        string               `json:"id"`
	Email     string               `json:"email"`
	Role      string               `json:"role"`
	BlockedAt *time.Time           `json:"blocked_at,omitempty"`
	Token     *TokenStatusResponse `json:"token"`
}

type VerifyResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func NormalizeEmail(raw string) (string, error) {
	addr, err := mail.ParseAddress(strings.TrimSpace(raw))
	if err != nil {
		return "", ErrInvalidEmail
	}
	return strings.ToLower(addr.Address), nil
}

func ToUserResponse(user model.User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Role:      user.Role,
		BlockedAt: user.BlockedAt,
	}
}

// CanManageTarget reports whether actorRole may manage (revoke token / block) a user with targetRole.
func CanManageTarget(actorRole, targetRole string) bool {
	if actorRole == RoleSuperAdmin {
		return true
	}
	if actorRole == RoleAdmin {
		return targetRole == RoleAdmin || targetRole == RoleUser
	}
	return false
}

func IsBlocked(user *model.User) bool {
	return user != nil && user.BlockedAt != nil
}
