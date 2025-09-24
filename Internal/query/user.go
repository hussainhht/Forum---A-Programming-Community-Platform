package query

import (
	"database/sql"
	"fmt"
	"forum/Internal/model"
)

var (
	GetUserByEmailOrUsernameQuery = "SELECT id, username, email, password, fullname, photo FROM users WHERE email = ? OR username = ? LIMIT 1"
	SelectUserWhereEmailQuery     = "SELECT email FROM users WHERE email = ? LIMIT 1"
	SelectUserWhereUsernameQuery  = "SELECT username FROM users WHERE username = ? LIMIT 1"
	InsertUserQuery               = "INSERT INTO users (username, email, password, fullname, photo) VALUES (?, ?, ?, ?, ?)"
	SelectUserWhereIDQuery        = "SELECT username FROM users WHERE id = ? LIMIT 1"
)

func InsertUser(DB *sql.DB, username, email, password, fullname, photo string) error {
	_, err := DB.Exec(InsertUserQuery, username, email, password, fullname, photo)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func SelectUserWhereUsername(DB *sql.DB, username string) (string, error) {
	var existingUsername string
	err := DB.QueryRow(SelectUserWhereUsernameQuery, username).Scan(&existingUsername)
	if err != nil {
		return "", fmt.Errorf("error selecting user: %v", err)
	}
	return existingUsername, nil
}

func SelectUserWhereEmail(DB *sql.DB, email string) (string, error) {
	var existingEmail string
	err := DB.QueryRow(SelectUserWhereEmailQuery, email).Scan(&existingEmail)
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
		fullname string
		photo    string
	)
	err := DB.QueryRow(GetUserByEmailOrUsernameQuery, identifier, identifier).Scan(&userID, &username, &email, &password, &fullname, &photo)
	if err != nil {
		return model.User{}, fmt.Errorf("error selecting user: %v", err)
	}
	return model.User{
		ID:       userID,
		Username: username,
		Email:    email,
		Password: password,
		FullName: fullname,
		Photo:    photo,
	}, nil
}

func GetUserByID(DB *sql.DB, userID int) (*model.User, error) {
	var user model.User
	err := DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error selecting user by ID: %v", err)
	}
	return &user, nil
}

func GetUsernameByUserID(DB *sql.DB, userID int) (string, error) {
	var username string
	err := DB.QueryRow(SelectUserWhereIDQuery, userID).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("error selecting username: %v", err)
	}
	return username, nil
}
