package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var mainCli = &cobra.Command{
	Use:   "advisor-go",
	Short: "Uma CLI para gerenciar tarefas com lembretes :D",
	Long:  "Este é um projeto criado com o objetivo de estudar em entender a linguagem Go. Também, tem-se como objetivo usar cronjobs e a construção da base para um microserviço responsável por notificações em cron,",
}

func Execute() {
	if err := mainCli.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	mainCli.AddCommand(addCli)
	mainCli.AddCommand(listCli)
	mainCli.AddCommand(dueCli)
	mainCli.AddCommand(deployApi)
}
