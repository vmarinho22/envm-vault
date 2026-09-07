package router

import (
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func registerTeamRoutes(v1 fiber.Router, deps Deps) {
	teams := v1.Group("/teams", deps.JWT.Protect())
	teams.Get("", deps.Team.List)
	teams.Post("", middleware.RequireAdmin(), deps.Team.Create)
	teams.Patch("/:id", middleware.RequireAdmin(), deps.Team.Update)
	teams.Delete("/:id", middleware.RequireAdmin(), deps.Team.Delete)
}
