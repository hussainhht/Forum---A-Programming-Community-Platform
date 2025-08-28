package service

import (
	"forum/Internal/query"
	"net/http"
)

func (service *Service) IsValidSession(sessionID string) bool {
	// Implement session validation logic

	if sessionID == "" {
		return false
	}

	user, err := query.SelectUserFromSession(service.DB, sessionID)
	if err != nil || user == nil {
		return false
	}

	return true
}

func (service *Service) GetSessionIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
