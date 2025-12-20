package cmd

import (
	"fmt"

	"github.com/andre-felipe-wonsik-alves/advisor-go/internal/notify"
	"github.com/andre-felipe-wonsik-alves/advisor-go/internal/task"
	"github.com/spf13/cobra"
)

var dueCmd = &cobra.Command{
	Use:   "due",
	Short: "Verifica tarefas com vencimentos próximos/vencidos.",
	RunE: func(cmd *cobra.Command, args []string) error {
		tasks, err := task.Due()

		for i := range tasks {
			notify.NotifyToTerminal(tasks[i])
		}

		if err != nil {
			return fmt.Errorf("Erro na busca de tarefas pendendtes.")
		}
		return nil
	},
}
