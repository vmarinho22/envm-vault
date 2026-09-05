package middleware

import (
	"strings"

	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserRoleKey contextKey = "userRole"
const UserIDKey contextKey = "userID"
const UserEmailKey contextKey = "userEmail"

type JWTMiddleware struct {
	secret []byte
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{secret: []byte(secret)}
}

func (m *JWTMiddleware) Protect() fiber.Handler {
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

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return m.secret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		role, _ := claims["role"].(string)
		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)

		c.Locals(string(UserRoleKey), role)
		c.Locals(string(UserIDKey), userID)
		c.Locals(string(UserEmailKey), email)

		return c.Next()
	}
}

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals(string(UserRoleKey)).(string)
		if role != service.RoleSuperAdmin && role != service.RoleAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "insufficient permissions",
			})
		}
		return c.Next()
	}
}

func GetUserRole(c *fiber.Ctx) string {
	role, _ := c.Locals(string(UserRoleKey)).(string)
	return role
}

func GetUserID(c *fiber.Ctx) string {
	id, _ := c.Locals(string(UserIDKey)).(string)
	return id
}
