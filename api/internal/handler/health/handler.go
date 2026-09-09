package health

import (
	"github.com/vmarinho22/envm-vault/api/internal/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Check(c *fiber.Ctx) error {
	dbStatus := "ok"

	if err := database.Ping(h.db); err != nil {
		dbStatus = "error"
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "degraded",
			"db":     dbStatus,
		})
	}

	return c.JSON(fiber.Map{
		"status": "ok",
		"db":     dbStatus,
	})
}
