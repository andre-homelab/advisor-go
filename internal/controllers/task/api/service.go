package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

var (
	ErrTaskNotFound       = errors.New("tarefa não encontrada")
	ErrInvalidInput       = errors.New("dados de entrada inválidos")
	ErrParentTaskNotFound = errors.New("tarefa pai não encontrada")
)

type Store interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	List(ctx context.Context) ([]models.Task, error)
	Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error)
	Delete(ctx context.Context, id string) error
}

type Service struct {
	repo Store
}

func NewService(repo Store) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ ERRO ] Problema ao listar as Tasks: %w", err)
	}
	return tasks, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*models.Task, error) {
	task, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (s *Service) Create(ctx context.Context, title, description string, priority models.Priority, reminderAt time.Time) (*models.Task, error) {
	return s.CreateWithParent(ctx, title, description, priority, reminderAt, nil)
}

func (s *Service) CreateWithParent(ctx context.Context, title, description string, priority models.Priority, reminderAt time.Time, parentID *string) (*models.Task, error) {
	if parentID != nil && *parentID == "" {
		return nil, ErrInvalidInput
	}
	if parentID != nil {
		parentTask, err := s.repo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("[ ERRO ] Problema ao validar tarefa pai: %w", err)
		}
		if parentTask == nil {
			return nil, ErrParentTaskNotFound
		}
	}

	newTask := models.Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		ReminderAt:  reminderAt,
		Done:        false,
		ParentID:    parentID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, &newTask); err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (s *Service) Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
	if value, ok := changes["parent_id"]; ok {
		parentID, ok := value.(string)
		if !ok {
			return nil, ErrInvalidInput
		}
		if parentID == "" {
			return nil, ErrInvalidInput
		}
		if parentID == id {
			return nil, ErrInvalidInput
		}
		parentTask, err := s.repo.GetByID(ctx, parentID)
		if err != nil {
			return nil, fmt.Errorf("[ ERRO ] Problema ao validar tarefa pai: %w", err)
		}
		if parentTask == nil {
			return nil, ErrParentTaskNotFound
		}
	}

	task, err := s.repo.Patch(ctx, id, changes)

	if err != nil {
		return nil, fmt.Errorf("[ ERRO ] Problema dentro do Patch: %w", err)
	}
	return task, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)

	if err != nil {
		return fmt.Errorf("[ ERRO ] Problema dentro do Delete: %w", err)
	}

	return nil
}

func (s *Service) Complete(ctx context.Context, id string) (*models.Task, error) {
	changes := map[string]any{}
	changes["done"] = true

	task, err := s.repo.Patch(ctx, id, changes)

	if err != nil {
		return nil, fmt.Errorf("[ ERRO ] Problema ao completar tarefa: %w", err)
	}

	return task, nil
}

func ParsePriority(s string) (models.Priority, error) {
	switch s {
	case "low", "baixa":
		return models.PriorityLow, nil
	case "medium", "media", "média":
		return models.PriorityMedium, nil
	case "high", "alta":
		return models.PriorityHigh, nil
	default:
		return "", ErrInvalidInput
	}
}
