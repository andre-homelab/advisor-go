package cli

import (
	"github.com/andre-felipe-wonsik-alves/inputs/api"
	taskApi "github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
	"github.com/spf13/cobra"
)

func NewDeployAPICli(service *taskApi.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Sobe a API REST",
		RunE: func(cli *cobra.Command, args []string) error {
			return api.Execute(cli.Context(), service)
		},
	}
}
