package router

import (
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func registerProjectRoutes(v1 fiber.Router, deps Deps) {
	projects := v1.Group("/projects", deps.JWT.Protect())
	projects.Get("", deps.Project.List)
	projects.Get("/by-slug/:slug", deps.Project.GetBySlug)
	projects.Post("", deps.Project.Create)
	projects.Patch("/:id", deps.Project.Update)
	projects.Delete("/:id", middleware.RequireAdmin(), deps.Project.Delete)
}
