package repository

import (
	"context"
	"errors"

	"github.com/andre-felipe-wonsik-alves/internal/models"
	"gorm.io/gorm"
)

type DBStore struct {
	db *gorm.DB
}

func NewDBStore(db *gorm.DB) *DBStore {
	return &DBStore{db: db}
}

func (s *DBStore) Create(ctx context.Context, t *models.Task) error {
	return s.db.WithContext(ctx).Create(t).Error
}

func (s *DBStore) GetByID(ctx context.Context, id string) (*models.Task, error) {
	var t models.Task
	err := s.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		First(&t, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *DBStore) List(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := s.db.WithContext(ctx).
		Preload("Parent").
		Preload("Children").
		Order("created_at desc").
		Find(&tasks).Error
	return tasks, err
}

func (s *DBStore) Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
	tx := s.db.WithContext(ctx).Model(&models.Task{}).Where("id = ?", id).Updates(changes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, nil
	}
	return s.GetByID(ctx, id)
}

func (s *DBStore) Delete(ctx context.Context, id string) error {
	tx := s.db.WithContext(ctx).Delete(&models.Task{}, "id = ?", id)
	return tx.Error
}
