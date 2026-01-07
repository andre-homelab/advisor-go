package task

import (
	"fmt"
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

func Create(title string, description string, priority models.Priority, reminderAt time.Time) (models.Task, error) {
	now := time.Now()

	newTask := models.Task{
		Title:       title,
		Description: description,
		Priority:    priority,
		ReminderAt:  reminderAt,
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// jsonPath := env.GetEnv("CANONICAL_TASKS_PATH", "../data/tasks.json")
	// store := NewJSONStore(jsonPath)

	// tasks, err := store.Load()
	// if err != nil {
	// 	return newTask, fmt.Errorf("erro ao carregar tarefas: %w", err)
	// }

	// tasks = append(tasks, newTask)

	// if err := store.Save(tasks); err != nil {
	// 	return newTask, fmt.Errorf("erro ao salvar tarefas: %w", err)
	// }

	// return newTask, nil

	fmt.Println("Create")
	return newTask, nil
}
