package cli

import (
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task"
	"github.com/spf13/cobra"
)

var listCli = &cobra.Command{
	Use:   "list",
	Short: "Lista todas as tarefas.",
	RunE: func(cli *cobra.Command, args []string) error {
		task.List()

		return nil
	},
}
