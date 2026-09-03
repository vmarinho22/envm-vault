package router

import (
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func registerUserRoutes(v1 fiber.Router, deps Deps) {
	g := v1.Group("/auth")
	g.Post("/register",
		deps.JWT.Protect(),
		middleware.RequireAdmin(),
		deps.User.Register,
	)
}
