package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

var schema = []string{
	`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);`,
	`CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`,
}

func Init(dbFile string) error {
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	for _, query := range schema {
		if _, err := db.Exec(query); err != nil {
			db.Close()
			return err
		}
	}

	DB = db
	log.Println("База данных инициализирована")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}