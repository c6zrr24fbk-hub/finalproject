package main

import (
	"final_project/pkg/db"
	"final_project/pkg/server"
	"log"
	"os"
)

func main() {
	dbFile := "scheduler.db"
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.Close()

	webDir := "./web"
	if err := server.StartServer(webDir); err != nil {
		log.Printf("Ошибка запуска сервера: %v", err)
		os.Exit(1)
	}
}