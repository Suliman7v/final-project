package main

import (
	"fmt"
	"log"
	"os"

	"final-project/pkg/db"
	"final-project/pkg/server"
)

func main() {
	connStr := os.Getenv("TODO_DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:password@localhost:5432/todo?sslmode=disable"
	}

	err := db.Init(connStr)
	if err != nil {
		log.Printf("Ошибка инициализации БД: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("Сервер запускается...")

	err = server.Start()
	if err != nil {
		log.Printf("Ошибка запуска сервера: %v", err)
		os.Exit(1)
	}
}
