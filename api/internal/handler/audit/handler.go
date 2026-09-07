package audit

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/handler/httperr"
	"github.com/envm-vault/envm-vault/api/internal/middleware"
	auditservice "github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct {
	audit *auditservice.Service
}

func New(audit *auditservice.Service) *Handler {
	return &Handler{audit: audit}
}

func (h *Handler) List(c *fiber.Ctx) error {
	filter, err := parseListFilter(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.audit.List(c.UserContext(), middleware.GetUserRole(c), filter)
	if err != nil {
		return httperr.Write(c, err)
	}
	return c.JSON(resp)
}

func parseListFilter(c *fiber.Ctx) (auditservice.ListFilter, error) {
	filter := auditservice.ListFilter{
		Page:     1,
		PageSize: auditservice.DefaultPageSize,
	}

	if raw := strings.TrimSpace(c.Query("page")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			return filter, fmt.Errorf("invalid page")
		}
		filter.Page = n
	}

	if raw := strings.TrimSpace(c.Query("page_size")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			return filter, fmt.Errorf("invalid page_size")
		}
		filter.PageSize = n
	}

	sort := strings.ToLower(strings.TrimSpace(c.Query("sort", "desc")))
	switch sort {
	case "asc":
		filter.SortAsc = true
	case "desc", "":
		filter.SortAsc = false
	default:
		return filter, fmt.Errorf("invalid sort; allowed: asc, desc")
	}

	actorIDs, err := parseUUIDList(c, "actor_user_id")
	if err != nil {
		return filter, fmt.Errorf("invalid actor_user_id")
	}
	filter.ActorUserIDs = actorIDs

	targetIDs, err := parseUUIDList(c, "target_user_id")
	if err != nil {
		return filter, fmt.Errorf("invalid target_user_id")
	}
	filter.TargetUserIDs = targetIDs

	from, err := optionalTime(c.Query("occurred_from"), false)
	if err != nil {
		return filter, fmt.Errorf("invalid occurred_from")
	}
	filter.OccurredFrom = from

	to, err := optionalTime(c.Query("occurred_to"), true)
	if err != nil {
		return filter, fmt.Errorf("invalid occurred_to")
	}
	filter.OccurredTo = to

	return filter, nil
}

// parseUUIDList accepts repeated query keys and/or comma-separated values.
func parseUUIDList(c *fiber.Ctx, key string) ([]uuid.UUID, error) {
	rawValues := c.Context().QueryArgs().PeekMulti(key)
	parts := make([]string, 0, len(rawValues))
	for _, v := range rawValues {
		parts = append(parts, string(v))
	}
	if len(parts) == 0 {
		if single := strings.TrimSpace(c.Query(key)); single != "" {
			parts = append(parts, single)
		}
	}

	ids := make([]uuid.UUID, 0)
	seen := make(map[uuid.UUID]struct{})
	for _, part := range parts {
		for _, token := range strings.Split(part, ",") {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			id, err := uuid.Parse(token)
			if err != nil {
				return nil, err
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func optionalTime(raw string, endOfDay bool) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		utc := t.UTC()
		return &utc, nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, err
	}
	if endOfDay {
		t = t.Add(24*time.Hour - time.Nanosecond)
	}
	utc := t.UTC()
	return &utc, nil
}
