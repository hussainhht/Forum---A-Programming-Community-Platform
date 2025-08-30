package service

import (
	"fmt"
	"forum/Internal/query"
)

func (service *Service) CommentReaction(commentID int, userID int, reactionType string) error {
	if reactionType != "like" && reactionType != "dislike" {
		return fmt.Errorf("invalid reaction type: %s", reactionType)
	}
	
	return query.InsertCommentReaction(service.DB, commentID, userID, reactionType)
}
