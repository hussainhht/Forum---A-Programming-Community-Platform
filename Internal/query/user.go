package query

import (
	"database/sql"
	"fmt"
	"forum/Internal/model"
)

var (
	GetUserByEmailOrUsername = "SELECT id, username, email, password FROM users WHERE email = ? OR username = ? LIMIT 1"
	selectUserWhereEmail     = "SELECT email FROM users WHERE email = ? LIMIT 1"
	selectUserWhereUsername  = "SELECT username FROM users WHERE username = ? LIMIT 1"
	insertUserQuery          = "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
)

func InsertUser(DB *sql.DB, username, email, password string) error {
	_, err := DB.Exec(insertUserQuery, username, email, password)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func SelectUserWhereUsername(DB *sql.DB, username string) (string, error) {
	var existingUsername string
	err := DB.QueryRow(selectUserWhereUsername, username).Scan(&existingUsername)
	if err != nil {
		return "", fmt.Errorf("error selecting user: %v", err)
	}
	return existingUsername, nil
}

func SelectUserWhereEmail(DB *sql.DB, email string) (string, error) {
	var existingEmail string
	err := DB.QueryRow(selectUserWhereEmail, email).Scan(&existingEmail)
	if err != nil {
		return "", err
	}
	return existingEmail, nil
}

func GetUserByUsernameOrEmail(DB *sql.DB, identifier string) (model.User, error) {
	var (
		userID   int
		username string
		email    string
		password string
	)
	err := DB.QueryRow(GetUserByEmailOrUsername, identifier, identifier).Scan(&userID, &username, &email, &password)
	if err != nil {
		return model.User{}, fmt.Errorf("error selecting user: %v", err)
	}
	return model.User{
		ID:       userID,
		Username: username,
		Email:    email,
		Password: password,
	}, nil
}
