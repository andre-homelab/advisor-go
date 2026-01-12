package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"

	taskApi "github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
	"github.com/andre-felipe-wonsik-alves/internal/models"
	"github.com/spf13/cobra"
)

func NewCompleteCli(service *taskApi.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "complete",
		Short: "Completa uma task por meio do ID.",
		RunE: func(cli *cobra.Command, args []string) error {
			ctx := cli.Context()
			reader := bufio.NewReader(os.Stdin)

			tasks, err := service.List(ctx)

			if err != nil {
				log.Fatalln("Um erro aconteceu durante a busca das Tasks: ", err)
				return nil

			}

			for _, task := range tasks {
				err := showTaskId(task)

				if err != nil {
					log.Fatalln("Um erro aconteceu durante a listagem: ", err)
					return nil
				}
			}

			ID, err := promptNonEmpty(reader, "Digite ID da tarefa (obrigatório): ")
			if err != nil {
				return err
			}

			service.Complete(ctx, ID)
			return nil
		},
	}
}

func showTaskId(task models.Task) error {
	fmt.Println("\n______________________________")
	fmt.Printf("ID: %s \n| > Título: %s\n", task.ID, task.Title)
	fmt.Println()

	return nil
}
