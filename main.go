package main

import (
	"final-project/pkg/db"
	"final-project/pkg/server"
	"log"
	"os"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	err := db.Init(dbFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start())
}
