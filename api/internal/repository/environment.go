package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/vmarinho22/envm-vault/api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EnvironmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (r *EnvironmentRepository) WithDB(db *gorm.DB) *EnvironmentRepository {
	return &EnvironmentRepository{db: db}
}

func (r *EnvironmentRepository) List(ctx context.Context, titleQuery string) ([]model.Environment, error) {
	q := r.db.WithContext(ctx).Order("title ASC")
	if term := strings.TrimSpace(titleQuery); term != "" {
		q = q.Where(`title ILIKE ? ESCAPE '\'`, "%"+escapeILIKE(term)+"%")
	}

	var environments []model.Environment
	if err := q.Find(&environments).Error; err != nil {
		return nil, fmt.Errorf("list environments: %w", err)
	}
	return environments, nil
}

func (r *EnvironmentRepository) Create(ctx context.Context, environment *model.Environment) (*model.Environment, error) {
	if err := r.db.WithContext(ctx).Create(environment).Error; err != nil {
		return nil, fmt.Errorf("create environment: %w", err)
	}
	return environment, nil
}

func (r *EnvironmentRepository) Update(ctx context.Context, environment *model.Environment) (*model.Environment, error) {
	if err := r.db.WithContext(ctx).Save(environment).Error; err != nil {
		return nil, fmt.Errorf("update environment: %w", err)
	}
	return environment, nil
}

func (r *EnvironmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Environment, error) {
	var environment model.Environment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&environment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find environment by id: %w", err)
	}
	return &environment, nil
}

func (r *EnvironmentRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Environment, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var environments []model.Environment
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&environments).Error; err != nil {
		return nil, fmt.Errorf("find environments by ids: %w", err)
	}
	return environments, nil
}

func (r *EnvironmentRepository) ExistsByTitle(ctx context.Context, title string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Environment{}).Where("LOWER(title) = LOWER(?)", title)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check environment title exists: %w", err)
	}
	return count > 0, nil
}

func (r *EnvironmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Environment{})
	if res.Error != nil {
		return fmt.Errorf("delete environment: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
