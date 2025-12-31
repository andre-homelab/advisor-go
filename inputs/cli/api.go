package cli

import (
	"github.com/andre-felipe-wonsik-alves/inputs/api"
	"github.com/spf13/cobra"
)

var deployApi = &cobra.Command{
	Use:   "api",
	Short: "Sobe a API REST",
	RunE: func(cli *cobra.Command, args []string) error {
		api.Execute()

		return nil
	},
}
