package task

import (
	"strings"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

func ParsePriority(input string) models.Priority {
	s := strings.ToLower(strings.TrimSpace(input))

	switch s {
	case "medium", "média":
		return models.PriorityMedium
	case "high", "alta":
		return models.PriorityHigh
	default:
		return models.PriorityLow
	}
}
