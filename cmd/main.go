package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"log"
	"net/http"
	"os"

	"github.com/DinizJ/desafio/internal/handler"
	"github.com/DinizJ/desafio/internal/repository"
	"github.com/DinizJ/desafio/internal/service"
)

func main() {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	log.Printf("Conectando em: %s", dsn)

	//conectar ao banco
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar: %s", err)
	}
	defer db.Close()

	//Ping
	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao fazer ping: %v", err)
	}
	log.Printf("Conectado com sucesso")

	//Inicializa as layers
	repo := repository.NewTaskRepository(db)
	svc := service.TaskService(repo)
	hdl := handler.TaskHandler(svc)

	//Config das rotas
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/tasks", hdl.CreateTask)
	mux.HandleFunc("GET /api/v1/tasks", hdl.ListTask)
	mux.HandleFunc("GET /api/v1/tasks/{id}", hdl.GetTask)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", hdl.UpdateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", hdl.DeleteTask)
	mux.HandleFunc("PATCH /api/v1/tasks/{id}/complete", hdl.CompleteTask)

	// Roda servidor
	log.Println("Servidor rodando em :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Erro ao rodar servidor: %v", err)
	}
}
