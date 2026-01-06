package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/controllers/task"
	"github.com/andre-felipe-wonsik-alves/internal/models"
)

var (
	ErrTaskNotFound = errors.New("tarefa não encontrada")
	ErrInvalidInput = errors.New("dados de entrada inválidos")
)

type Service struct {
	store task.Store
}

func NewService(store task.Store) *Service {
	return &Service{store: store}
}

func (s *Service) List(ctx context.Context) ([]models.Task, error) {
	tasks, err := s.store.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar tarefas: %w", err)
	}
	return tasks, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*models.Task, error) {
	task, err := s.store.GetByID(ctx, id)

	if err != nil {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (s *Service) Create(ctx context.Context, title, description string, priority models.Priority, reminderAt time.Time) (*models.Task, error) {
	newTask := models.Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		ReminderAt:  reminderAt,
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.store.Create(ctx, &newTask); err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (s *Service) Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
	return s.store.Patch(ctx, id, changes)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) Complete(ctx context.Context, id string) (*models.Task, error) {
	changes := map[string]any{}

	changes["done"] = true

	return s.store.Patch(ctx, id, changes)
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

func PriorityToString(p models.Priority) string {
	switch p {
	case models.PriorityLow:
		return "baixa"
	case models.PriorityMedium:
		return "média"
	case models.PriorityHigh:
		return "alta"
	default:
		return string(p)
	}
}
