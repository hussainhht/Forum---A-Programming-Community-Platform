package service

import (
	"fmt"
	"forum/Internal/query"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

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
	if !service.IsValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}

	// Validate password
	if !service.IsValidPassword(password) {
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

func hashPassword(password string) string {
	// use bcrypt to hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

var emailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[A-Za-z]{2,}$`)

func (service *Service) IsValidPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range pw {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case regexp.MustCompile(`[!@#\$%\^&\*]`).MatchString(string(ch)):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func (service *Service) IsValidEmail(email string) bool {
	return emailRX.MatchString(email)
}
