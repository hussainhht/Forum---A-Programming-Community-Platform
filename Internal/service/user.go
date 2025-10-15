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

func (service *Service) RegisterUser(username, email, password, confirmPassword, fullname, photo string) error {

	// Validate input
	if username == "" || email == "" || password == "" || confirmPassword == "" || fullname == "" {
		return fmt.Errorf("all fields are required")
	}

	// Check username length
	if len(username) < 3 || len(username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}

	// Validate email format
	if !service.IsValidEmail(email) {
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

	if photo == "" {
		photo = "default.png"
	}

	err = query.InsertUser(service.DB, username, email, hashPassword(password), fullname, photo)
	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}

	return nil
}

func (service *Service) LogoutUser(session_id string) {
	user, _ := service.GetUserFromSessionID(session_id)
	if user.ID != 0 {
		_ = query.RemoveSession(service.DB, user.ID)
	}

}

// UpdateUserFullName updates the full name for a user
func (service *Service) UpdateUserFullName(userID int, fullname string) error {
	if fullname == "" {
		return fmt.Errorf("full name cannot be empty")
	}
	if len(fullname) < 2 || len(fullname) > 100 {
		return fmt.Errorf("full name must be between 2 and 100 characters")
	}
	return query.UpdateUserFullName(service.DB, userID, fullname)
}

// UpdateUserPassword updates the password for a user after verifying current password
func (service *Service) UpdateUserPassword(userID int, currentPassword, newPassword string) error {
	// Get user
	user, err := query.GetUserByID(service.DB, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if len(newPassword) < 6 {
		return fmt.Errorf("new password must be at least 6 characters")
	}

	// Hash new password
	hashedPassword := hashPassword(newPassword)

	// Update password
	return query.UpdateUserPassword(service.DB, userID, hashedPassword)
}

// UpdateUserPhoto updates the profile photo for a user
func (service *Service) UpdateUserPhoto(userID int, photoPath string) error {
	if photoPath == "" {
		return fmt.Errorf("photo path cannot be empty")
	}
	return query.UpdateUserPhoto(service.DB, userID, photoPath)
}
