package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func InsertUser(username, email, password string) error {
	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
	_, err = DB.Exec(query, username, email, string(hashed))
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (int, string, string, error) {
	var id int
	var username, hashedPassword string
	query := `SELECT id, username, password FROM users WHERE lower(email) = lower(?) LIMIT 1`
	err := DB.QueryRow(query, email).Scan(&id, &username, &hashedPassword)
	if err != nil {
		return 0, "", "", err
	}
	return id, username, hashedPassword, nil
}

func InsertComment(postID, userID int, content string) error {
	if content == "" {
		return errors.New("comment cannot be empty")
	}
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, postID, userID, content)
	return err
}

func InsertReaction(userID, postID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return errors.New("invalid reaction type")
	}
	query := `INSERT INTO reactions (user_id, post_id, reaction_type) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, userID, postID, reactionType)
	return err
}

func CreateSession(db *sql.DB, userID int, sessionID string, expiresAt time.Time) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, sessionID, userID, expiresAt)
	return err
}

func GetUserBySession(db *sql.DB, sessionID string) (int, error) {
	var userID int
	var expires string
	err := db.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE id = ?`, sessionID).
		Scan(&userID, &expires)
	if err != nil {
		return 0, err
	}
	// TODO: Check if session is expired
	return userID, nil
}

func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE id = ?`, sessionID)
	return err
}
