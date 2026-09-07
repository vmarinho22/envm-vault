package secret

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/envm-vault/envm-vault/api/internal/repository"
	"github.com/envm-vault/envm-vault/api/internal/secretcrypto"
	"github.com/envm-vault/envm-vault/api/internal/service"
	"github.com/envm-vault/envm-vault/api/internal/service/audit"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	maxSecretNameLen           = 128
	productionEnvironmentTitle = "production"
)

func assertProductionDownloadAllowed(actorRole, envTitle string) error {
	if envTitle == productionEnvironmentTitle && actorRole != service.RoleSuperAdmin {
		return service.ErrProductionDownloadForbidden
	}
	return nil
}

var secretNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type Service struct {
	secrets  *repository.SecretRepository
	projects *repository.ProjectRepository
	crypto   secretcrypto.Crypter
	audit    audit.Logger
}

func New(
	secrets *repository.SecretRepository,
	projects *repository.ProjectRepository,
	crypto secretcrypto.Crypter,
	auditLog audit.Logger,
) *Service {
	return &Service{
		secrets:  secrets,
		projects: projects,
		crypto:   crypto,
		audit:    audit.OrNop(auditLog),
	}
}

type ValueMeta struct {
	HasValue  bool      `json:"has_value"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type SecretResponse struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Values    map[string]ValueMeta `json:"values"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type ListResponse struct {
	Items []SecretResponse `json:"items"`
}

type RevealResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DownloadItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DownloadResponse struct {
	Items []DownloadItem `json:"items"`
}

func normalizeName(raw string) (string, error) {
	name := strings.TrimSpace(raw)
	if name == "" || utf8.RuneCountInString(name) > maxSecretNameLen {
		return "", service.ErrInvalidSecretName
	}
	if !secretNamePattern.MatchString(name) {
		return "", service.ErrInvalidSecretName
	}
	return name, nil
}

func parseActorID(raw string) *uuid.UUID {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil
	}
	return &id
}

func (s *Service) loadProject(ctx context.Context, projectID string) (*model.Project, error) {
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
	return project, nil
}

func projectEnvSet(project *model.Project) map[uuid.UUID]struct{} {
	set := make(map[uuid.UUID]struct{}, len(project.Environments))
	for _, e := range project.Environments {
		set[e.ID] = struct{}{}
	}
	return set
}

func toResponse(secret model.Secret) SecretResponse {
	values := make(map[string]ValueMeta, len(secret.Values))
	for _, v := range secret.Values {
		values[v.EnvironmentID.String()] = ValueMeta{
			HasValue:  true,
			UpdatedAt: v.UpdatedAt,
		}
	}
	return SecretResponse{
		ID:        secret.ID.String(),
		Name:      secret.Name,
		Values:    values,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}
}

func (s *Service) List(ctx context.Context, projectID string) (*ListResponse, error) {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	secrets, err := s.secrets.ListByProjectID(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	items := make([]SecretResponse, 0, len(secrets))
	for _, sec := range secrets {
		items = append(items, toResponse(sec))
	}
	return &ListResponse{Items: items}, nil
}

func (s *Service) Create(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectID, name string,
	values map[string]string,
) (*SecretResponse, error) {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	normalized, err := normalizeName(name)
	if err != nil {
		return nil, err
	}
	taken, err := s.secrets.ExistsByName(ctx, project.ID, normalized, nil)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, service.ErrSecretNameTaken
	}

	secret := &model.Secret{
		ProjectID: project.ID,
		Name:      normalized,
	}
	created, err := s.secrets.Create(ctx, secret)
	if err != nil {
		return nil, err
	}

	envSet := projectEnvSet(project)
	for envIDRaw, plaintext := range values {
		envID, parseErr := uuid.Parse(strings.TrimSpace(envIDRaw))
		if parseErr != nil {
			_ = s.secrets.Delete(ctx, created.ID)
			return nil, service.ErrEnvironmentNotFound
		}
		if _, ok := envSet[envID]; !ok {
			_ = s.secrets.Delete(ctx, created.ID)
			return nil, service.ErrEnvironmentNotInProject
		}
		if setErr := s.encryptAndUpsert(ctx, created.ID, envID, plaintext); setErr != nil {
			_ = s.secrets.Delete(ctx, created.ID)
			return nil, setErr
		}
	}

	fresh, err := s.secrets.FindByID(ctx, created.ID)
	if err != nil {
		return nil, err
	}

	pid := project.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionSecretCreated,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &pid,
		ProjectTitle: project.Title,
		ProjectSlug:  project.Slug,
	})

	resp := toResponse(*fresh)
	return &resp, nil
}

