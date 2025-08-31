package service

import (
	"fmt"
	"forum/Internal/query"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (service *Service) LoginUser(emailOrUsername, password string) (string, error) {

	existingUser, err := query.GetUserByUsernameOrEmail(service.DB, emailOrUsername)
	if err != nil {
		return "", err
	}
	storedPassword := existingUser.Password
	// compare plain password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials: %v", err)
	}

	// remove any existing sessions for this user
	_ = query.RemoveSession(service.DB, existingUser.ID)

	newSessionID := uuid.New().String()

	// create a new session
	if err := query.CreateSession(service.DB, existingUser.ID, newSessionID, time.Now().Add(24*time.Hour)); err != nil {
		return "", err
	}

	return newSessionID, nil
}

func (service *Service) RegisterUser(username, email, password, confirmPassword string) error {

	// Validate input
	if username == "" || email == "" || password == "" || confirmPassword == "" {
		return fmt.Errorf("all fields are required")
	}

	// Check username length
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	// Validate email format
	if !service.isValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	// Validate password
	if !service.isValidPassword(password) {
		return fmt.Errorf("password must be at least 8 characters and contain uppercase, lowercase, digit, and special character")
	}

	// Check if passwords match
	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	// Check if username already exists
	_, err := query.SelectUserWhereUsername(service.DB, username)
	if err == nil {
		return fmt.Errorf("username already exists")
	}

	// Check if email already exists
	_, err = query.SelectUserWhereEmail(service.DB, email)
	if err == nil {
		return fmt.Errorf("email already exists")
	}

	err = query.InsertUser(service.DB, username, email, hashPassword(password))
	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}

	return nil
}

func (service *Service) LogoutUser(session_id string) {
	userID, _ := service.GetUserIDFromSessionID(session_id)
	if userID != 0 {
		_ = query.RemoveSession(service.DB, userID)
	}

}
