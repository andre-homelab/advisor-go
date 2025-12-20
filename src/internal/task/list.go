package task

import (
	"fmt"
)

func List() {
	jsonPath := Env("CANONICAL_TASKS_PATH", "../data/tasks.json")
	store := NewJSONStore(jsonPath)

	tasks, err := store.Load()

	if err != nil {
		fmt.Println("erro no carregamento do arquivo")
		return
	}

	fmt.Print("\n-=== TAREFAS EXISTENTES ===-\n\n")

	for i := range tasks {
		fmt.Printf("ID: %s \n| > Título: %s\n| > Descrição: %s\n| > Prioridade: %s \n| > Lembrete: %s", tasks[i].ID, tasks[i].Title, tasks[i].Description, tasks[i].Priority, tasks[i].ReminderAt)
		fmt.Print("\n\n")
	}

	fmt.Println("-==========================-")
}
