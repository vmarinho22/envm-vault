package service

import (
	"net/mail"
	"strings"

	"github.com/envm-vault/envm-vault/api/internal/model"
)

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
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
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}
}
