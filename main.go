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
		log.Fatal("Ошибка инициализации БД:", err)
	}

	webDir := "./web"
	if err := server.StartServer(webDir); err != nil {
		log.Fatal(err)
	}
}