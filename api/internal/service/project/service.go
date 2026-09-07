package project

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/envm-vault/envm-vault/api/internal/slug"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxProjectTitleLen = 80
	DefaultPageSize    = 12
)

type Service struct {
	projects     *repository.ProjectRepository
	teams        *repository.TeamRepository
	environments *repository.EnvironmentRepository
	tags         *repository.TagRepository
	audit        audit.Logger
}

func New(
	projects *repository.ProjectRepository,
	teams *repository.TeamRepository,
	environments *repository.EnvironmentRepository,
	tags *repository.TagRepository,
	auditLog audit.Logger,
) *Service {
	return &Service{
		projects:     projects,
		teams:        teams,
		environments: environments,
		tags:         tags,
		audit:        audit.OrNop(auditLog),
	}
}

type TeamSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type EnvironmentSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TagSummary struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

type ProjectResponse struct {
	ID           string               `json:"id"`
	Title        string               `json:"title"`
	Slug         string               `json:"slug"`
	Teams        []TeamSummary        `json:"teams"`
	Environments []EnvironmentSummary `json:"environments"`
	Tags         []TagSummary         `json:"tags"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

type ListResponse struct {
	Items    []ProjectResponse `json:"items"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int64             `json:"total"`
}

type ListFilter struct {
	TitleQuery string
	TeamIDs    []string
	TagIDs     []string
	Page       int
	PageSize   int
}

func toResponse(p model.Project) ProjectResponse {
	teams := make([]TeamSummary, 0, len(p.Teams))
	for _, t := range p.Teams {
		teams = append(teams, TeamSummary{ID: t.ID.String(), Name: t.Name, Slug: t.Slug})
	}
	envs := make([]EnvironmentSummary, 0, len(p.Environments))
	for _, e := range p.Environments {
		envs = append(envs, EnvironmentSummary{ID: e.ID.String(), Title: e.Title})
	}
	tags := make([]TagSummary, 0, len(p.Tags))
	for _, t := range p.Tags {
		tags = append(tags, TagSummary{ID: t.ID.String(), Title: t.Title, Color: t.Color})
	}
	return ProjectResponse{
		ID:           p.ID.String(),
		Title:        p.Title,
		Slug:         p.Slug,
		Teams:        teams,
		Environments: envs,
		Tags:         tags,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func requireAdmin(role string) error {
	if role != service.RoleSuperAdmin && role != service.RoleAdmin {
		return service.ErrForbidden
	}
	return nil
}

func normalizeTitle(raw string) (string, string, error) {
	title := strings.TrimSpace(raw)
	if title == "" || utf8.RuneCountInString(title) > maxProjectTitleLen {
		return "", "", service.ErrInvalidProjectTitle
	}
	s := slug.FromName(title)
	if s == "" {
		return "", "", service.ErrInvalidProjectTitle
	}
	return title, s, nil
}

func parseActorID(actorID string) *uuid.UUID {
	id, err := uuid.Parse(actorID)
	if err != nil {
		return nil
	}
	return &id
}

func parseUUIDs(raw []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(raw))
	seen := make(map[uuid.UUID]struct{}, len(raw))
	for _, s := range raw {
		id, err := uuid.Parse(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out, nil
}

func (s *Service) resolveRelations(ctx context.Context, teamIDs, environmentIDs, tagIDs []string) ([]model.Team, []model.Environment, []model.Tag, error) {
	if len(teamIDs) == 0 || len(environmentIDs) == 0 {
		return nil, nil, nil, service.ErrInvalidProjectRelations
	}

	tIDs, err := parseUUIDs(teamIDs)
	if err != nil {
		return nil, nil, nil, service.ErrTeamNotFound
	}
	eIDs, err := parseUUIDs(environmentIDs)
	if err != nil {
		return nil, nil, nil, service.ErrEnvironmentNotFound
	}
	var tgIDs []uuid.UUID
	if len(tagIDs) > 0 {
		tgIDs, err = parseUUIDs(tagIDs)
		if err != nil {
			return nil, nil, nil, service.ErrTagNotFound
		}
	}

	teams, err := s.teams.FindByIDs(ctx, tIDs)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(teams) != len(tIDs) {
		return nil, nil, nil, service.ErrTeamNotFound
	}

	envs, err := s.environments.FindByIDs(ctx, eIDs)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(envs) != len(eIDs) {
		return nil, nil, nil, service.ErrEnvironmentNotFound
	}

	tags, err := s.tags.FindByIDs(ctx, tgIDs)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(tags) != len(tgIDs) {
		return nil, nil, nil, service.ErrTagNotFound
	}

	return teams, envs, tags, nil
}

func (s *Service) List(ctx context.Context, filter ListFilter) (*ListResponse, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > DefaultPageSize {
		pageSize = DefaultPageSize
	}

	var teamUUIDs []uuid.UUID
	if len(filter.TeamIDs) > 0 {
		ids, err := parseUUIDs(filter.TeamIDs)
		if err != nil {
			return nil, service.ErrTeamNotFound
		}
		teamUUIDs = ids
	}
	var tagUUIDs []uuid.UUID
	if len(filter.TagIDs) > 0 {
		ids, err := parseUUIDs(filter.TagIDs)
		if err != nil {
			return nil, service.ErrTagNotFound
		}
		tagUUIDs = ids
	}

	projects, total, err := s.projects.List(ctx, repository.ProjectListFilter{
		TitleQuery: filter.TitleQuery,
		TeamIDs:    teamUUIDs,
		TagIDs:     tagUUIDs,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		return nil, err
	}

	items := make([]ProjectResponse, 0, len(projects))
	for _, p := range projects {
		items = append(items, toResponse(p))
	}
	return &ListResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *Service) GetBySlug(ctx context.Context, projectSlug string) (*ProjectResponse, error) {
	project, err := s.projects.FindBySlug(ctx, strings.TrimSpace(projectSlug))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrProjectNotFound
		}
		return nil, err
	}
	resp := toResponse(*project)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, actorID, actorEmail, actorRole, title string, teamIDs, environmentIDs, tagIDs []string) (*ProjectResponse, error) {
	normalizedTitle, projectSlug, err := normalizeTitle(title)
	if err != nil {
		return nil, err
	}

	teams, envs, tags, err := s.resolveRelations(ctx, teamIDs, environmentIDs, tagIDs)
	if err != nil {
		return nil, err
	}

	taken, err := s.projects.ExistsBySlug(ctx, projectSlug, nil)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, service.ErrProjectSlugTaken
	}

	created, err := s.projects.Create(ctx, &model.Project{
		Title:        normalizedTitle,
		Slug:         projectSlug,
		Teams:        teams,
		Environments: envs,
		Tags:         tags,
	})
	if err != nil {
		return nil, err
	}

	projectID := created.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionProjectCreated,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &projectID,
		ProjectTitle: created.Title,
		ProjectSlug:  created.Slug,
	})

	resp := toResponse(*created)
	return &resp, nil
}

