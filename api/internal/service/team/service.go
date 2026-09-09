package team

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
	"github.com/vmarinho22/envm-vault/api/internal/slug"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxTeamNameLen = 80

type Service struct {
	teams *repository.TeamRepository
	users *repository.UserRepository
	audit audit.Logger
}

func New(teams *repository.TeamRepository, users *repository.UserRepository, auditLog audit.Logger) *Service {
	return &Service{teams: teams, users: users, audit: audit.OrNop(auditLog)}
}

type AccountableResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type TeamResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Accountable AccountableResponse `json:"accountable"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

func toResponse(t model.Team) TeamResponse {
	return TeamResponse{
		ID:   t.ID.String(),
		Name: t.Name,
		Slug: t.Slug,
		Accountable: AccountableResponse{
			ID:    t.Accountable.ID.String(),
			Email: t.Accountable.Email,
		},
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func requireAdmin(role string) error {
	if role != service.RoleSuperAdmin && role != service.RoleAdmin {
		return service.ErrForbidden
	}
	return nil
}

func normalizeName(raw string) (string, string, error) {
	name := strings.TrimSpace(raw)
	if name == "" || utf8.RuneCountInString(name) > maxTeamNameLen {
		return "", "", service.ErrInvalidTeamName
	}
	s := slug.FromName(name)
	if s == "" {
		return "", "", service.ErrInvalidTeamName
	}
	return name, s, nil
}

func parseActorID(actorID string) *uuid.UUID {
	id, err := uuid.Parse(actorID)
	if err != nil {
		return nil
	}
	return &id
}

func (s *Service) List(ctx context.Context, nameQuery string) ([]TeamResponse, error) {
	teams, err := s.teams.List(ctx, nameQuery)
	if err != nil {
		return nil, err
	}
	out := make([]TeamResponse, 0, len(teams))
	for _, t := range teams {
		out = append(out, toResponse(t))
	}
	return out, nil
}

func (s *Service) Create(ctx context.Context, actorID, actorEmail, actorRole, name, accountableID string) (*TeamResponse, error) {
	if err := requireAdmin(actorRole); err != nil {
		return nil, err
	}

	normalizedName, teamSlug, err := normalizeName(name)
	if err != nil {
		return nil, err
	}

	accID, err := uuid.Parse(strings.TrimSpace(accountableID))
	if err != nil {
		return nil, service.ErrUserNotFound
	}

	if _, err := s.users.FindByID(ctx, accID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}

	taken, err := s.teams.ExistsBySlug(ctx, teamSlug, nil)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, service.ErrTeamSlugTaken
	}

	created, err := s.teams.Create(ctx, &model.Team{
		Name:          normalizedName,
		Slug:          teamSlug,
		AccountableID: accID,
	})
	if err != nil {
		return nil, err
	}

	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionTeamCreated,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		TargetUserID: &accID,
	})

	resp := toResponse(*created)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, actorID, actorEmail, actorRole, teamID string, name *string, accountableID *string) (*TeamResponse, error) {
	if err := requireAdmin(actorRole); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(strings.TrimSpace(teamID))
	if err != nil {
		return nil, service.ErrTeamNotFound
	}

	team, err := s.teams.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrTeamNotFound
		}
		return nil, err
	}

	if name == nil && accountableID == nil {
		resp := toResponse(*team)
		return &resp, nil
	}

	if name != nil {
		normalizedName, teamSlug, nameErr := normalizeName(*name)
		if nameErr != nil {
			return nil, nameErr
		}
		if teamSlug != team.Slug {
			taken, existsErr := s.teams.ExistsBySlug(ctx, teamSlug, &id)
			if existsErr != nil {
				return nil, existsErr
			}
			if taken {
				return nil, service.ErrTeamSlugTaken
			}
		}
		team.Name = normalizedName
		team.Slug = teamSlug
	}

	if accountableID != nil {
		accID, parseErr := uuid.Parse(strings.TrimSpace(*accountableID))
		if parseErr != nil {
			return nil, service.ErrUserNotFound
		}
		if _, findErr := s.users.FindByID(ctx, accID); findErr != nil {
			if errors.Is(findErr, gorm.ErrRecordNotFound) {
				return nil, service.ErrUserNotFound
			}
			return nil, findErr
		}
		team.AccountableID = accID
	}

	updated, err := s.teams.Update(ctx, team)
	if err != nil {
		return nil, err
	}

	target := updated.AccountableID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionTeamUpdated,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		TargetUserID: &target,
	})

	resp := toResponse(*updated)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, actorID, actorEmail, actorRole, teamID string) error {
	if err := requireAdmin(actorRole); err != nil {
		return err
	}

	id, err := uuid.Parse(strings.TrimSpace(teamID))
	if err != nil {
		return service.ErrTeamNotFound
	}

	team, err := s.teams.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrTeamNotFound
		}
		return err
	}

	if err := s.teams.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrTeamNotFound
		}
		return err
	}

	target := team.AccountableID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionTeamDeleted,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		TargetUserID: &target,
	})
	return nil
}
