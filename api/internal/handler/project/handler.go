package project

import (
	"strings"

	"github.com/vmarinho22/envm-vault/api/internal/handler/httperr"
	"github.com/vmarinho22/envm-vault/api/internal/middleware"
	projectservice "github.com/vmarinho22/envm-vault/api/internal/service/project"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	project *projectservice.Service
}

func New(project *projectservice.Service) *Handler {
	return &Handler{project: project}
}

type createRequest struct {
	Title          string   `json:"title"`
	TeamIDs        []string `json:"team_ids"`
	EnvironmentIDs []string `json:"environment_ids"`
	TagIDs         []string `json:"tag_ids"`
}

type updateRequest struct {
	Title          *string   `json:"title"`
	TeamIDs        *[]string `json:"team_ids"`
	EnvironmentIDs *[]string `json:"environment_ids"`
	TagIDs         *[]string `json:"tag_ids"`
}

func parseCSV(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func (h *Handler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", projectservice.DefaultPageSize)

	resp, err := h.project.List(c.UserContext(), projectservice.ListFilter{
		TitleQuery: c.Query("title"),
		TeamIDs:    parseCSV(c.Query("team_ids")),
		TagIDs:     parseCSV(c.Query("tag_ids")),
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(resp)
}

func (h *Handler) GetBySlug(c *fiber.Ctx) error {
	project, err := h.project.GetBySlug(c.UserContext(), c.Params("slug"))
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(project)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	project, err := h.project.Create(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		req.Title,
		req.TeamIDs,
		req.EnvironmentIDs,
		req.TagIDs,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(project)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	project, err := h.project.Update(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
		req.Title,
		req.TeamIDs,
		req.EnvironmentIDs,
		req.TagIDs,
	)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(project)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	if err := h.project.Delete(
		c.UserContext(),
		middleware.GetUserID(c),
		middleware.GetUserEmail(c),
		middleware.GetUserRole(c),
		c.Params("id"),
	); err != nil {
		return httperr.Write(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
