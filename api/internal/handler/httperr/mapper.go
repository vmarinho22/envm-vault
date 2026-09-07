package httperr

import (
	"errors"
	"log"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func AppHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	var ferr *fiber.Error
	if errors.As(err, &ferr) {
		code = ferr.Code
		if code < fiber.StatusInternalServerError && ferr.Message != "" {
			msg = ferr.Message
		}
	}

	if code >= fiber.StatusInternalServerError {
		log.Printf("request error method=%s path=%s: %v", c.Method(), c.Path(), err)
	}

	return c.Status(code).JSON(fiber.Map{"error": msg})
}

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
	case errors.Is(err, service.ErrUserBlocked):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "user is blocked"})
	case errors.Is(err, service.ErrEmailAlreadyExists):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already registered"})
	case errors.Is(err, service.ErrInvalidRole):
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "invalid role; allowed: admin, user"})
	case errors.Is(err, service.ErrTokenNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "access token not found"})
	case errors.Is(err, service.ErrInvalidTeamName):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid team name"})
	case errors.Is(err, service.ErrTeamNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "team not found"})
	case errors.Is(err, service.ErrTeamSlugTaken):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "team slug already taken"})
	default:
		log.Printf("unhandled error method=%s path=%s: %v", c.Method(), c.Path(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
