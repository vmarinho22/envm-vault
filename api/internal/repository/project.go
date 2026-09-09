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

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) WithDB(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

type ProjectListFilter struct {
	TitleQuery string
	TeamIDs    []uuid.UUID
	TagIDs     []uuid.UUID
	Page       int
	PageSize   int
}

func projectPreload(q *gorm.DB) *gorm.DB {
	return q.Preload("Teams").Preload("Environments").Preload("Tags")
}

func (r *ProjectRepository) List(ctx context.Context, filter ProjectListFilter) ([]model.Project, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.Project{})

	if term := strings.TrimSpace(filter.TitleQuery); term != "" {
		q = q.Where(`title ILIKE ? ESCAPE '\'`, "%"+escapeILIKE(term)+"%")
	}
	if len(filter.TeamIDs) > 0 {
		q = q.Where(
			`EXISTS (SELECT 1 FROM project_teams pt WHERE pt.project_id = projects.id AND pt.team_id IN ?)`,
			filter.TeamIDs,
		)
	}
	if len(filter.TagIDs) > 0 {
		q = q.Where(
			`EXISTS (SELECT 1 FROM project_tags ptag WHERE ptag.project_id = projects.id AND ptag.tag_id IN ?)`,
			filter.TagIDs,
		)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count projects: %w", err)
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

	var projects []model.Project
	if err := projectPreload(q.Order("title ASC").Offset(offset).Limit(pageSize)).Find(&projects).Error; err != nil {
		return nil, 0, fmt.Errorf("list projects: %w", err)
	}
	return projects, total, nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	var project model.Project
	err := projectPreload(r.db.WithContext(ctx)).Where("id = ?", id).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find project by id: %w", err)
	}
	return &project, nil
}

func (r *ProjectRepository) FindBySlug(ctx context.Context, slug string) (*model.Project, error) {
	var project model.Project
	err := projectPreload(r.db.WithContext(ctx)).Where("slug = ?", slug).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find project by slug: %w", err)
	}
	return &project, nil
}

func (r *ProjectRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Project{}).Where("slug = ?", slug)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check project slug exists: %w", err)
	}
	return count > 0, nil
}

func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) (*model.Project, error) {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
	return r.FindByID(ctx, project.ID)
}

func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) (*model.Project, error) {
	if err := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: false}).Save(project).Error; err != nil {
		return nil, fmt.Errorf("update project: %w", err)
	}
	return r.FindByID(ctx, project.ID)
}

func (r *ProjectRepository) ReplaceAssociations(
	ctx context.Context,
	project *model.Project,
	teams []model.Team,
	environments []model.Environment,
	tags []model.Tag,
) error {
	db := r.db.WithContext(ctx)
	if err := db.Model(project).Association("Teams").Replace(teams); err != nil {
		return fmt.Errorf("replace project teams: %w", err)
	}
	if err := db.Model(project).Association("Environments").Replace(environments); err != nil {
		return fmt.Errorf("replace project environments: %w", err)
	}
	if err := db.Model(project).Association("Tags").Replace(tags); err != nil {
		return fmt.Errorf("replace project tags: %w", err)
	}
	return nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	project, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	db := r.db.WithContext(ctx)
	if err := db.Model(project).Association("Teams").Clear(); err != nil {
		return fmt.Errorf("clear project teams: %w", err)
	}
	if err := db.Model(project).Association("Environments").Clear(); err != nil {
		return fmt.Errorf("clear project environments: %w", err)
	}
	if err := db.Model(project).Association("Tags").Clear(); err != nil {
		return fmt.Errorf("clear project tags: %w", err)
	}
	res := db.Where("id = ?", id).Delete(&model.Project{})
	if res.Error != nil {
		return fmt.Errorf("delete project: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
