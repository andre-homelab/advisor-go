package task

import (
	"fmt"
	"time"
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

	store := NewJSONStore("/home/andre/Documentos/Git/advisor-go/src/data/tasks.json")

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