func (s *Service) UpdateName(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectID, secretID, name string,
) (*SecretResponse, error) {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	sid, err := uuid.Parse(strings.TrimSpace(secretID))
	if err != nil {
		return nil, service.ErrSecretNotFound
	}
	secret, err := s.secrets.FindByProjectAndID(ctx, project.ID, sid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrSecretNotFound
		}
		return nil, err
	}
	normalized, err := normalizeName(name)
	if err != nil {
		return nil, err
	}
	if normalized != secret.Name {
		taken, existsErr := s.secrets.ExistsByName(ctx, project.ID, normalized, &sid)
		if existsErr != nil {
			return nil, existsErr
		}
		if taken {
			return nil, service.ErrSecretNameTaken
		}
		if err := s.secrets.UpdateName(ctx, sid, normalized); err != nil {
			return nil, err
		}
	}

	fresh, err := s.secrets.FindByID(ctx, sid)
	if err != nil {
		return nil, err
	}

	pid := project.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionSecretUpdated,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &pid,
		ProjectTitle: project.Title,
		ProjectSlug:  project.Slug,
	})

	resp := toResponse(*fresh)
	return &resp, nil
}

func (s *Service) Delete(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectID, secretID string,
) error {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(strings.TrimSpace(secretID))
	if err != nil {
		return service.ErrSecretNotFound
	}
	if _, err := s.secrets.FindByProjectAndID(ctx, project.ID, sid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrSecretNotFound
		}
		return err
	}
	if err := s.secrets.Delete(ctx, sid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrSecretNotFound
		}
		return err
	}

	pid := project.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionSecretDeleted,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &pid,
		ProjectTitle: project.Title,
		ProjectSlug:  project.Slug,
	})
	return nil
}

func (s *Service) encryptAndUpsert(ctx context.Context, secretID, environmentID uuid.UUID, plaintext string) error {
	aad := secretcrypto.AAD(secretID, environmentID)
	blob, version, err := s.crypto.Encrypt([]byte(plaintext), aad)
	if err != nil {
		return fmt.Errorf("encrypt secret value: %w", err)
	}
	return s.secrets.UpsertValue(ctx, &model.SecretValue{
		SecretID:      secretID,
		EnvironmentID: environmentID,
		Ciphertext:    blob,
		KeyVersion:    version,
	})
}

func (s *Service) SetValue(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectID, secretID, environmentID, plaintext string,
) (*SecretResponse, error) {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	sid, err := uuid.Parse(strings.TrimSpace(secretID))
	if err != nil {
		return nil, service.ErrSecretNotFound
	}
	eid, err := uuid.Parse(strings.TrimSpace(environmentID))
	if err != nil {
		return nil, service.ErrEnvironmentNotFound
	}
	if _, ok := projectEnvSet(project)[eid]; !ok {
		return nil, service.ErrEnvironmentNotInProject
	}
	if _, err := s.secrets.FindByProjectAndID(ctx, project.ID, sid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrSecretNotFound
		}
		return nil, err
	}
	if err := s.encryptAndUpsert(ctx, sid, eid, plaintext); err != nil {
		return nil, err
	}

	fresh, err := s.secrets.FindByID(ctx, sid)
	if err != nil {
		return nil, err
	}

	pid := project.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionSecretValueSet,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &pid,
		ProjectTitle: project.Title,
		ProjectSlug:  project.Slug,
	})

	resp := toResponse(*fresh)
	return &resp, nil
}

func (s *Service) RevealValue(
	ctx context.Context,
	projectID, secretID, environmentID string,
) (*RevealResponse, error) {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	sid, err := uuid.Parse(strings.TrimSpace(secretID))
	if err != nil {
		return nil, service.ErrSecretNotFound
	}
	eid, err := uuid.Parse(strings.TrimSpace(environmentID))
	if err != nil {
		return nil, service.ErrEnvironmentNotFound
	}
	if _, ok := projectEnvSet(project)[eid]; !ok {
		return nil, service.ErrEnvironmentNotInProject
	}
	secret, err := s.secrets.FindByProjectAndID(ctx, project.ID, sid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrSecretNotFound
		}
		return nil, err
	}
	value, err := s.secrets.FindValue(ctx, sid, eid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrSecretValueNotFound
		}
		return nil, err
	}
	plaintext, err := s.crypto.Decrypt(value.Ciphertext, secretcrypto.AAD(sid, eid), value.KeyVersion)
	if err != nil {
		return nil, fmt.Errorf("decrypt secret value: %w", err)
	}
	return &RevealResponse{Name: secret.Name, Value: string(plaintext)}, nil
}

