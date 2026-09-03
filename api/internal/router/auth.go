package router

import "github.com/gofiber/fiber/v2"

func registerAuthRoutes(v1 fiber.Router, deps Deps) {
	g := v1.Group("/auth")
	g.Post("/login", deps.Auth.Login)
	g.Post("/verify", deps.Auth.Verify)
}
