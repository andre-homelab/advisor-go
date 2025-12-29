package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/andre-felipe-wonsik-alves/docs"
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task"
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
)

// @title           Task Notification API
// @version         1.0
// @description     API REST para gerenciamento de tarefas com lembretes
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	taskStore := task.NewJSONStore("data/tasks.json")

	taskService := api.NewService(taskStore)
	taskHandler := api.NewTaskHandler(taskService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RequestID)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", taskHandler.ListTasks)
			r.Post("/", taskHandler.CreateTask)
			r.Get("/{id}", taskHandler.GetTask)
			r.Put("/{id}", taskHandler.UpdateTask)
			r.Delete("/{id}", taskHandler.DeleteTask)
			r.Patch("/{id}/complete", taskHandler.CompleteTask)
			r.Get("/due", taskHandler.GetDueTasks)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Println("Documentação disponível em http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