func (s *Service) ClearValue(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectID, secretID, environmentID string,
) error {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return err
	}
	sid, err := uuid.Parse(strings.TrimSpace(secretID))
	if err != nil {
		return service.ErrSecretNotFound
	}
	eid, err := uuid.Parse(strings.TrimSpace(environmentID))
	if err != nil {
		return service.ErrEnvironmentNotFound
	}
	if _, err := s.secrets.FindByProjectAndID(ctx, project.ID, sid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrSecretNotFound
		}
		return err
	}
	if err := s.secrets.DeleteValue(ctx, sid, eid); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service.ErrSecretValueNotFound
		}
		return err
	}

	pid := project.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionSecretValueCleared,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &pid,
		ProjectTitle: project.Title,
		ProjectSlug:  project.Slug,
	})
	return nil
}

func (s *Service) ExportEnvironment(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectID, environmentID string,
	auditDownload bool,
) (*DownloadResponse, error) {
	project, err := s.loadProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	eid, err := uuid.Parse(strings.TrimSpace(environmentID))
	if err != nil {
		return nil, service.ErrEnvironmentNotFound
	}
	var envTitle string
	for _, e := range project.Environments {
		if e.ID == eid {
			envTitle = e.Title
			break
		}
	}
	if envTitle == "" {
		return nil, service.ErrEnvironmentNotInProject
	}
	if err := assertProductionDownloadAllowed(actorRole, envTitle); err != nil {
		return nil, err
	}

	items, err := s.decryptEnvironment(ctx, project.ID, eid)
	if err != nil {
		return nil, err
	}

	if auditDownload {
		pid := project.ID
		s.audit.Record(ctx, audit.Event{
			Action:       audit.ActionSecretDownloaded,
			ActorUserID:  parseActorID(actorID),
			ActorEmail:   actorEmail,
			ActorRole:    actorRole,
			ProjectID:    &pid,
			ProjectTitle: project.Title,
			ProjectSlug:  project.Slug,
		})
	}

	return &DownloadResponse{Items: items}, nil
}

func (s *Service) DownloadBySlug(
	ctx context.Context,
	actorID, actorEmail, actorRole, projectSlug, envTitle string,
) (*DownloadResponse, error) {
	title := strings.TrimSpace(envTitle)
	if title == "" {
		return nil, service.ErrEnvironmentParamRequired
	}
	if err := assertProductionDownloadAllowed(actorRole, title); err != nil {
		return nil, err
	}

	project, err := s.projects.FindBySlug(ctx, strings.TrimSpace(projectSlug))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, service.ErrProjectNotFound
		}
		return nil, err
	}

	var envID *uuid.UUID
	for _, e := range project.Environments {
		if e.Title == title {
			id := e.ID
			envID = &id
			break
		}
	}
	if envID == nil {
		return nil, service.ErrEnvironmentNotInProject
	}

	items, err := s.decryptEnvironment(ctx, project.ID, *envID)
	if err != nil {
		return nil, err
	}

	pid := project.ID
	s.audit.Record(ctx, audit.Event{
		Action:       audit.ActionSecretDownloaded,
		ActorUserID:  parseActorID(actorID),
		ActorEmail:   actorEmail,
		ActorRole:    actorRole,
		ProjectID:    &pid,
		ProjectTitle: project.Title,
		ProjectSlug:  project.Slug,
	})

	return &DownloadResponse{Items: items}, nil
}

func (s *Service) decryptEnvironment(ctx context.Context, projectID, environmentID uuid.UUID) ([]DownloadItem, error) {
	values, err := s.secrets.ListValuesForEnvironment(ctx, projectID, environmentID)
	if err != nil {
		return nil, err
	}
	items := make([]DownloadItem, 0, len(values))
	for _, v := range values {
		plaintext, decErr := s.crypto.Decrypt(v.Ciphertext, secretcrypto.AAD(v.SecretID, environmentID), v.KeyVersion)
		if decErr != nil {
			return nil, fmt.Errorf("decrypt secret value: %w", decErr)
		}
		name := v.Secret.Name
		items = append(items, DownloadItem{Name: name, Value: string(plaintext)})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items, nil
}

// FormatDotEnv renders download items as a dotenv file body.
func FormatDotEnv(items []DownloadItem) string {
	var b strings.Builder
	for _, item := range items {
		b.WriteString(item.Name)
		b.WriteByte('=')
		b.WriteString(quoteDotEnv(item.Value))
		b.WriteByte('\n')
	}
	return b.String()
}

func quoteDotEnv(value string) string {
	if value == "" {
		return `""`
	}
	needsQuote := strings.ContainsAny(value, " \t\n\r\"'#")
	if !needsQuote {
		return value
	}
	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\"`)
	escaped = strings.ReplaceAll(escaped, "\n", `\n`)
	escaped = strings.ReplaceAll(escaped, "\r", `\r`)
	return `"` + escaped + `"`
}
