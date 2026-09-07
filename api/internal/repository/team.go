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

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) WithDB(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func escapeILIKE(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

func (r *TeamRepository) List(ctx context.Context, nameQuery string) ([]model.Team, error) {
	q := r.db.WithContext(ctx).Preload("Accountable").Order("name ASC")
	if term := strings.TrimSpace(nameQuery); term != "" {
		q = q.Where(`name ILIKE ? ESCAPE '\'`, "%"+escapeILIKE(term)+"%")
	}

	var teams []model.Team
	if err := q.Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}
	return teams, nil
}

func (r *TeamRepository) Create(ctx context.Context, team *model.Team) (*model.Team, error) {
	if err := r.db.WithContext(ctx).Create(team).Error; err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}
	if err := r.db.WithContext(ctx).Preload("Accountable").First(team, "id = ?", team.ID).Error; err != nil {
		return nil, fmt.Errorf("reload team after create: %w", err)
	}
	return team, nil
}

func (r *TeamRepository) Update(ctx context.Context, team *model.Team) (*model.Team, error) {
	if err := r.db.WithContext(ctx).Save(team).Error; err != nil {
		return nil, fmt.Errorf("update team: %w", err)
	}
	if err := r.db.WithContext(ctx).Preload("Accountable").First(team, "id = ?", team.ID).Error; err != nil {
		return nil, fmt.Errorf("reload team after update: %w", err)
	}
	return team, nil
}

func (r *TeamRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Team, error) {
	var team model.Team
	err := r.db.WithContext(ctx).Preload("Accountable").Where("id = ?", id).First(&team).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("find team by id: %w", err)
	}
	return &team, nil
}

func (r *TeamRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]model.Team, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var teams []model.Team
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&teams).Error; err != nil {
		return nil, fmt.Errorf("find teams by ids: %w", err)
	}
	return teams, nil
}

func (r *TeamRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Team{}).Where("slug = ?", slug)
	if excludeID != nil {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, fmt.Errorf("check team slug exists: %w", err)
	}
	return count > 0, nil
}

func (r *TeamRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Team{})
	if res.Error != nil {
		return fmt.Errorf("delete team: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
