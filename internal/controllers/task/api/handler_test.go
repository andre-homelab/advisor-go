package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/models"
	"github.com/go-chi/chi/v5"
)

type stubStore struct {
	createFn func(ctx context.Context, task *models.Task) error
	getFn    func(ctx context.Context, id string) (*models.Task, error)
	listFn   func(ctx context.Context) ([]models.Task, error)
	patchFn  func(ctx context.Context, id string, changes map[string]any) (*models.Task, error)
	deleteFn func(ctx context.Context, id string) error

	lastCreated     *models.Task
	lastPatchID     string
	lastPatchChange map[string]any
}

func (s *stubStore) Create(ctx context.Context, task *models.Task) error {
	s.lastCreated = task
	if s.createFn != nil {
		return s.createFn(ctx, task)
	}
	return nil
}

func (s *stubStore) GetByID(ctx context.Context, id string) (*models.Task, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (s *stubStore) List(ctx context.Context) ([]models.Task, error) {
	if s.listFn != nil {
		return s.listFn(ctx)
	}
	return nil, errors.New("not implemented")
}

func (s *stubStore) Patch(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
	s.lastPatchID = id
	s.lastPatchChange = changes
	if s.patchFn != nil {
		return s.patchFn(ctx, id, changes)
	}
	return nil, nil
}

func (s *stubStore) Delete(ctx context.Context, id string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}

func newRequestWithID(method, path, id string, body *bytes.Reader) *http.Request {
	req := httptest.NewRequest(method, path, body)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", id)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) ErrorResponse {
	t.Helper()
	var errResp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	return errResp
}

func TestTaskHandler_ListTasks(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		store := &stubStore{
			listFn: func(ctx context.Context) ([]models.Task, error) {
				return []models.Task{{ID: "1", Title: "A"}}, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

		handler.ListTasks(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var got []models.Task
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(got) != 1 || got[0].ID != "1" {
			t.Fatalf("unexpected tasks: %+v", got)
		}
	})

	t.Run("error", func(t *testing.T) {
		store := &stubStore{
			listFn: func(ctx context.Context) ([]models.Task, error) {
				return nil, errors.New("db down")
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

		handler.ListTasks(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "Erro ao carregar tarefas" {
			t.Fatalf("error = %q, want %q", errResp.Error, "Erro ao carregar tarefas")
		}
	})
}

func TestTaskHandler_ListSubtasks(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		store := &stubStore{
			getFn: func(ctx context.Context, id string) (*models.Task, error) {
				return nil, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodGet, "/task/1/subtasks", "1", bytes.NewReader(nil))

		handler.ListSubtasks(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "Tarefa não encontrada" {
			t.Fatalf("error = %q, want %q", errResp.Error, "Tarefa não encontrada")
		}
	})

	t.Run("success", func(t *testing.T) {
		parentID := "parent-1"
		store := &stubStore{
			getFn: func(ctx context.Context, id string) (*models.Task, error) {
				return &models.Task{
					ID: parentID,
					Children: []models.Task{
						{ID: "child-1", Title: "A"},
						{ID: "child-2", Title: "B"},
					},
				}, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodGet, "/task/parent-1/subtasks", parentID, bytes.NewReader(nil))

		handler.ListSubtasks(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		var got []models.Task
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(got) != 2 || got[0].ID != "child-1" || got[1].ID != "child-2" {
			t.Fatalf("unexpected subtasks: %+v", got)
		}
	})
}

func TestTaskHandler_CreateTask(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		store := &stubStore{}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("{"))

		handler.CreateTask(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "JSON inválido" {
			t.Fatalf("error = %q, want %q", errResp.Error, "JSON inválido")
		}
	})

	t.Run("missing title", func(t *testing.T) {
		store := &stubStore{}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"","priority":"low"}`))

		handler.CreateTask(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "Título é obrigatório" {
			t.Fatalf("error = %q, want %q", errResp.Error, "Título é obrigatório")
		}
	})

	t.Run("invalid priority", func(t *testing.T) {
		store := &stubStore{}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(`{"title":"A","priority":"invalid"}`))

		handler.CreateTask(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "Prioridade inválida" {
			t.Fatalf("error = %q, want %q", errResp.Error, "Prioridade inválida")
		}
	})

	t.Run("success", func(t *testing.T) {
		store := &stubStore{
			createFn: func(ctx context.Context, task *models.Task) error {
				task.ID = "task-1"
				return nil
			},
			getFn: func(ctx context.Context, id string) (*models.Task, error) {
				return &models.Task{ID: id, Title: "A"}, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		reminder := time.Date(2025, 6, 1, 10, 0, 0, 0, time.UTC)
		body := `{"title":"A","description":"B","priority":"high","reminder_at":"` + reminder.Format(time.RFC3339) + `"}`
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))

		handler.CreateTask(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
		}
		if store.lastCreated == nil || store.lastCreated.Priority != models.PriorityHigh {
			t.Fatalf("priority in create = %v, want %v", store.lastCreated.Priority, models.PriorityHigh)
		}
		var got models.Task
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.ID != "task-1" || got.Title != "A" {
			t.Fatalf("unexpected task: %+v", got)
		}
	})
}

func TestTaskHandler_GetTask(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		store := &stubStore{
			getFn: func(ctx context.Context, id string) (*models.Task, error) {
				return nil, errors.New("missing")
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodGet, "/tasks/1", "1", bytes.NewReader(nil))

		handler.GetTask(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "Tarefa não encontrada" {
			t.Fatalf("error = %q, want %q", errResp.Error, "Tarefa não encontrada")
		}
	})

	t.Run("success", func(t *testing.T) {
		store := &stubStore{
			getFn: func(ctx context.Context, id string) (*models.Task, error) {
				return &models.Task{ID: id, Title: "A"}, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodGet, "/tasks/1", "1", bytes.NewReader(nil))

		handler.GetTask(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		var got models.Task
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.ID != "1" {
			t.Fatalf("id = %q, want %q", got.ID, "1")
		}
	})
}

func TestTaskHandler_PatchTask(t *testing.T) {
	t.Run("invalid json", func(t *testing.T) {
		store := &stubStore{}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodPatch, "/tasks/1", "1", bytes.NewReader([]byte("{")))

		handler.PatchTask(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if !strings.Contains(rec.Body.String(), "JSON inválido") {
			t.Fatalf("body = %q, want JSON inválido", rec.Body.String())
		}
	})

	t.Run("no fields", func(t *testing.T) {
		store := &stubStore{}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodPatch, "/tasks/1", "1", bytes.NewReader([]byte(`{}`)))

		handler.PatchTask(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
		if !strings.Contains(rec.Body.String(), "nenhum campo para atualizar") {
			t.Fatalf("body = %q, want nenhum campo para atualizar", rec.Body.String())
		}
	})

	t.Run("patch error", func(t *testing.T) {
		store := &stubStore{
			patchFn: func(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
				return nil, errors.New("boom")
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodPatch, "/tasks/1", "1", bytes.NewReader([]byte(`{"title":"A"}`)))

		handler.PatchTask(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		if !strings.Contains(rec.Body.String(), "erro ao atualizar") {
			t.Fatalf("body = %q, want erro ao atualizar", rec.Body.String())
		}
	})

	t.Run("not found", func(t *testing.T) {
		store := &stubStore{
			patchFn: func(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
				return nil, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodPatch, "/tasks/1", "1", bytes.NewReader([]byte(`{"title":"A"}`)))

		handler.PatchTask(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
		if !strings.Contains(rec.Body.String(), "não encontrado") {
			t.Fatalf("body = %q, want não encontrado", rec.Body.String())
		}
	})

	t.Run("success", func(t *testing.T) {
		store := &stubStore{
			patchFn: func(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
				return &models.Task{ID: id, Title: "Updated"}, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodPatch, "/tasks/1", "1", bytes.NewReader([]byte(`{"title":"Updated","priority":"high"}`)))

		handler.PatchTask(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		if store.lastPatchID != "1" {
			t.Fatalf("patch id = %q, want %q", store.lastPatchID, "1")
		}
		if store.lastPatchChange["title"] != "Updated" {
			t.Fatalf("title change = %v, want %v", store.lastPatchChange["title"], "Updated")
		}
		if store.lastPatchChange["priority"] != "high" {
			t.Fatalf("priority change = %v, want %v", store.lastPatchChange["priority"], "high")
		}
	})
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		store := &stubStore{
			deleteFn: func(ctx context.Context, id string) error {
				return errors.New("db down")
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodDelete, "/tasks/1", "1", bytes.NewReader(nil))

		handler.DeleteTask(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "Erro ao deletar tarefa" {
			t.Fatalf("error = %q, want %q", errResp.Error, "Erro ao deletar tarefa")
		}
	})

	t.Run("success", func(t *testing.T) {
		store := &stubStore{
			deleteFn: func(ctx context.Context, id string) error {
				return nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodDelete, "/tasks/1", "1", bytes.NewReader(nil))

		handler.DeleteTask(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
		}
	})
}

func TestTaskHandler_CompleteTask(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		store := &stubStore{
			patchFn: func(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
				return nil, errors.New("db down")
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodPatch, "/tasks/1/complete", "1", bytes.NewReader(nil))

		handler.CompleteTask(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		errResp := decodeError(t, rec)
		if errResp.Error != "Erro ao completar tarefa" {
			t.Fatalf("error = %q, want %q", errResp.Error, "Erro ao completar tarefa")
		}
	})

	t.Run("success", func(t *testing.T) {
		store := &stubStore{
			patchFn: func(ctx context.Context, id string, changes map[string]any) (*models.Task, error) {
				return &models.Task{ID: id, Done: true}, nil
			},
		}
		handler := NewTaskHandler(NewService(store))

		rec := httptest.NewRecorder()
		req := newRequestWithID(http.MethodPatch, "/tasks/1/complete", "1", bytes.NewReader(nil))

		handler.CompleteTask(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		var got models.Task
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if !got.Done {
			t.Fatalf("done = %v, want true", got.Done)
		}
	})
}
