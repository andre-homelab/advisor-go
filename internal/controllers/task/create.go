package task

import (
	"fmt"
	"time"

	env "github.com/andre-felipe-wonsik-alves/internal"
)

func Create(title string, description string, priority Priority, reminderAt time.Time) (Task, error) {
	now := time.Now()

	newTask := Task{
		ID:          fmt.Sprintf("%d", now.UnixNano()),
		Title:       title,
		Description: description,
		Priority:    priority,
		ReminderAt:  reminderAt,
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	jsonPath := env.Env("CANONICAL_TASKS_PATH", "../data/tasks.json")
	store := NewJSONStore(jsonPath)

	tasks, err := store.Load()
	if err != nil {
		return newTask, fmt.Errorf("erro ao carregar tarefas: %w", err)
	}

	tasks = append(tasks, newTask)

	if err := store.Save(tasks); err != nil {
		return newTask, fmt.Errorf("erro ao salvar tarefas: %w", err)
	}

	return newTask, nil
}
