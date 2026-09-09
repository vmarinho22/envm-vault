package router

import (
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func registerProjectRoutes(v1 fiber.Router, deps Deps) {
	// API-key download must be registered before the JWT /projects group.
	// Fiber runs matching USE middleware (Group handlers) first when they appear
	// earlier in the stack; a Group("/projects", JWT) would otherwise reject
	// Bearer ev_… tokens with "invalid or expired token". Use route-level
	// APIKey (not Group Use) so JWT project routes are not forced through it.
	v1.Get("/projects/:slug/env", deps.APIKey.Protect(), deps.Secret.DownloadBySlug)

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
}
