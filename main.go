package main

import (
	"fmt"
	"log"
	"os"

	"final-project/pkg/db"
	"final-project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
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
