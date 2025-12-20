package cmd

import (
	"fmt"
	"time"

	"github.com/andre-felipe-wonsik-alves/advisor-go/internal/notify"
	"github.com/andre-felipe-wonsik-alves/advisor-go/internal/task"
	"github.com/spf13/cobra"
)

var dueCmd = &cobra.Command{
	Use:   "due",
	Short: "Verifica tarefas com vencimentos próximos/vencidos.",
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonPath := task.Env("CANONICAL_TASKS_PATH", "../data/tasks.json")
		store := task.NewJSONStore(jsonPath)

		tasks, err := store.Load()

		if err != nil {
			return fmt.Errorf("Erro ao carregar tarefas: %w", err)
		}

		now := time.Now()
		found := 0

		for _, t := range tasks {
			if task.IsDue(t, now) {
				notify.NotifyToTerminal(t)
				found++
			}
		}

		if found == 0 {
			fmt.Println("Nenhuma tarefa vencida no momento.")
		}

		return nil
	},
}
