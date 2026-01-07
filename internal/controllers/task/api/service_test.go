package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/models"
)

type fakeStore struct {
	createFn func(ctx context.Context, task *models.Task) error
}

func (f *fakeStore) Create(ctx context.Context, task *models.Task) error {
	if f.createFn == nil {
		return nil
	}
	return f.createFn(ctx, task)
}

func (f *fakeStore) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return nil, nil
}

func (f *fakeStore) List(ctx context.Context) ([]models.Task, error) {
	return nil, nil
}

func (f *fakeStore) Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
	return nil, nil
}

func (f *fakeStore) Delete(ctx context.Context, id string) error {
	return nil
}

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	reminderAt := time.Date(2024, 9, 10, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name      string
		createErr error
		wantErr   bool
	}{
		{
			name:      "creates task with expected fields",
			createErr: nil,
			wantErr:   false,
		},
		{
			name:      "propagates repository error",
			createErr: errors.New("db failure"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var captured *models.Task
			store := &fakeStore{
				createFn: func(ctx context.Context, task *models.Task) error {
					captured = task
					return tt.createErr
				},
			}
			service := NewService(store)

			start := time.Now()
			task, err := service.Create(context.Background(), "titulo", "descricao", models.PriorityHigh, reminderAt)
			end := time.Now()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !errors.Is(err, tt.createErr) {
					t.Fatalf("expected error %v, got %v", tt.createErr, err)
				}
				if task != nil {
					t.Fatalf("expected nil task, got %+v", task)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if task == nil {
				t.Fatal("expected task, got nil")
			}
			if captured == nil {
				t.Fatal("expected repository Create to receive task")
			}
			if task != captured {
				t.Fatalf("expected returned task to match repository argument")
			}
			if task.Title != "titulo" {
				t.Fatalf("expected title %q, got %q", "titulo", task.Title)
			}
			if task.Description != "descricao" {
				t.Fatalf("expected description %q, got %q", "descricao", task.Description)
			}
			if task.Priority != models.PriorityHigh {
				t.Fatalf("expected priority %q, got %q", models.PriorityHigh, task.Priority)
			}
			if !task.ReminderAt.Equal(reminderAt) {
				t.Fatalf("expected reminderAt %v, got %v", reminderAt, task.ReminderAt)
			}
			if task.Done {
				t.Fatalf("expected done to be false")
			}
			if task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
				t.Fatalf("expected CreatedAt/UpdatedAt to be set")
			}
			if task.CreatedAt.Before(start) || task.CreatedAt.After(end) {
				t.Fatalf("expected CreatedAt within test window")
			}
			if task.UpdatedAt.Before(start) || task.UpdatedAt.After(end) {
				t.Fatalf("expected UpdatedAt within test window")
			}
			if task.UpdatedAt.Before(task.CreatedAt) {
				t.Fatalf("expected UpdatedAt to be >= CreatedAt")
			}
		})
	}
}
