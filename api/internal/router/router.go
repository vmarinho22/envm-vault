package router

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App, deps Deps) {
	v1 := app.Group("/v1")
	registerHealthRoutes(v1, deps)
	registerAuthRoutes(v1, deps)
	registerUserRoutes(v1, deps)
	registerTokenRoutes(v1, deps)
	registerAuditRoutes(v1, deps)
	registerTeamRoutes(v1, deps)
	registerTagRoutes(v1, deps)
	registerEnvironmentRoutes(v1, deps)
	registerProjectRoutes(v1, deps)
	registerSecretRoutes(v1, deps)

	app.Get("/health", deps.Health.Check)
}
