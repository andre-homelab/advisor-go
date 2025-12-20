package cmd

import (
	"github.com/andre-felipe-wonsik-alves/advisor-go/internal/task"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todas as tarefas.",
	RunE: func(cmd *cobra.Command, args []string) error {
		task.List()

		return nil
	},
}
