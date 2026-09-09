package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// NewRateLimiter returns a per-IP rate limiter for API routes (/v1/*).
// Health probes and non-API paths (SPA/static) are skipped.
func NewRateLimiter(max int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: window,
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			if path == "/health" || path == "/v1/health" {
				return true
			}
			return !strings.HasPrefix(path, "/v1")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "rate limit exceeded",
			})
		},
	})
}
