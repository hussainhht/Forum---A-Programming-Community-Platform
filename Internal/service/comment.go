package service

import (
	"fmt"
	"forum/Internal/model"
	"forum/Internal/query"
)

func (service *Service) CommentReaction(commentID int, userID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return fmt.Errorf("invalid reaction type: %s", reactionType)
	}

	return query.InsertCommentReaction(service.DB, commentID, userID, reactionType)
}

func (service *Service) GetCommentsByPostID(postID int) ([]model.Comment, error) {
	// check if post exists
	if _, err := query.GetPostByID(service.DB, fmt.Sprintf("%d", postID)); err != nil {
		return nil, fmt.Errorf("post with ID %d does not exist: %w", postID, err)
	}

	return query.GetCommentsByPostID(service.DB, postID)
}

func (service *Service) CreateComment(sessionID string, postID int, content string) error {
	userID, err := service.GetUserIDFromSessionID(sessionID)
	if err != nil {
		return fmt.Errorf("failed to get user ID from session: %w", err)
	}

	return query.InsertComment(service.DB, postID, userID, content)
}
