package main

import (
	"database/sql"
	"log"
	"os"
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
	db err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar: %s", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Erro ao fazer ping: %v", err)
	}
	log.Printf("Conectado com sucesso")

}
