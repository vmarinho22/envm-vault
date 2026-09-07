package service

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCode        = errors.New("invalid or expired verification code")
	ErrTooManyAttempts    = errors.New("too many verification attempts")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidRole        = errors.New("invalid role")
	ErrTokenNotFound      = errors.New("access token not found")
	ErrInvalidTeamName    = errors.New("invalid team name")
	ErrTeamNotFound       = errors.New("team not found")
	ErrTeamSlugTaken      = errors.New("team slug already taken")
)
