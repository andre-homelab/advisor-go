package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	taskApi "github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task/repository"
	"github.com/andre-felipe-wonsik-alves/internal/database"
	"github.com/andre-felipe-wonsik-alves/internal/misc"
	"github.com/spf13/cobra"
)

func NewRootCli(taskSvc *taskApi.Service) *cobra.Command {
	root := &cobra.Command{
		Use:   "advisor-go",
		Short: "Uma CLI para gerenciar tarefas com lembretes :D",
		Long:  "…",
	}

	root.AddCommand(NewAddCli(taskSvc))
	root.AddCommand(NewListCli(taskSvc))
	root.AddCommand(NewCompleteCli(taskSvc))
	root.AddCommand(NewDeployAPICli(taskSvc))

	return root
}

func Execute() {
	misc.PrintBanner()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect()
	if err != nil {
		fmt.Println("Erro na conexão com o banco:", err)
		os.Exit(1)
	}

	database.AutoMigrate(db)

	repo := repository.NewDBStore(db)
	taskSvc := taskApi.NewService(repo)

	root := NewRootCli(taskSvc)

	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
