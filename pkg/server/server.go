package server

import (
	"final-project/pkg/api"
	"fmt"
	"net/http"
	"os"
)

func Start() error {
	value, exists := os.LookupEnv("TODO_PORT")
	port := "7540"

	if exists && value != "" {
		port = value
	}

	api.Init()
	
	http.Handle("/", http.FileServer(http.Dir("web")))

	fmt.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
