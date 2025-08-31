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

func (s *Service) GetPostByID(postID string) (*model.Post, error) {
	var post *model.Post
	var err error
	post, err = query.GetPostByID(s.DB, postID)
	if err != nil {
		return nil, err
	}
	post.Username, err = query.GetUsernameByUserID(s.DB, post.UserID)
	if err != nil {
		return nil, err
	}
	fmt.Println(post.Username)
	post.LikeCount, err = query.GetPostLikeCount(s.DB, post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments, err = s.GetCommentsByPostID(post.ID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Service) PostReaction(postID int, userID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return fmt.Errorf("invalid reaction type: %s", reactionType)
	}

	// check if user already interacted with the post
	existingReaction, err := query.GetPostReaction(s.DB, postID, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing reaction: %w", err)
	}

	if existingReaction != nil {
		// If the reaction already exists, update it
		if existingReaction.ReactionType == reactionType {
			return nil // No change needed
		}
		return query.UpdatePostReaction(s.DB, postID, userID, reactionType)
	}

	return query.InsertPostReaction(s.DB, postID, userID, reactionType)
}
