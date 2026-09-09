package maps

import (
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapTeam(err error) *Mapped {
	return firstMatch(
		match(err, service.ErrInvalidTeamName, fiber.StatusBadRequest, "invalid team name"),
		match(err, service.ErrTeamNotFound, fiber.StatusNotFound, "team not found"),
		match(err, service.ErrTeamSlugTaken, fiber.StatusConflict, "team slug already taken"),
	)
}
