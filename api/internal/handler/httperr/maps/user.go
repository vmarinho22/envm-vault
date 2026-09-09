package maps

import (
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapUser(err error) *Mapped {
	return firstMatch(
		match(err, service.ErrUserNotFound, fiber.StatusNotFound, "user not found"),
		match(err, service.ErrUserBlocked, fiber.StatusForbidden, "user is blocked"),
	)
}
