package router

import "github.com/gofiber/fiber/v2"

func registerHealthRoutes(v1 fiber.Router, deps Deps) {
	v1.Get("/health", deps.Health.Check)
}
