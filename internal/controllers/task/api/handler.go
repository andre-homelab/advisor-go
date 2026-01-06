package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andre-felipe-wonsik-alves/internal/controllers/task"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskService Service
}

func NewTaskHandler(taskService *Service) *TaskHandler {
	return &TaskHandler{taskService: *taskService}
}

type CreateTaskRequest struct {
	Title       string    `json:"title" example:"Reunião importante"`
	Description string    `json:"description" example:"Apresentar projeto ao time"`
	Priority    string    `json:"priority" example:"high" enums:"low,medium,high"`
	ReminderAt  time.Time `json:"reminder_at" example:"2025-12-27T15:00:00Z"`
}

type PatchTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Priority    *string    `json:"priority,omitempty" enums:"low,medium,high"`
	ReminderAt  *time.Time `json:"reminder_at,omitempty"`
	Done        *bool      `json:"done,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Tarefa não encontrada"`
	Message string `json:"message,omitempty" example:"ID inválido fornecido"`
}

// @Summary     Listar todas as tarefas
// @Description Retorna lista de todas as tarefas cadastradas
// @Tags        tasks
// @Produce     json
// @Success     200 {array} models.Task
// @Failure     500 {object} ErrorResponse
// @Router      /tasks [get]
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := h.taskService.List(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao carregar tarefas", err)
		return
	}
	respondJSON(w, http.StatusOK, tasks)
}

// @Summary     Criar nova tarefa
// @Description Adiciona uma nova tarefa ao sistema
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       task body CreateTaskRequest true "Dados da tarefa"
// @Success     201 {object} models.Task
// @Failure     400 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "JSON inválido", err)
		return
	}

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "Título é obrigatório", nil)
		return
	}

	priority := task.ParsePriority(req.Priority)

	newTask, err := h.taskService.Create(r.Context(), req.Title, req.Description, priority, req.ReminderAt)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Erro ao criar tarefa", err)
		return
	}

	respondJSON(w, http.StatusCreated, newTask)
}

// @Summary     Buscar tarefa por ID
// @Description Retorna uma tarefa específica pelo ID
// @Tags        tasks
// @Produce     json
// @Param       id path string true "ID da tarefa"
// @Success     200 {object} models.Task
// @Failure     404 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /tasks/{id} [get]
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	t, err := h.taskService.GetByID(r.Context(), id)
	if err != nil {
		if err == ErrTaskNotFound {
			respondError(w, http.StatusNotFound, "Tarefa não encontrada", nil)
			return
		}

		respondError(w, http.StatusInternalServerError, "Erro ao buscar tarefa", err)
		return
	}

	respondJSON(w, http.StatusOK, t)
}

// @Summary     Atualizar campos específicos de uma tarefa
// @Description Atualiza dados de uma tarefa existente
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id path string true "ID da tarefa"
// @Param       task body PatchTaskRequest true "Dados para atualização"
// @Success     200 {object} models.Task
// @Failure     400 {object} ErrorResponse
// @Failure     404 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /tasks/{id} [patch]
func (h *TaskHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req PatchTaskRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	changes := map[string]any{}

	if req.Title != nil {
		changes["title"] = *req.Title
	}
	if req.Description != nil {
		changes["description"] = *req.Description
	}
	if req.Done != nil {
		changes["done"] = *req.Done
	}
	if req.Priority != nil {
		changes["priority"] = *req.Priority
	}
	if req.ReminderAt != nil {
		changes["reminderAt"] = *req.ReminderAt
	}

	if len(changes) == 0 {
		http.Error(w, "nenhum campo para atualizar", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.Patch(r.Context(), id, changes)

	if err != nil {
		http.Error(w, "erro ao atualizar", http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// @Summary     Deletar tarefa
// @Description Remove uma tarefa do sistema
// @Tags        tasks
// @Produce     json
// @Param       id path string true "ID da tarefa"
// @Success     204 "Tarefa removida com sucesso"
// @Failure     404 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.taskService.Delete(r.Context(), id)
	if err != nil {
		if err == ErrTaskNotFound {
			respondError(w, http.StatusNotFound, "Tarefa não encontrada", nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "Erro ao deletar tarefa", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Summary     Marcar tarefa como concluída
// @Description Marca uma tarefa específica como concluída
// @Tags        tasks
// @Produce     json
// @Param       id path string true "ID da tarefa"
// @Success     200 {object} models.Task
// @Failure     404 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /tasks/{id}/complete [patch]
func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	completed, err := h.taskService.Complete(r.Context(), id)
	if err != nil {
		if err == ErrTaskNotFound {
			respondError(w, http.StatusNotFound, "Tarefa não encontrada", nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "Erro ao completar tarefa", err)
		return
	}

	respondJSON(w, http.StatusOK, completed)
}

// @Summary     Listar tarefas vencidas
// @Description Retorna todas as tarefas cujo lembrete já passou
// @Tags        tasks
// @Produce     json
// @Success     200 {array} models.Task
// @Failure     500 {object} ErrorResponse
// @Router      /tasks/due [get]
// func (h *TaskHandler) GetDueTasks(w http.ResponseWriter, r *http.Request) {
// 	dueTasks, err := h.taskService.GetDue()
// 	if err != nil {
// 		respondError(w, http.StatusInternalServerError, "Erro ao buscar tarefas vencidas", err)
// 		return
// 	}

// 	respondJSON(w, http.StatusOK, dueTasks)
// }

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string, err error) {
	errResp := ErrorResponse{Error: message}
	if err != nil {
		errResp.Message = err.Error()
	}
	respondJSON(w, status, errResp)
}
