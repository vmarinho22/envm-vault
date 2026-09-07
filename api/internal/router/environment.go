package router

import (
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func registerEnvironmentRoutes(v1 fiber.Router, deps Deps) {
	environments := v1.Group("/environments", deps.JWT.Protect())
	environments.Get("", deps.Environment.List)
	environments.Post("", middleware.RequireAdmin(), deps.Environment.Create)
	environments.Patch("/:id", middleware.RequireAdmin(), deps.Environment.Update)
	environments.Delete("/:id", middleware.RequireAdmin(), deps.Environment.Delete)
}
