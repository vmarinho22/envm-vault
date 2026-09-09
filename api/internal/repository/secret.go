package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/vmarinho22/envm-vault/api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SecretRepository struct {
	db *gorm.DB
}

func NewSecretRepository(db *gorm.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

func (r *SecretRepository) WithDB(db *gorm.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

func (r *SecretRepository) ListByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.Secret, error) {
	var secrets []model.Secret
	err := r.db.WithContext(ctx).
		Preload("Values").
		Joins("JOIN project_secrets ps ON ps.secret_id = secrets.id").
		Where("ps.project_id = ?", projectID).
		Order("secrets.name ASC").
		Find(&secrets).Error
	if err != nil {
		return nil, fmt.Errorf("list secrets by project: %w", err)
	}
	return secrets, nil
}

func (r *SecretRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Secret, error) {
	var secret model.Secret
	err := r.db.WithContext(ctx).Preload("Values").Where("id = ?", id).First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find secret by id: %w", err)
	}
	return &secret, nil
}

func (r *SecretRepository) FindByProjectAndID(ctx context.Context, projectID, secretID uuid.UUID) (*model.Secret, error) {
	var secret model.Secret
	err := r.db.WithContext(ctx).
		Preload("Values").
		Joins("JOIN project_secrets ps ON ps.secret_id = secrets.id").
		Where("secrets.id = ? AND ps.project_id = ?", secretID, projectID).
		First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find secret by project and id: %w", err)
	}
	return &secret, nil
}

func (r *SecretRepository) ExistsByName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Secret{}).Where("name = ?", name)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check secret name exists: %w", err)
	}
	return count > 0, nil
}

func (r *SecretRepository) IsAttached(ctx context.Context, projectID, secretID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ProjectSecret{}).
		Where("project_id = ? AND secret_id = ?", projectID, secretID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check secret attached: %w", err)
	}
	return count > 0, nil
}

func (r *SecretRepository) Attach(ctx context.Context, projectID, secretID uuid.UUID) error {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).
		Create(&model.ProjectSecret{ProjectID: projectID, SecretID: secretID}).Error
	if err != nil {
		return fmt.Errorf("attach secret to project: %w", err)
	}
	return nil
}

func (r *SecretRepository) Detach(ctx context.Context, projectID, secretID uuid.UUID) error {
	if err := r.DeletePrefixesByProjectAndSecret(ctx, projectID, secretID); err != nil {
		return err
	}
	res := r.db.WithContext(ctx).
		Where("project_id = ? AND secret_id = ?", projectID, secretID).
		Delete(&model.ProjectSecret{})
	if res.Error != nil {
		return fmt.Errorf("detach secret from project: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SecretRepository) DetachAllByProjectID(ctx context.Context, projectID uuid.UUID) error {
	if err := r.DeletePrefixesByProjectID(ctx, projectID); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&model.ProjectSecret{}).Error; err != nil {
		return fmt.Errorf("detach all secrets for project: %w", err)
	}
	return nil
}

func (r *SecretRepository) CountProjects(ctx context.Context, secretID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ProjectSecret{}).
		Where("secret_id = ?", secretID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count secret projects: %w", err)
	}
	return count, nil
}

type SecretCatalogFilter struct {
	TitleQuery        string
	ExcludeProjectID  *uuid.UUID
	Page              int
	PageSize          int
}

// ListCatalog returns a page of catalog secrets with total count.
func (r *SecretRepository) ListCatalog(ctx context.Context, filter SecretCatalogFilter) ([]model.Secret, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Secret{})
	if term := strings.TrimSpace(filter.TitleQuery); term != "" {
		q = q.Where(`name ILIKE ? ESCAPE '\'`, "%"+escapeILIKE(term)+"%")
	}
	if filter.ExcludeProjectID != nil {
		q = q.Where(
			`NOT EXISTS (
				SELECT 1 FROM project_secrets ps
				WHERE ps.secret_id = secrets.id AND ps.project_id = ?
			)`,
			*filter.ExcludeProjectID,
		)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count catalog secrets: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 12
	}
	offset := (page - 1) * pageSize

	var secrets []model.Secret
	if err := q.Order("name ASC").Offset(offset).Limit(pageSize).Find(&secrets).Error; err != nil {
		return nil, 0, fmt.Errorf("list catalog secrets: %w", err)
	}
	return secrets, total, nil
}

func (r *SecretRepository) ListProjectsBySecretID(ctx context.Context, secretID uuid.UUID) ([]model.Project, error) {
	var projects []model.Project
	err := r.db.WithContext(ctx).
		Joins("JOIN project_secrets ps ON ps.project_id = projects.id").
		Where("ps.secret_id = ?", secretID).
		Order("projects.title ASC").
		Find(&projects).Error
	if err != nil {
		return nil, fmt.Errorf("list projects by secret: %w", err)
	}
	return projects, nil
}

func (r *SecretRepository) Create(ctx context.Context, secret *model.Secret) (*model.Secret, error) {
	if err := r.db.WithContext(ctx).Create(secret).Error; err != nil {
		return nil, fmt.Errorf("create secret: %w", err)
	}
	return r.FindByID(ctx, secret.ID)
}

