package maps

import (
	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapProject(err error) *Mapped {
	return firstMatch(
		match(err, service.ErrInvalidProjectTitle, fiber.StatusBadRequest, "invalid project title"),
		match(err, service.ErrProjectNotFound, fiber.StatusNotFound, "project not found"),
		match(err, service.ErrProjectSlugTaken, fiber.StatusConflict, "project slug already taken"),
		match(err, service.ErrInvalidProjectRelations, fiber.StatusBadRequest, "project requires at least one team and one environment"),
	)
}
