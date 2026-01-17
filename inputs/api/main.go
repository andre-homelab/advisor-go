package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/andre-felipe-wonsik-alves/docs"
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
	taskApi "github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
	"github.com/andre-felipe-wonsik-alves/internal/database"
	"github.com/andre-felipe-wonsik-alves/internal/misc"
)

// @title           Task Notification API
// @version         1.0
// @description     API REST para gerenciamento de tarefas com lembretes
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

func Execute(ctx context.Context, service *taskApi.Service) error {
	misc.PrintBanner()

	log.Println("Testando conexão com o banco de dados")
	db, err := database.Connect()

	if err != nil {
		log.Fatal("Erro na conexão com o banco de dados: ", err)
	}

	log.Println("Conexão com o banco estabelecida!")

	log.Println("Iniciando migrations...")
	database.AutoMigrate(db)

	log.Println("Migration concluída com sucesso!")

	taskHandler := api.NewTaskHandler(service)

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
			r.Patch("/{id}", taskHandler.PatchTask)
			r.Delete("/{id}", taskHandler.DeleteTask)
			r.Patch("/{id}/complete", taskHandler.CompleteTask)
			r.Get("/{id}/subtasks", taskHandler.ListSubtasks)

		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Servidor rodando em http://localhost:8080")
	log.Println("Documentação disponível em http://localhost:8080/swagger/index.html")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	fmt.Print("\n\n")
	return nil
}