func (r *SecretRepository) UpdateName(ctx context.Context, id uuid.UUID, name string) error {
	res := r.db.WithContext(ctx).Model(&model.Secret{}).Where("id = ?", id).Update("name", name)
	if res.Error != nil {
		return fmt.Errorf("update secret name: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SecretRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Secret{})
	if res.Error != nil {
		return fmt.Errorf("delete secret: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SecretRepository) UpsertValue(ctx context.Context, value *model.SecretValue) error {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "secret_id"}, {Name: "environment_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"ciphertext", "key_version", "updated_at"}),
	}).Create(value).Error
	if err != nil {
		return fmt.Errorf("upsert secret value: %w", err)
	}
	return nil
}

func (r *SecretRepository) FindValue(ctx context.Context, secretID, environmentID uuid.UUID) (*model.SecretValue, error) {
	var value model.SecretValue
	err := r.db.WithContext(ctx).
		Where("secret_id = ? AND environment_id = ?", secretID, environmentID).
		First(&value).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find secret value: %w", err)
	}
	return &value, nil
}

func (r *SecretRepository) DeleteValue(ctx context.Context, secretID, environmentID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("secret_id = ? AND environment_id = ?", secretID, environmentID).
		Delete(&model.SecretValue{})
	if res.Error != nil {
		return fmt.Errorf("delete secret value: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ListValuesForEnvironment returns secret values for secrets attached to the project
// that have a value for the given environment.
func (r *SecretRepository) ListValuesForEnvironment(ctx context.Context, projectID, environmentID uuid.UUID) ([]model.SecretValue, error) {
	var values []model.SecretValue
	err := r.db.WithContext(ctx).
		Joins("JOIN secrets ON secrets.id = secret_values.secret_id").
		Joins("JOIN project_secrets ps ON ps.secret_id = secrets.id").
		Preload("Secret").
		Where("ps.project_id = ? AND secret_values.environment_id = ?", projectID, environmentID).
		Order("secrets.name ASC").
		Find(&values).Error
	if err != nil {
		return nil, fmt.Errorf("list secret values for environment: %w", err)
	}
	return values, nil
}

func (r *SecretRepository) UpsertPrefix(ctx context.Context, row *model.ProjectSecretPrefix) error {
	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "project_id"},
			{Name: "secret_id"},
			{Name: "environment_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"prefix", "updated_at"}),
	}).Create(row).Error
	if err != nil {
		return fmt.Errorf("upsert project secret prefix: %w", err)
	}
	return nil
}

func (r *SecretRepository) DeletePrefix(ctx context.Context, projectID, secretID, environmentID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("project_id = ? AND secret_id = ? AND environment_id = ?", projectID, secretID, environmentID).
		Delete(&model.ProjectSecretPrefix{})
	if res.Error != nil {
		return fmt.Errorf("delete project secret prefix: %w", res.Error)
	}
	return nil
}

func (r *SecretRepository) DeletePrefixesByProjectAndSecret(ctx context.Context, projectID, secretID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("project_id = ? AND secret_id = ?", projectID, secretID).
		Delete(&model.ProjectSecretPrefix{}).Error; err != nil {
		return fmt.Errorf("delete prefixes by project and secret: %w", err)
	}
	return nil
}

func (r *SecretRepository) DeletePrefixesByProjectID(ctx context.Context, projectID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Delete(&model.ProjectSecretPrefix{}).Error; err != nil {
		return fmt.Errorf("delete prefixes by project: %w", err)
	}
	return nil
}

// DeletePrefixesOutsideEnvironments removes prefixes for envs no longer linked to the project.
func (r *SecretRepository) DeletePrefixesOutsideEnvironments(ctx context.Context, projectID uuid.UUID, keepEnvIDs []uuid.UUID) error {
	q := r.db.WithContext(ctx).Where("project_id = ?", projectID)
	if len(keepEnvIDs) == 0 {
		if err := q.Delete(&model.ProjectSecretPrefix{}).Error; err != nil {
			return fmt.Errorf("delete all project prefixes: %w", err)
		}
		return nil
	}
	if err := q.Where("environment_id NOT IN ?", keepEnvIDs).Delete(&model.ProjectSecretPrefix{}).Error; err != nil {
		return fmt.Errorf("delete prefixes outside environments: %w", err)
	}
	return nil
}

func (r *SecretRepository) ListPrefixesByProjectID(ctx context.Context, projectID uuid.UUID) ([]model.ProjectSecretPrefix, error) {
	var rows []model.ProjectSecretPrefix
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list prefixes by project: %w", err)
	}
	return rows, nil
}

func (r *SecretRepository) ListPrefixesForEnvironment(
	ctx context.Context,
	projectID, environmentID uuid.UUID,
) ([]model.ProjectSecretPrefix, error) {
	var rows []model.ProjectSecretPrefix
	err := r.db.WithContext(ctx).
		Where("project_id = ? AND environment_id = ?", projectID, environmentID).
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list prefixes for environment: %w", err)
	}
	return rows, nil
}
