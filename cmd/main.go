package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

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
	svc := service.NewTaskService(repo)
	hdl := handler.NewTaskHandler(svc)

	//Config das rotas
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/tasks", hdl.CreateTask).Methods("POST")
	router.HandleFunc("/api/v1/tasks", hdl.ListTask).Methods("GET")
	router.HandleFunc("/api/v1/tasks/{id}", hdl.GetTask).Methods("GET")
	router.HandleFunc("/api/v1/tasks/{id}", hdl.UpdateTask).Methods("PUT")
	router.HandleFunc("/api/v1/tasks/{id}", hdl.DeleteTask).Methods("DELETE")
	router.HandleFunc("/api/v1/tasks/{id}/complete", hdl.CompleteTask).Methods("PATCH")

	// Roda servidor
	log.Println("Servidor rodando em :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Erro ao rodar servidor: %v", err)
	}
}
