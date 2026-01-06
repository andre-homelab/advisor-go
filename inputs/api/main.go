package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/andre-felipe-wonsik-alves/docs"
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task/api"
	"github.com/andre-felipe-wonsik-alves/internal/controllers/task/repository"
	"github.com/andre-felipe-wonsik-alves/internal/database"
)

// @title           Task Notification API
// @version         1.0
// @description     API REST para gerenciamento de tarefas com lembretes
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

func Execute() {
	printBanner()

	log.Println("Testando conexão com o banco de dados")
	db, err := database.Connect()

	if err != nil {
		log.Fatal("Erro na conexão com o banco de dados: ", err)
	}

	log.Println("Conexão com o banco estabelecida!")

	log.Println("Iniciando migrations...")
	database.AutoMigrate(db)

	log.Println("Migration concluída com sucesso!")

	repo := repository.NewDBStore(db)

	taskService := api.NewService(*repo)
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
			r.Patch("/{id}", taskHandler.PatchTask)
			r.Delete("/{id}", taskHandler.DeleteTask)
			r.Patch("/{id}/complete", taskHandler.CompleteTask)
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
	fmt.Print("\n\n")
}

func printBanner() {
	const banner = `
 $$$$$$\   $$$$$$\           $$$$$$\  $$$$$$$\  $$\    $$\ $$$$$$\  $$$$$$\   $$$$$$\  $$$$$$$\  
$$  __$$\ $$  __$$\         $$  __$$\ $$  __$$\ $$ |   $$ |\_$$  _|$$  __$$\ $$  __$$\ $$  __$$\ 
$$ /  \__|$$ /  $$ |        $$ /  $$ |$$ |  $$ |$$ |   $$ |  $$ |  $$ /  \__|$$ /  $$ |$$ |  $$ |
$$ |$$$$\ $$ |  $$ |$$$$$$\ $$$$$$$$ |$$ |  $$ |\$$\  $$  |  $$ |  \$$$$$$\  $$ |  $$ |$$$$$$$  |
$$ |\_$$ |$$ |  $$ |\______|$$  __$$ |$$ |  $$ | \$$\$$  /   $$ |   \____$$\ $$ |  $$ |$$  __$$< 
$$ |  $$ |$$ |  $$ |        $$ |  $$ |$$ |  $$ |  \$$$  /    $$ |  $$\   $$ |$$ |  $$ |$$ |  $$ |
\$$$$$$  | $$$$$$  |        $$ |  $$ |$$$$$$$  |   \$  /   $$$$$$\ \$$$$$$  | $$$$$$  |$$ |  $$ |
 \______/  \______/         \__|  \__|\_______/     \_/    \______| \______/  \______/ \__|  \__|
`

	fmt.Print("\n\n")
	fmt.Print(banner)
	fmt.Print("\n\n")
}
