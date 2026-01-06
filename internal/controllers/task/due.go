package task

import (
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

func Due() ([]models.Task, error) {
	// jsonPath := env.GetEnv("CANONICAL_TASKS_PATH", "../data/tasks.json")
	// store := NewJSONStore(jsonPath)

	// tasks, err := store.Load()
	// if err != nil {
	// 	return nil, fmt.Errorf("erro ao carregar tarefas: %w", err)
	// }

	// now := time.Now()
	dueTasks := make([]models.Task, 0)

	// for _, t := range tasks {
	// 	if isDue(t, now) {
	// 		dueTasks = append(dueTasks, t)
	// 	}
	// }

	return dueTasks, nil
}

func isDue(task models.Task, now time.Time) bool {
	if task.Done || task.ReminderAt.IsZero() {
		return false
	}

	return !task.ReminderAt.After(now) // se já passou, true
}
