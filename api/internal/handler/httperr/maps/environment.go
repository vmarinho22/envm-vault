package maps

import (
	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapEnvironment(err error) *Mapped {
	return firstMatch(
		match(err, service.ErrInvalidEnvironmentTitle, fiber.StatusBadRequest, "invalid environment title"),
		match(err, service.ErrEnvironmentNotFound, fiber.StatusNotFound, "environment not found"),
		match(err, service.ErrEnvironmentTitleTaken, fiber.StatusConflict, "environment title already taken"),
		match(err, service.ErrProductionEnvironmentProtected, fiber.StatusForbidden, "production environment cannot be modified or deleted"),
	)
}
