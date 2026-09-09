package environment

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/vmarinho22/envm-vault/api/internal/model"
	"github.com/vmarinho22/envm-vault/api/internal/repository"
	"github.com/vmarinho22/envm-vault/api/internal/service"
	"github.com/vmarinho22/envm-vault/api/internal/service/audit"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxEnvironmentTitleLen = 80

type Service struct {
	environments *repository.EnvironmentRepository
	audit        audit.Logger
}

func New(environments *repository.EnvironmentRepository, auditLog audit.Logger) *Service {
	return &Service{environments: environments, audit: audit.OrNop(auditLog)}
}

type EnvironmentResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(e model.Environment) EnvironmentResponse {
	return EnvironmentResponse{
		ID:        e.ID.String(),
		Title:     e.Title,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func requireAdmin(role string) error {
	if role != service.RoleSuperAdmin && role != service.RoleAdmin {
		return service.ErrForbidden
	}
	return nil
}

func AssertMutableEnvironment(title string) error {
	if title == service.ProductionEnvironmentTitle {
		return service.ErrProductionEnvironmentProtected
	}
	return nil
}

func normalizeTitle(raw string) (string, error) {
	title := strings.TrimSpace(raw)
	if title == "" || utf8.RuneCountInString(title) > maxEnvironmentTitleLen {
		return "", service.ErrInvalidEnvironmentTitle
	}
	return title, nil
}

func parseActorID(actorID string) *uuid.UUID {
	id, err := uuid.Parse(actorID)
	if err != nil {
		return nil
	}
	return &id
}

func (s *Service) List(ctx context.Context, titleQuery string) ([]EnvironmentResponse, error) {
	environments, err := s.environments.List(ctx, titleQuery)
	if err != nil {
		return nil, err
	}
	out := make([]EnvironmentResponse, 0, len(environments))
	for _, e := range environments {
		out = append(out, toResponse(e))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, actorID, actorEmail, actorRole, title string) (*EnvironmentResponse, error) {
	if err := requireAdmin(actorRole); err != nil {
		return nil, err
	}

	normalizedTitle, err := normalizeTitle(title)
	if err != nil {
		return nil, err
	}

	taken, err := s.environments.ExistsByTitle(ctx, normalizedTitle, nil)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, service.ErrEnvironmentTitleTaken
	}

	created, err := s.environments.Create(ctx, &model.Environment{
		Title: normalizedTitle,
	})
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionEnvironmentCreated,
		ActorUserID: parseActorID(actorID),
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})

	resp := toResponse(*created)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, actorID, actorEmail, actorRole, environmentID string, title *string) (*EnvironmentResponse, error) {
	if err := requireAdmin(actorRole); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(strings.TrimSpace(environmentID))
	if err != nil {
		return nil, service.ErrEnvironmentNotFound
	}

	environment, err := s.environments.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrEnvironmentNotFound
		}
		return nil, err
	}

	if title == nil {
		resp := toResponse(*environment)
		return &resp, nil
	}

	if err := AssertMutableEnvironment(environment.Title); err != nil {
		return nil, err
	}

	normalizedTitle, titleErr := normalizeTitle(*title)
	if titleErr != nil {
		return nil, titleErr
	}
	if !strings.EqualFold(normalizedTitle, environment.Title) {
		taken, existsErr := s.environments.ExistsByTitle(ctx, normalizedTitle, &id)
		if existsErr != nil {
			return nil, existsErr
		}
		if taken {
			return nil, service.ErrEnvironmentTitleTaken
		}
	}
	environment.Title = normalizedTitle

	updated, err := s.environments.Update(ctx, environment)
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionEnvironmentUpdated,
		ActorUserID: parseActorID(actorID),
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})

	resp := toResponse(*updated)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, actorID, actorEmail, actorRole, environmentID string) error {
	if err := requireAdmin(actorRole); err != nil {
		return err
	}

	id, err := uuid.Parse(strings.TrimSpace(environmentID))
	if err != nil {
		return service.ErrEnvironmentNotFound
	}

	environment, err := s.environments.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrEnvironmentNotFound
		}
		return err
	}
	if err := AssertMutableEnvironment(environment.Title); err != nil {
		return err
	}

	if err := s.environments.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrEnvironmentNotFound
		}
		return err
	}

	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionEnvironmentDeleted,
		ActorUserID: parseActorID(actorID),
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})
	return nil
}
