package router

import "github.com/gofiber/fiber/v2"

func registerTokenRoutes(v1 fiber.Router, deps Deps) {
	g := v1.Group("/tokens", deps.JWT.Protect())
	g.Get("/me", deps.Token.GetMine)
	g.Post("", deps.Token.Create)
	g.Delete("", deps.Token.RevokeOwn)
}
