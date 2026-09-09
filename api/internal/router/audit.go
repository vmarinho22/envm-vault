package router

import (
	"github.com/vmarinho22/envm-vault/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func registerAuditRoutes(v1 fiber.Router, deps Deps) {
	g := v1.Group("/audit-logs", deps.JWT.Protect(), middleware.RequireAdmin())
	g.Get("", deps.Audit.List)
}
