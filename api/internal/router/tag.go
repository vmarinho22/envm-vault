package router

import (
	"github.com/gofiber/fiber/v2"
)

func registerTagRoutes(v1 fiber.Router, deps Deps) {
	tags := v1.Group("/tags", deps.JWT.Protect())
	tags.Get("", deps.Tag.List)
	tags.Post("", deps.Tag.Create)
	tags.Patch("/:id", deps.Tag.Update)
	tags.Delete("/:id", deps.Tag.Delete)
}
