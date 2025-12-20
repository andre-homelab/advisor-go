package task

import (
	"fmt"
	"time"
)

func Due() ([]Task, error) {
	jsonPath := Env("CANONICAL_TASKS_PATH", "../data/tasks.json")
	store := NewJSONStore(jsonPath)

	tasks, err := store.Load()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar tarefas: %w", err)
	}

	now := time.Now()
	dueTasks := make([]Task, 0)

	for _, t := range tasks {
		if isDue(t, now) {
			dueTasks = append(dueTasks, t)
		}
	}

	return dueTasks, nil
}

func isDue(task Task, now time.Time) bool {
	if task.Done || task.ReminderAt.IsZero() {
		return false
	}

	return !task.ReminderAt.After(now) // se já passou, true
}
