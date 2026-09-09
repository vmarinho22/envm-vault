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

	projects.Get("/:id/secrets", deps.Secret.List)
	projects.Post("/:id/secrets", deps.Secret.Create)
	projects.Post("/:id/secrets/attach", deps.Secret.Attach)
	projects.Patch("/:id/secrets/:secretId", deps.Secret.UpdateName)
	projects.Delete("/:id/secrets/:secretId", deps.Secret.Delete)
	projects.Put("/:id/secrets/:secretId/values/:environmentId", deps.Secret.SetValue)
	projects.Get("/:id/secrets/:secretId/values/:environmentId", deps.Secret.RevealValue)
	projects.Delete("/:id/secrets/:secretId/values/:environmentId", deps.Secret.ClearValue)
	projects.Get("/:id/environments/:environmentId/values", deps.Secret.ExportEnvironment)

	// API-key download (personal token). Register outside JWT group.
	download := v1.Group("/projects", deps.APIKey.Protect())
	download.Get("/:slug/env", deps.Secret.DownloadBySlug)
}
