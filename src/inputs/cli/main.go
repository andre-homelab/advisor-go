package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var mainCmd = &cobra.Command{
	Use:   "advisor-go",
	Short: "Uma CLI para gerenciar tarefas com lembretes :D",
	Long:  "Este é um projeto criado com o objetivo de estudar em entender a linguagem Go. Também, tem-se como objetivo usar cronjobs e a construção da base para um microserviço responsável por notificações em cron,",
}

func Execute() {
	if err := mainCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	mainCmd.AddCommand(addCmd)
	mainCmd.AddCommand(listCmd)
	mainCmd.AddCommand(dueCmd)
}
