package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/andre-felipe-wonsik-alves/src/docs" // importa docs gerados pelo swag
	"github.com/andre-felipe-wonsik-alves/src/internal/handler"
	"github.com/andre-felipe-wonsik-alves/src/internal/task"
)

// @title           Task Notification API
// @version         1.0
// @description     API REST para gerenciamento de tarefas com lembretes
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte API
// @contact.email  seu-email@exemplo.com

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// Inicializa o store JSON
	taskStore := task.NewJSONStore("data/tasks.json")

	// Inicializa service e handlers
	taskService := task.NewService(taskStore)
	taskHandler := handler.NewTaskHandler(taskService)

	// Configura router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.RequestID)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// Rotas da API
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

	// Health check
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
