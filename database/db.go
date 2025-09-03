package database

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite3", "./forum.db")
	if err != nil {
		return nil, err
	}

	schema, err := os.ReadFile("./database/db.sql")
	if err != nil {
		return nil, err
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		return nil, err
	}
	return DB, nil
}
