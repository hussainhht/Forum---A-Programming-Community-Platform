package service

import (
	"fmt"
	"forum/Internal/model"
	"forum/Internal/query"
)

func (service *Service) CreatePost(sessionID, title, content string, categories []string) error {
	DB := service.DB

	// get user from session
	userID, err := query.GetUserIDFromSession(DB, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get user ID from session: %w", err)
	}

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = query.InsertPost(tx, title, content, categories, userID)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAllPosts retrieves all posts for display on home page
func (s *Service) GetAllPosts() ([]model.Post, error) {
	return query.GetAllPosts(s.DB)
}
