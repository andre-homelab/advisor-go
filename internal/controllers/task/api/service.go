package api

import (
	"errors"
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/controllers/task"
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

func (s *Service) List() ([]task.Task, error) {
	return s.store.Load()
}

func (s *Service) GetByID(id string) (*task.Task, error) {
	tasks, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	for i := range tasks {
		if tasks[i].ID == id {
			return &tasks[i], nil
		}
	}

	return nil, ErrTaskNotFound
}

func (s *Service) Create(title, description string, priority task.Priority, reminderAt time.Time) (*task.Task, error) {
	tasks, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	newTask := task.Task{
		ID:          generateID(),
		Title:       title,
		Description: description,
		Priority:    priority,
		ReminderAt:  reminderAt,
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask)
	if err := s.store.Save(tasks); err != nil {
		return nil, err
	}

	return &newTask, nil
}

func (s *Service) Update(id string, title, description *string, priority *task.Priority, reminderAt *time.Time, done *bool) (*task.Task, error) {
	tasks, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	found := false
	var updated *task.Task

	for i := range tasks {
		if tasks[i].ID == id {
			if title != nil {
				tasks[i].Title = *title
			}
			if description != nil {
				tasks[i].Description = *description
			}
			if priority != nil {
				tasks[i].Priority = *priority
			}
			if reminderAt != nil {
				tasks[i].ReminderAt = *reminderAt
			}
			if done != nil {
				tasks[i].Done = *done
			}

			tasks[i].UpdatedAt = time.Now()
			updated = &tasks[i]
			found = true
			break
		}
	}

	if !found {
		return nil, ErrTaskNotFound
	}

	if err := s.store.Save(tasks); err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(id string) error {
	tasks, err := s.store.Load()
	if err != nil {
		return err
	}

	found := false
	newTasks := make([]task.Task, 0, len(tasks))

	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		} else {
			found = true
		}
	}

	if !found {
		return ErrTaskNotFound
	}

	return s.store.Save(newTasks)
}

func (s *Service) Complete(id string) (*task.Task, error) {
	tasks, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	found := false
	var completed *task.Task

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			tasks[i].UpdatedAt = time.Now()
			completed = &tasks[i]
			found = true
			break
		}
	}

	if !found {
		return nil, ErrTaskNotFound
	}

	if err := s.store.Save(tasks); err != nil {
		return nil, err
	}

	return completed, nil
}

func (s *Service) GetDue() ([]task.Task, error) {
	tasks, err := s.store.Load()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	dueTasks := make([]task.Task, 0)

	for _, t := range tasks {
		if !t.Done && !t.ReminderAt.IsZero() && t.ReminderAt.Before(now) {
			dueTasks = append(dueTasks, t)
		}
	}

	return dueTasks, nil
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

func ParsePriority(s string) (task.Priority, error) {
	switch s {
	case "low", "baixa":
		return task.PriorityLow, nil
	case "medium", "media", "média":
		return task.PriorityMedium, nil
	case "high", "alta":
		return task.PriorityHigh, nil
	default:
		return "", ErrInvalidInput
	}
}

func PriorityToString(p task.Priority) string {
	switch p {
	case task.PriorityLow:
		return "baixa"
	case task.PriorityMedium:
		return "média"
	case task.PriorityHigh:
		return "alta"
	default:
		return string(p)
	}
}
