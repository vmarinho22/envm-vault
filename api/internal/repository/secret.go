package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/envm-vault/envm-vault/api/internal/model"
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
		Where("project_id = ?", projectID).
		Order("name ASC").
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
		Where("id = ? AND project_id = ?", secretID, projectID).
		First(&secret).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find secret by project and id: %w", err)
	}
	return &secret, nil
}

func (r *SecretRepository) ExistsByName(ctx context.Context, projectID uuid.UUID, name string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Secret{}).Where("project_id = ? AND name = ?", projectID, name)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check secret name exists: %w", err)
	}
	return count > 0, nil
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

// ListValuesForEnvironment returns secrets that have a value for the given project+environment.
func (r *SecretRepository) ListValuesForEnvironment(ctx context.Context, projectID, environmentID uuid.UUID) ([]model.SecretValue, error) {
	var values []model.SecretValue
	err := r.db.WithContext(ctx).
		Joins("JOIN secrets ON secrets.id = secret_values.secret_id").
		Preload("Secret").
		Where("secrets.project_id = ? AND secret_values.environment_id = ?", projectID, environmentID).
		Order("secrets.name ASC").
		Find(&values).Error
	if err != nil {
		return nil, fmt.Errorf("list secret values for environment: %w", err)
	}
	return values, nil
}

// DeleteValuesForEnvironments removes secret values for the given environments within a project.
func (r *SecretRepository) DeleteValuesForEnvironments(ctx context.Context, projectID uuid.UUID, environmentIDs []uuid.UUID) error {
	if len(environmentIDs) == 0 {
		return nil
	}
	err := r.db.WithContext(ctx).
		Where(
			`secret_id IN (SELECT id FROM secrets WHERE project_id = ?) AND environment_id IN ?`,
			projectID,
			environmentIDs,
		).
		Delete(&model.SecretValue{}).Error
	if err != nil {
		return fmt.Errorf("delete secret values for environments: %w", err)
	}
	return nil
}

// DeleteByProjectID removes all secrets (and cascading values) for a project.
func (r *SecretRepository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&model.Secret{}).Error; err != nil {
		return fmt.Errorf("delete secrets by project: %w", err)
	}
	return nil
}
