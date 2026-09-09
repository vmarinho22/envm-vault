package middleware

import (
	auditservice "github.com/vmarinho22/envm-vault/api/internal/service/audit"
	"github.com/gofiber/fiber/v2"
)

// NewAudit attaches request origin (IP, User-Agent, Origin, auth channel) to the request context.
func NewAudit(audit auditservice.Logger) fiber.Handler {
	audit = auditservice.OrNop(audit)
	return func(c *fiber.Ctx) error {
		ctx := audit.AttachOrigin(c.UserContext(), auditservice.RequestOrigin{
			IP:          c.IP(),
			UserAgent:   c.Get("User-Agent"),
			Origin:      c.Get("Origin"),
			AuthChannel: auditservice.ChannelLoginAttempt,
		})
		c.SetUserContext(ctx)
		return c.Next()
	}
}
