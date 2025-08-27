package database

import (
	"database/sql"
	"log"
	"time"

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

func InsertPost(title, content string, userID int, categoryID int, createdAt time.Time) {
	insertPost := `INSERT INTO posts (title, content, user_id, category_id, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := DB.Exec(insertPost, title, content, userID, categoryID, createdAt)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadPost() {

}
