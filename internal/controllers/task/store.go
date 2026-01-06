package task

import (
	"context"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

type Store interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	List(ctx context.Context) ([]models.Task, error)
	Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error)
	Delete(ctx context.Context, id string) error
}

type DBTaskStore struct {
	repo DBRepository
}

type DBRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	List(ctx context.Context) ([]models.Task, error)
	Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error)
	Delete(ctx context.Context, id string) error
}

func NewDBTaskStore(repo DBRepository) *DBTaskStore {
	return &DBTaskStore{
		repo: repo,
	}
}

func (s *DBTaskStore) Create(ctx context.Context, task *models.Task) error {
	return s.repo.Create(ctx, task)
}

func (s *DBTaskStore) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DBTaskStore) List(ctx context.Context) ([]models.Task, error) {
	return s.repo.List(ctx)
}

func (s *DBTaskStore) Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
	return s.repo.Patch(ctx, id, changes)
}

func (s *DBTaskStore) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
