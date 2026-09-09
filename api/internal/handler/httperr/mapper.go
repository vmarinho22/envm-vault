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
	case errors.Is(err, service.ErrInvalidTagTitle):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tag title"})
	case errors.Is(err, service.ErrInvalidTagColor):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tag color"})
	case errors.Is(err, service.ErrTagNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tag not found"})
	case errors.Is(err, service.ErrTagTitleTaken):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "tag title already taken"})
	case errors.Is(err, service.ErrInvalidEnvironmentTitle):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid environment title"})
	case errors.Is(err, service.ErrEnvironmentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "environment not found"})
	case errors.Is(err, service.ErrEnvironmentTitleTaken):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "environment title already taken"})
	case errors.Is(err, service.ErrInvalidProjectTitle):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid project title"})
	case errors.Is(err, service.ErrProjectNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "project not found"})
	case errors.Is(err, service.ErrProjectSlugTaken):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "project slug already taken"})
	case errors.Is(err, service.ErrInvalidProjectRelations):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "project requires at least one team and one environment"})
	case errors.Is(err, service.ErrInvalidSecretName):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid secret name"})
	case errors.Is(err, service.ErrSecretNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "secret not found"})
	case errors.Is(err, service.ErrSecretNameTaken):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "secret name already taken"})
	case errors.Is(err, service.ErrSecretAlreadyAttached):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "secret already attached to this project"})
	case errors.Is(err, service.ErrSecretValueNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "secret value not found"})
	case errors.Is(err, service.ErrEnvironmentNotInProject):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "environment is not linked to this project"})
	case errors.Is(err, service.ErrEnvironmentParamRequired):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "environment query parameter is required"})
	case errors.Is(err, service.ErrProductionDownloadForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "only super_admin can download production environment"})
	case errors.Is(err, service.ErrProductionEnvironmentProtected):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "production environment cannot be modified or deleted"})
	default:
		log.Printf("unhandled error method=%s path=%s: %v", c.Method(), c.Path(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
