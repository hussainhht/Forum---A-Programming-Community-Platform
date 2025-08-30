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
