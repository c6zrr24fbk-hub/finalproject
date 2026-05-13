package server

import (
	"final_project/pkg/api"
	"log"
	"net/http"
	"os"
)

func StartServer(webDir string) error {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}
	
	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost:%s", port)
	return http.ListenAndServe(":"+port, nil)
}