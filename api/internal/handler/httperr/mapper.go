package httperr

import (
	"errors"
	"log"

	"github.com/envm-vault/envm-vault/api/internal/handler/httperr/maps"
	"github.com/gofiber/fiber/v2"
)

func AppHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	var ferr *fiber.Error
	if errors.As(err, &ferr) {
		code = ferr.Code
		if code < fiber.StatusInternalServerError && ferr.Message != "" {
			msg = ferr.Message
		}
	}

	if code >= fiber.StatusInternalServerError {
		log.Printf("request error method=%s path=%s: %v", c.Method(), c.Path(), err)
	}

	return c.Status(code).JSON(fiber.Map{"error": msg})
}

func Write(c *fiber.Ctx, err error) error {
	if mapped := maps.Lookup(err); mapped != nil {
		return c.Status(mapped.Status).JSON(fiber.Map{"error": mapped.Message})
	}

	log.Printf("unhandled error method=%s path=%s: %v", c.Method(), c.Path(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}
