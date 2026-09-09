package maps

import (
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapAuth(err error) *Mapped {
	return firstMatch(
		match(err, service.ErrInvalidEmail, fiber.StatusBadRequest, "invalid email"),
		match(err, service.ErrInvalidCode, fiber.StatusUnauthorized, "invalid or expired verification code"),
		match(err, service.ErrTooManyAttempts, fiber.StatusTooManyRequests, "too many verification attempts; request a new code"),
		match(err, service.ErrUnauthorized, fiber.StatusUnauthorized, "unauthorized"),
		match(err, service.ErrForbidden, fiber.StatusForbidden, "forbidden"),
		match(err, service.ErrEmailAlreadyExists, fiber.StatusConflict, "email already registered"),
		match(err, service.ErrInvalidRole, fiber.StatusUnprocessableEntity, "invalid role; allowed: admin, user"),
	)
}