func (s *Service) Update(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectID string,
	title *string,
	teamIDs, environmentIDs, tagIDs *[]string,
) (*ProjectResponse, error) {
	id, err := uuid.Parse(strings.TrimSpace(projectID))
	if err != nil {
		return nil, service.ErrProjectNotFound
	}

	project, err := s.projects.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrProjectNotFound
		}
		return nil, err
	}

	if title == nil && teamIDs == nil && environmentIDs == nil && tagIDs == nil {
		resp := toResponse(*project)
		return &resp, nil
	}

	if title != nil {
		normalizedTitle, projectSlug, titleErr := normalizeTitle(*title)
		if titleErr != nil {
			return nil, titleErr
		}
		if projectSlug != project.Slug {
			taken, existsErr := s.projects.ExistsBySlug(ctx, projectSlug, &id)
			if existsErr != nil {
				return nil, existsErr
			}
			if taken {
				return nil, service.ErrProjectSlugTaken
			}
		}
		project.Title = normalizedTitle
		project.Slug = projectSlug
	}

	replaceRelations := teamIDs != nil || environmentIDs != nil || tagIDs != nil
	if replaceRelations {
		nextTeams := teamIDs
		nextEnvs := environmentIDs
		nextTags := tagIDs
		if nextTeams == nil {
			ids := make([]string, 0, len(project.Teams))
			for _, t := range project.Teams {
				ids = append(ids, t.ID.String())
			}
			nextTeams = &ids
		}
		if nextEnvs == nil {
			ids := make([]string, 0, len(project.Environments))
			for _, e := range project.Environments {
				ids = append(ids, e.ID.String())
			}
			nextEnvs = &ids
		}
		if nextTags == nil {
			ids := make([]string, 0, len(project.Tags))
			for _, t := range project.Tags {
				ids = append(ids, t.ID.String())
			}
			nextTags = &ids
		}

		teams, envs, tags, relErr := s.resolveRelations(ctx, *nextTeams, *nextEnvs, *nextTags)
		if relErr != nil {
			return nil, relErr
		}
		if err := s.projects.ReplaceAssociations(ctx, project, teams, envs, tags); err != nil {
			return nil, err
		}
	}

	updated, err := s.projects.Update(ctx, project)
	if err != nil {
		return nil, err
	}

	pid := updated.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionProjectUpdated,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &pid,
		ProjectTitle: updated.Title,
		ProjectSlug:  updated.Slug,
	})

	resp := toResponse(*updated)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, actorID, actorEmail, actorRole, projectID string) error {
	if err := requireAdmin(actorRole); err != nil {
		return err
	}

	id, err := uuid.Parse(strings.TrimSpace(projectID))
	if err != nil {
		return service.ErrProjectNotFound
	}

	project, err := s.projects.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrProjectNotFound
		}
		return err
	}

	if err := s.projects.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrProjectNotFound
		}
		return err
	}

	projectUUID := project.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionProjectDeleted,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &projectUUID,
		ProjectTitle: project.Title,
		ProjectSlug:  project.Slug,
	})
	return nil
}
