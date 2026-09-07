package audit

import (
	"context"
	"time"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/google/uuid"
)

const DefaultPageSize = 25

// ListFilter is the domain filter for listing audit logs.
type ListFilter struct {
	Page          int
	PageSize      int
	SortAsc       bool
	ActorUserIDs  []uuid.UUID
	TargetUserIDs []uuid.UUID
	OccurredFrom  *time.Time
	OccurredTo    *time.Time
}

// LogResponse is one audit log row for the API.
type LogResponse struct {
	ID           string     `json:"id"`
	Action       string     `json:"action"`
	ActorUserID  *string    `json:"actor_user_id"`
	ActorEmail   string     `json:"actor_email"`
	ActorRole    string     `json:"actor_role"`
	TargetUserID *string    `json:"target_user_id"`
	TargetEmail  *string    `json:"target_email"`
	AuthChannel  string     `json:"auth_channel"`
	IP           string     `json:"ip"`
	UserAgent    string     `json:"user_agent"`
	Origin       string     `json:"origin"`
	OccurredAt   time.Time  `json:"occurred_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ListResponse is a paginated audit log page.
type ListResponse struct {
	Items    []LogResponse `json:"items"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
	Total    int64         `json:"total"`
}

func (s *Service) List(ctx context.Context, actorRole string, filter ListFilter) (*ListResponse, error) {
	if actorRole != service.RoleSuperAdmin && actorRole != service.RoleAdmin {
		return nil, service.ErrForbidden
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > DefaultPageSize {
		pageSize = DefaultPageSize
	}

	rows, total, err := s.logs.List(ctx, repository.AuditLogListFilter{
		Page:          page,
		PageSize:      pageSize,
		SortAsc:       filter.SortAsc,
		ActorUserIDs:  filter.ActorUserIDs,
		TargetUserIDs: filter.TargetUserIDs,
		OccurredFrom:  filter.OccurredFrom,
		OccurredTo:    filter.OccurredTo,
	})
	if err != nil {
		return nil, err
	}

	emails, err := s.targetEmails(ctx, rows)
	if err != nil {
		return nil, err
	}

	items := make([]LogResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, toLogResponse(row, emails))
	}

	return &ListResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *Service) targetEmails(ctx context.Context, rows []model.AuditLog) (map[uuid.UUID]string, error) {
	if s.users == nil {
		return map[uuid.UUID]string{}, nil
	}

	seen := make(map[uuid.UUID]struct{})
	ids := make([]uuid.UUID, 0)
	for _, row := range rows {
		if row.TargetUserID == nil {
			continue
		}
		id := *row.TargetUserID
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return map[uuid.UUID]string{}, nil
	}

	users, err := s.users.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	out := make(map[uuid.UUID]string, len(users))
	for _, u := range users {
		out[u.ID] = u.Email
	}
	return out, nil
}

func toLogResponse(row model.AuditLog, emails map[uuid.UUID]string) LogResponse {
	resp := LogResponse{
		ID:          row.ID.String(),
		Action:      row.Action,
		ActorEmail:  row.ActorEmail,
		ActorRole:   row.ActorRole,
		AuthChannel: row.AuthChannel,
		IP:          row.IP,
		UserAgent:   row.UserAgent,
		Origin:      row.Origin,
		OccurredAt:  row.OccurredAt,
		CreatedAt:   row.CreatedAt,
	}
	if row.ActorUserID != nil {
		id := row.ActorUserID.String()
		resp.ActorUserID = &id
	}
	if row.TargetUserID != nil {
		id := row.TargetUserID.String()
		resp.TargetUserID = &id
		if email, ok := emails[*row.TargetUserID]; ok {
			resp.TargetEmail = &email
		}
	}
	return resp
}
