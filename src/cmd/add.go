package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/andre-felipe-wonsik-alves/advisor-go/internal/task"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adiciona uma nova tarefa",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		title, err := promptNonEmpty(reader, "Título da tarefa (obrigatório): ")

		if err != nil {
			return err
		}

		description, err := prompt(reader, "Descrição (opcional): ")

		if err != nil {
			return err
		}

		priorityStr, err := prompt(reader, "Prioridade da tarefa (média, alta ou baixa): ")

		if err != nil {
			return err
		}

		priority := parsePriority(priorityStr)

		fmt.Println("\nInforme a data/hora do lembrete no formato:")
		fmt.Println("  02/01/2006 15:04")
		reminderStr, err := promptNonEmpty(reader, "Lembrar em: ")

		if err != nil {
			return err
		}

		reminderAt, err := parseReminder(reminderStr)

		if err != nil {
			return err
		}

		newTask, err := task.Create(title, description, priority, reminderAt)

		if err != nil {
			return err
		}

		fmt.Println("\nTarefa adicionada com sucesso!")
		fmt.Printf("ID: %s\n", newTask.ID)
		fmt.Printf("Título: %s\n", newTask.Title)
		fmt.Printf("Lembrar em: %s\n", newTask.ReminderAt.Format("02/01/2006 15:04"))

		return nil

	},
}

func prompt(reader *bufio.Reader, label string) (string, error) {
	fmt.Print(label)

	text, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

func promptNonEmpty(reader *bufio.Reader, label string) (string, error) {
	for {
		value, err := prompt(reader, label)

		if err != nil {
			return "", err
		}

		if value != "" {
			return value, nil
		}

		fmt.Println("Campo obrigatório!")
	}
}

func parsePriority(input string) task.Priority {
	s := strings.ToLower(strings.TrimSpace(input))

	switch s {
	case "medium", "média":
		return task.PriorityMedium
	case "high", "alta":
		return task.PriorityHigh
	default:
		return task.PriorityLow
	}
}

func parseReminder(input string) (time.Time, error) {
	layout := "02/01/2006 15:04"
	return time.ParseInLocation(layout, input, time.Local)
}
