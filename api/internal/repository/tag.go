package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/envm-vault/envm-vault/api/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) WithDB(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) List(ctx context.Context, titleQuery string) ([]model.Tag, error) {
	q := r.db.WithContext(ctx).Order("title ASC")
	if term := strings.TrimSpace(titleQuery); term != "" {
		q = q.Where(`title ILIKE ? ESCAPE '\'`, "%"+escapeILIKE(term)+"%")
	}

	var tags []model.Tag
	if err := q.Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	return tags, nil
}

func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
		return nil, fmt.Errorf("create tag: %w", err)
	}
	return tag, nil
}

func (r *TagRepository) Update(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if err := r.db.WithContext(ctx).Save(tag).Error; err != nil {
		return nil, fmt.Errorf("update tag: %w", err)
	}
	return tag, nil
}

func (r *TagRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find tag by id: %w", err)
	}
	return &tag, nil
}

func (r *TagRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Tag, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var tags []model.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("find tags by ids: %w", err)
	}
	return tags, nil
}

func (r *TagRepository) ExistsByTitle(ctx context.Context, title string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Tag{}).Where("LOWER(title) = LOWER(?)", title)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check tag title exists: %w", err)
	}
	return count > 0, nil
}

func (r *TagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Tag{})
	if res.Error != nil {
		return fmt.Errorf("delete tag: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
