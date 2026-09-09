package maps

import (
	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapTag(err error) *Mapped {
	return firstMatch(
		match(err, service.ErrInvalidTagTitle, fiber.StatusBadRequest, "invalid tag title"),
		match(err, service.ErrInvalidTagColor, fiber.StatusBadRequest, "invalid tag color"),
		match(err, service.ErrTagNotFound, fiber.StatusNotFound, "tag not found"),
		match(err, service.ErrTagTitleTaken, fiber.StatusConflict, "tag title already taken"),
	)
}
