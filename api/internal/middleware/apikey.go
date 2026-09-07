package middleware

import (
	"errors"
	"strings"

	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	auditservice "github.com/envm-vault/envm-vault/api/internal/service/audit"
	tokenservice "github.com/envm-vault/envm-vault/api/internal/service/token"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// APIKeyMiddleware authenticates requests with a personal access token (ev_…).
type APIKeyMiddleware struct {
	tokens *repository.AccessTokenRepository
	users  *repository.UserRepository
	audit  auditservice.Logger
}

func NewAPIKeyMiddleware(
	tokens *repository.AccessTokenRepository,
	users *repository.UserRepository,
	audit auditservice.Logger,
) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		tokens: tokens,
		users:  users,
		audit:  auditservice.OrNop(audit),
	}
}

func (m *APIKeyMiddleware) Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization header",
			})
		}

		plaintext := strings.TrimSpace(parts[1])
		if plaintext == "" || !strings.HasPrefix(plaintext, "ev_") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or revoked access token",
			})
		}

		token, err := m.tokens.FindByHash(c.UserContext(), tokenservice.HashToken(plaintext))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid or revoked access token",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		if token.RevokedAt != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or revoked access token",
			})
		}

		user, err := m.users.FindByID(c.UserContext(), token.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid or revoked access token",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
		if service.IsBlocked(user) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "user is blocked",
			})
		}

		c.Locals(string(UserRoleKey), user.Role)
		c.Locals(string(UserIDKey), user.ID.String())
		c.Locals(string(UserEmailKey), user.Email)
		c.SetUserContext(m.audit.WithChannel(c.UserContext(), auditservice.ChannelPersonalToken))

		return c.Next()
	}
}
