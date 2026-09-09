package router

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MountFrontend serves the Vue SPA from distDir when non-empty.
// Call after Register so /v1 and /health keep precedence.
func MountFrontend(app *fiber.App, distDir string) {
	distDir = strings.TrimSpace(distDir)
	if distDir == "" {
		return
	}
	info, err := os.Stat(distDir)
	if err != nil || !info.IsDir() {
		return
	}

	index := filepath.Join(distDir, "index.html")
	app.Static("/", distDir, fiber.Static{
		Index: "index.html",
	})
	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendFile(index)
	})
}
