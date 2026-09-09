package maps

import (
	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapSecret(err error) *Mapped {
	return firstMatch(
		match(err, service.ErrInvalidSecretName, fiber.StatusBadRequest, "invalid secret name"),
		match(err, service.ErrSecretNotFound, fiber.StatusNotFound, "secret not found"),
		match(err, service.ErrSecretNameTaken, fiber.StatusConflict, "secret name already taken"),
		match(err, service.ErrSecretAlreadyAttached, fiber.StatusConflict, "secret already attached to this project"),
		match(err, service.ErrSecretValueNotFound, fiber.StatusNotFound, "secret value not found"),
		match(err, service.ErrEnvironmentNotInProject, fiber.StatusBadRequest, "environment is not linked to this project"),
		match(err, service.ErrEnvironmentParamRequired, fiber.StatusBadRequest, "environment query parameter is required"),
		match(err, service.ErrProductionDownloadForbidden, fiber.StatusForbidden, "only super_admin can download production environment"),
		match(err, service.ErrInvalidSecretPrefix, fiber.StatusBadRequest, "invalid secret prefix"),
		match(err, service.ErrDownloadNameConflict, fiber.StatusConflict, "download names conflict after applying prefixes"),
	)
}
