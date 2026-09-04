package httperr

import (
	"errors"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func Write(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidEmail):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid email"})
	case errors.Is(err, service.ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	case errors.Is(err, service.ErrInvalidCode):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired verification code"})
	case errors.Is(err, service.ErrTooManyAttempts):
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "too many verification attempts; request a new code"})
	case errors.Is(err, service.ErrUnauthorized):
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	case errors.Is(err, service.ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	case errors.Is(err, service.ErrEmailAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
	case errors.Is(err, service.ErrInvalidRole):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid role; allowed: admin, user"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
