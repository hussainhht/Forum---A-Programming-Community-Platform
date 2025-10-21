package service

import (
	"fmt"
	"forum/Internal/model"
	"forum/Internal/query"
	"strconv"
)

func (service *Service) CommentReaction(commentID int, userID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return fmt.Errorf("invalid reaction type: %s", reactionType)
	}

	// check if user already interacted with the comment
	existingReaction, err := query.GetCommentReaction(service.DB, commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to check existing reaction: %w", err)
	}

	if existingReaction != nil {
		// If the reaction already exists, update it
		if existingReaction.ReactionType == reactionType {
			// If user clicks the same reaction again, remove the reaction (toggle off)
			return query.DeleteCommentReaction(service.DB, commentID, userID)
		}
		return query.UpdateCommentReaction(service.DB, commentID, userID, reactionType)
	}

	return query.InsertCommentReaction(service.DB, commentID, userID, reactionType)
}

func (service *Service) GetCommentsByPostID(postID int, userID int) ([]model.Comment, error) {
	// check if post exists
	if _, err := query.GetPostByID(service.DB, fmt.Sprintf("%d", postID)); err != nil {
		return nil, fmt.Errorf("post with ID %d does not exist: %w", postID, err)
	}

	var comments []model.Comment

	comments, err := query.GetCommentsByPostID(service.DB, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments for post ID %d: %w", postID, err)
	}

	for i := range comments {
		comments[i].Username, err = query.GetUsernameByUserID(service.DB, comments[i].UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get username for user ID %d: %w", comments[i].UserID, err)
		}

		// Get like and dislike counts
		comments[i].LikeCount, err = query.GetCommentLikeCount(service.DB, comments[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get like count for comment ID %d: %w", comments[i].ID, err)
		}

		comments[i].DislikeCount, err = query.GetCommentDislikeCount(service.DB, comments[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get dislike count for comment ID %d: %w", comments[i].ID, err)
		}

		// Get user's reaction if userID is provided (logged in)
		if userID > 0 {
			reaction, err := query.GetCommentReaction(service.DB, comments[i].ID, userID)
			if err == nil && reaction != nil {
				comments[i].UserLiked = reaction.ReactionType == "like"
				comments[i].UserDisliked = reaction.ReactionType == "dislike"
			}
		}
	}

	return comments, nil
}

func (service *Service) CreateComment(postID string, userID int, content string) error {
	// make sure content isnt empty
	if content == "" {
		return fmt.Errorf("comment content cannot be empty")
	}
	// Convert postID string to int
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
	}
	return query.InsertComment(service.DB, postIDInt, userID, content)
}
