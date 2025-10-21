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
	err = service.validateCategories(categories)
	if err != nil {
		return fmt.Errorf("failed to validate categories: %w", err)
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

// GetPostsByCategories retrieves posts filtered by selected categories
func (s *Service) GetPostsByCategories(categoryIDs []string) ([]model.Post, error) {
	return query.GetPostsByCategories(s.DB, categoryIDs)
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
	post.LikeCount, err = query.GetPostLikeCount(s.DB, post.ID)
	if err != nil {
		return nil, err
	}
	post.DislikeCount, err = query.GetPostDislikeCount(s.DB, post.ID)
	if err != nil {
		return nil, err
	}
	// Pass 0 for userID since we'll populate user reactions separately in the handler
	post.Comments, err = s.GetCommentsByPostID(post.ID, 0)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *Service) GetPostByIDWithUser(postID string, userID int) (*model.Post, error) {
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
	post.LikeCount, err = query.GetPostLikeCount(s.DB, post.ID)
	if err != nil {
		return nil, err
	}
	post.DislikeCount, err = query.GetPostDislikeCount(s.DB, post.ID)
	if err != nil {
		return nil, err
	}
	// Pass userID to populate user reactions
	post.Comments, err = s.GetCommentsByPostID(post.ID, userID)
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
			// If user clicks the same reaction again, remove the reaction (toggle off)
			return query.DeletePostReaction(s.DB, postID, userID)
		}
		return query.UpdatePostReaction(s.DB, postID, userID, reactionType)
	}

	return query.InsertPostReaction(s.DB, postID, userID, reactionType)
}

func (s *Service) GetFeaturedPosts() ([]model.Post, error) {
	return query.GetFeaturedPosts(s.DB)
}
func (s *Service) GetLatestPosts() ([]model.Post, error) {
	return query.GetLatestPosts(s.DB)
}

func (s *Service) GetDiscoverPosts(search, category, sort string, limit, offset int) ([]model.Post, bool, error) {
	return query.GetDiscoverPosts(s.DB, search, category, sort, limit, offset)
}

// GetDiscoverPostsMultiCategory retrieves posts filtered by multiple categories
func (s *Service) GetDiscoverPostsMultiCategory(search string, categories []string, sort string, limit, offset int) ([]model.Post, bool, error) {
	return query.GetDiscoverPostsMultiCategory(s.DB, search, categories, sort, limit, offset)
}

func (s *Service) GetPostReaction(postID, userID int) (*model.PostReaction, error) {
	return query.GetPostReaction(s.DB, postID, userID)
}
