package router

import (
	"github.com/gofiber/fiber/v2"
)

func registerSecretRoutes(v1 fiber.Router, deps Deps) {
	secrets := v1.Group("/secrets", deps.JWT.Protect())
	secrets.Get("", deps.Secret.SearchCatalog)
}
