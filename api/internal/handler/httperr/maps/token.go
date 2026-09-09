package maps

import (
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
)

func mapToken(err error) *Mapped {
	return match(err, service.ErrTokenNotFound, fiber.StatusNotFound, "access token not found")
}
