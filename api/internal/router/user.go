package router

import (
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func registerUserRoutes(v1 fiber.Router, deps Deps) {
	auth := v1.Group("/auth")
	auth.Post("/register",
		deps.JWT.Protect(),
		middleware.RequireAdmin(),
		deps.User.Register,
	)

	users := v1.Group("/users", deps.JWT.Protect(), middleware.RequireAdmin())
	users.Get("", deps.User.List)
	users.Delete("/:id/token", deps.Token.RevokeForUser)
}
