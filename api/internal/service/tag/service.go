package tag

import (
	"context"
	"errors"
	"regexp"
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

const maxTagTitleLen = 80

var hexColorPattern = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type Service struct {
	tags  *repository.TagRepository
	audit audit.Logger
}

func New(tags *repository.TagRepository, auditLog audit.Logger) *Service {
	return &Service{tags: tags, audit: audit.OrNop(auditLog)}
}

type TagResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(t model.Tag) TagResponse {
	return TagResponse{
		ID:        t.ID.String(),
		Title:     t.Title,
		Color:     t.Color,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func normalizeTitle(raw string) (string, error) {
	title := strings.TrimSpace(raw)
	if title == "" || utf8.RuneCountInString(title) > maxTagTitleLen {
		return "", service.ErrInvalidTagTitle
	}
	return title, nil
}

func normalizeColor(raw string) (string, error) {
	color := strings.TrimSpace(raw)
	if !hexColorPattern.MatchString(color) {
		return "", service.ErrInvalidTagColor
	}
	return strings.ToUpper(color), nil
}

func parseActorID(actorID string) *uuid.UUID {
	id, err := uuid.Parse(actorID)
	if err != nil {
		return nil
	}
	return &id
}

func (s *Service) List(ctx context.Context, titleQuery string) ([]TagResponse, error) {
	tags, err := s.tags.List(ctx, titleQuery)
	if err != nil {
		return nil, err
	}
	out := make([]TagResponse, 0, len(tags))
	for _, t := range tags {
		out = append(out, toResponse(t))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, actorID, actorEmail, actorRole, title, color string) (*TagResponse, error) {
	normalizedTitle, err := normalizeTitle(title)
	if err != nil {
		return nil, err
	}

	normalizedColor, err := normalizeColor(color)
	if err != nil {
		return nil, err
	}

	taken, err := s.tags.ExistsByTitle(ctx, normalizedTitle, nil)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, service.ErrTagTitleTaken
	}

	created, err := s.tags.Create(ctx, &model.Tag{
		Title: normalizedTitle,
		Color: normalizedColor,
	})
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionTagCreated,
		ActorUserID: parseActorID(actorID),
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})

	resp := toResponse(*created)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, actorID, actorEmail, actorRole, tagID string, title *string, color *string) (*TagResponse, error) {
	id, err := uuid.Parse(strings.TrimSpace(tagID))
	if err != nil {
		return nil, service.ErrTagNotFound
	}

	tag, err := s.tags.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrTagNotFound
		}
		return nil, err
	}

	if title == nil && color == nil {
		resp := toResponse(*tag)
		return &resp, nil
	}

	if title != nil {
		normalizedTitle, titleErr := normalizeTitle(*title)
		if titleErr != nil {
			return nil, titleErr
		}
		if !strings.EqualFold(normalizedTitle, tag.Title) {
			taken, existsErr := s.tags.ExistsByTitle(ctx, normalizedTitle, &id)
			if existsErr != nil {
				return nil, existsErr
			}
			if taken {
				return nil, service.ErrTagTitleTaken
			}
		}
		tag.Title = normalizedTitle
	}

	if color != nil {
		normalizedColor, colorErr := normalizeColor(*color)
		if colorErr != nil {
			return nil, colorErr
		}
		tag.Color = normalizedColor
	}

	updated, err := s.tags.Update(ctx, tag)
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionTagUpdated,
		ActorUserID: parseActorID(actorID),
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})

	resp := toResponse(*updated)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, actorID, actorEmail, actorRole, tagID string) error {
	id, err := uuid.Parse(strings.TrimSpace(tagID))
	if err != nil {
		return service.ErrTagNotFound
	}

	if _, err := s.tags.FindByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrTagNotFound
		}
		return err
	}

	if err := s.tags.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrTagNotFound
		}
		return err
	}

	s.audit.Record(ctx, audit.Event{
		Action:      audit.ActionTagDeleted,
		ActorUserID: parseActorID(actorID),
		ActorEmail:  actorEmail,
		ActorRole:   actorRole,
	})
	return nil
}
