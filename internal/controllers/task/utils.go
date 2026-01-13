package task

import (
	"errors"
	"strings"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

var ErrInvalidPriority = errors.New("prioridade inválida")

func ParsePriority(input string) (models.Priority, error) {
	s := strings.ToLower(strings.TrimSpace(input))

	switch s {
	case "low", "baixa":
		return models.PriorityLow, nil
	case "medium", "media", "média":
		return models.PriorityMedium, nil
	case "high", "alta":
		return models.PriorityHigh, nil
	default:
		return "", ErrInvalidPriority
	}
}
