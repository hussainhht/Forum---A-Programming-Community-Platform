package query

import (
	"database/sql"
	"fmt"
	"forum/Internal/model"
	"time"
)

var (
	deleteSessionQuery         = "DELETE FROM sessions WHERE user_id = ?"
	createSessionQuery         = "INSERT INTO sessions (user_id, session_id, expires_at) VALUES (?, ?, ?)"
	selectUserFromSessionQuery = "SELECT user_id FROM sessions WHERE session_id = ?"
)

func RemoveSession(db *sql.DB, userID int) error {
	_, err := db.Exec(deleteSessionQuery, userID)
	if err != nil {
		return fmt.Errorf("error removing session: %v", err)
	}
	return nil
}

func CreateSession(db *sql.DB, userID int, sessionID string, expiresAt time.Time) error {
	_, err := db.Exec(createSessionQuery, userID, sessionID, expiresAt)
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}
	return nil
}

func SelectUserFromSession(db *sql.DB, sessionID string) (*model.User, error) {
	var userID int
	err := db.QueryRow(selectUserFromSessionQuery, sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error selecting user from session: %v", err)
	}
	
	// Get full user details
	user, err := GetUserByID(db, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user details: %v", err)
	}
	
	return user, nil
}

func GetUserIDFromSession(db *sql.DB, sessionID string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("error getting user ID from session: %v", err)
	}
	return userID, nil
}
