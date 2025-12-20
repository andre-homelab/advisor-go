package notify

import (
	"fmt"

	"github.com/andre-felipe-wonsik-alves/advisor-go/internal/task"
)

func NotifyToTerminal(t task.Task) {
	fmt.Printf("\n[ TAREFA VENCIDA ]\nTítulo: %s \n| Prioridade: %s\n", t.Title, t.Priority)

	if t.Description != "" {
		fmt.Printf("| Descrição: %s\n", t.Description)
	}
}
