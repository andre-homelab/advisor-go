package task

import (
	"errors"
	"testing"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

func TestParsePriority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    models.Priority
		wantErr bool
	}{
		{name: "pt-br baixa maps to low", input: "baixa", want: models.PriorityLow},
		{name: "pt-br média maps to medium", input: "média", want: models.PriorityMedium},
		{name: "pt-br media maps to medium", input: "media", want: models.PriorityMedium},
		{name: "pt-br alta maps to high", input: "alta", want: models.PriorityHigh},
		{name: "low maps to low", input: "low", want: models.PriorityLow},
		{name: "medium maps to medium", input: "medium", want: models.PriorityMedium},
		{name: "high maps to high", input: "high", want: models.PriorityHigh},
		{name: "mixed case and spaces", input: "  HiGh  ", want: models.PriorityHigh},
		{name: "unknown returns error", input: "123", wantErr: true},
		{name: "empty returns error", input: "", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParsePriority(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParsePriority(%q) expected error, got nil", tt.input)
				}
				if !errors.Is(err, ErrInvalidPriority) {
					t.Fatalf("ParsePriority(%q) error = %v, want %v", tt.input, err, ErrInvalidPriority)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParsePriority(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("ParsePriority(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

}
