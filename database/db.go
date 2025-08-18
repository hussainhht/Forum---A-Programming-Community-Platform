package database

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// Make sure migration_history table exists first
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS migration_history (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT NOT NULL UNIQUE,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		log.Fatal(err)
	}

	Migrate()
}
