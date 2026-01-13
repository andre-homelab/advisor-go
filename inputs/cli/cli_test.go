package cli

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	taskApi "github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
	"github.com/andre-felipe-wonsik-alves/internal/models"
)

type fakeStore struct {
	createTask  *models.Task
	createErr   error
	listTasks   []models.Task
	listErr     error
	patchID     string
	patchChanges map[string]any
	patchResult *models.Task
	patchErr    error
}

func (f *fakeStore) Create(_ context.Context, task *models.Task) error {
	if task.ID == "" {
		task.ID = "task-123"
	}
	f.createTask = task
	return f.createErr
}

func (f *fakeStore) GetByID(_ context.Context, _ string) (*models.Task, error) {
	return nil, nil
}

func (f *fakeStore) List(_ context.Context) ([]models.Task, error) {
	return f.listTasks, f.listErr
}

func (f *fakeStore) Patch(_ context.Context, id string, changes map[string]any) (*models.Task, error) {
	f.patchID = id
	f.patchChanges = changes
	if f.patchErr != nil {
		return nil, f.patchErr
	}
	if f.patchResult != nil {
		return f.patchResult, nil
	}
	return &models.Task{ID: id, Done: true}, nil
}

func (f *fakeStore) Delete(_ context.Context, _ string) error {
	return nil
}

func withStdin(input string, fn func()) {
	original := os.Stdin
	reader, writer, _ := os.Pipe()
	_, _ = writer.WriteString(input)
	_ = writer.Close()
	os.Stdin = reader
	defer func() {
		os.Stdin = original
		_ = reader.Close()
	}()
	fn()
}

func captureStdout(fn func()) string {
	original := os.Stdout
	reader, writer, _ := os.Pipe()
	os.Stdout = writer
	fn()
	_ = writer.Close()
	os.Stdout = original
	output, _ := io.ReadAll(reader)
	_ = reader.Close()
	return string(output)
}

func TestNewAddCli_RunE_Success(t *testing.T) {
	fakeRepo := &fakeStore{}
	service := taskApi.NewService(fakeRepo)
	cmd := NewAddCli(service)
	cmd.SetContext(context.Background())

	input := strings.Join([]string{
		"Tarefa A",
		"Descricao A",
		"media",
		"02/01/2006 15:04",
		"",
	}, "\n")

	var err error
	var output string
	withStdin(input, func() {
		output = captureStdout(func() {
			err = cmd.Execute()
		})
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fakeRepo.createTask == nil {
		t.Fatalf("expected Create to be called")
	}
	if fakeRepo.createTask.Title != "Tarefa A" {
		t.Fatalf("unexpected title: %s", fakeRepo.createTask.Title)
	}
	if fakeRepo.createTask.Description != "Descricao A" {
		t.Fatalf("unexpected description: %s", fakeRepo.createTask.Description)
	}
	if fakeRepo.createTask.Priority != models.PriorityMedium {
		t.Fatalf("unexpected priority: %s", fakeRepo.createTask.Priority)
	}

	expectedReminder := time.Date(2006, 1, 2, 15, 4, 0, 0, time.Local)
	if !fakeRepo.createTask.ReminderAt.Equal(expectedReminder) {
		t.Fatalf("unexpected reminder: %v", fakeRepo.createTask.ReminderAt)
	}

	if !strings.Contains(output, "Tarefa adicionada com sucesso!") {
		t.Fatalf("expected success message, got %q", output)
	}
	if !strings.Contains(output, "ID: task-123") {
		t.Fatalf("expected task id in output, got %q", output)
	}
}

func TestNewListCli_RunE_Success(t *testing.T) {
	fakeRepo := &fakeStore{
		listTasks: []models.Task{
			{
				ID:          "task-1",
				Title:       "Tarefa 1",
				Description: "Descricao 1",
				Priority:    models.PriorityLow,
				ReminderAt:  time.Date(2025, 1, 5, 10, 30, 0, 0, time.Local),
			},
		},
	}
	service := taskApi.NewService(fakeRepo)
	cmd := NewListCli(service)
	cmd.SetContext(context.Background())

	var err error
	var output string
	output = captureStdout(func() {
		err = cmd.Execute()
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(output, "ID: task-1") {
		t.Fatalf("expected task id in output, got %q", output)
	}
	if !strings.Contains(output, "Título: Tarefa 1") {
		t.Fatalf("expected task title in output, got %q", output)
	}
}

func TestNewCompleteCli_RunE_Success(t *testing.T) {
	fakeRepo := &fakeStore{
		listTasks: []models.Task{
			{ID: "task-1", Title: "Tarefa 1"},
			{ID: "task-2", Title: "Tarefa 2"},
		},
	}
	service := taskApi.NewService(fakeRepo)
	cmd := NewCompleteCli(service)
	cmd.SetContext(context.Background())

	input := strings.Join([]string{
		"task-2",
		"",
	}, "\n")

	var err error
	var output string
	withStdin(input, func() {
		output = captureStdout(func() {
			err = cmd.Execute()
		})
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if fakeRepo.patchID != "task-2" {
		t.Fatalf("expected patch id task-2, got %s", fakeRepo.patchID)
	}
	if done, ok := fakeRepo.patchChanges["done"].(bool); !ok || !done {
		t.Fatalf("expected done=true in patch changes, got %v", fakeRepo.patchChanges)
	}
	if !strings.Contains(output, "ID: task-1") || !strings.Contains(output, "ID: task-2") {
		t.Fatalf("expected tasks to be listed, got %q", output)
	}
}
