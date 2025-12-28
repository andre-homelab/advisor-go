package cli

import (
	"fmt"

	"github.com/andre-felipe-wonsik-alves/internal/controllers/notify"
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task"
	"github.com/spf13/cobra"
)

var dueCli = &cobra.Command{
	Use:   "due",
	Short: "Verifica tarefas com vencimentos próximos/vencidos.",
	RunE: func(cli *cobra.Command, args []string) error {
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
