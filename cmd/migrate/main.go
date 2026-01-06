package main

import (
	"log"

	"github.com/andre-felipe-wonsik-alves/internal/database"
)

func main() {
	log.Println("Conectando ao banco de dados...")
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Erro na conexão com o banco de dados: ", err)
	}

	log.Println("Executando migrations...")
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Erro ao executar migrations: ", err)
	}

	log.Println("Migrations concluídas com sucesso!")
}
