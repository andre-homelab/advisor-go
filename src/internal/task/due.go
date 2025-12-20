package task

import "time"

func IsDue(task Task, now time.Time) bool {
	if task.Done || task.ReminderAt.IsZero() {
		return false
	}

	return !task.ReminderAt.After(now) // se já passou, true
}
