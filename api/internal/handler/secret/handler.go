package secret

import (
	"strings"

	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	secretservice "github.com/envm-vault/envm-vault/api/internal/service/secret"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	secret *secretservice.Service
}

func New(secret *secretservice.Service) *Handler {
	return &Handler{secret: secret}
}

type createRequest struct {
	Name   string            `json:"name"`
	Values map[string]string `json:"values"`
}

type updateNameRequest struct {
	Name string `json:"name"`
}

type setValueRequest struct {
	Value string `json:"value"`
}

func (h *Handler) List(c *fiber.Ctx) error {
	resp, err := h.secret.List(c.UserContext(), c.Params("id"))
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(resp)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	resp, err := h.secret.Create(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		req.Name,
		req.Values,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *Handler) UpdateName(c *fiber.Ctx) error {
	var req updateNameRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	resp, err := h.secret.UpdateName(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		c.Params("secretId"),
		req.Name,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(resp)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.secret.Delete(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		c.Params("secretId"),
	); err != nil {
		return httperr.Write(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) SetValue(c *fiber.Ctx) error {
	var req setValueRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	resp, err := h.secret.SetValue(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		c.Params("secretId"),
		c.Params("environmentId"),
		req.Value,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(resp)
}

func (h *Handler) RevealValue(c *fiber.Ctx) error {
	resp, err := h.secret.RevealValue(
		c.UserContext(),
		c.Params("id"),
		c.Params("secretId"),
		c.Params("environmentId"),
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(resp)
}

func (h *Handler) ClearValue(c *fiber.Ctx) error {
	if err := h.secret.ClearValue(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		c.Params("secretId"),
		c.Params("environmentId"),
	); err != nil {
		return httperr.Write(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ExportEnvironment(c *fiber.Ctx) error {
	resp, err := h.secret.ExportEnvironment(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		c.Params("environmentId"),
		true,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return writeDownload(c, resp)
}

func (h *Handler) DownloadBySlug(c *fiber.Ctx) error {
	resp, err := h.secret.DownloadBySlug(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("slug"),
		c.Params("envTitle"),
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return writeDownload(c, resp)
}

func writeDownload(c *fiber.Ctx, resp *secretservice.DownloadResponse) error {
	format := strings.ToLower(strings.TrimSpace(c.Query("format")))
	accept := strings.ToLower(c.Get("Accept"))
	asEnv := format == "env" || strings.Contains(accept, "text/plain")
	if asEnv {
		c.Set("Content-Type", "text/plain; charset=utf-8")
		return c.SendString(secretservice.FormatDotEnv(resp.Items))
	}
	return c.JSON(resp)
}
