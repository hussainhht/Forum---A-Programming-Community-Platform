package service

import (
	"encoding/json"
	"forum/Internal/model"
	"forum/Internal/query"
)


// Check if session exists and is valid
func (s *Service) ValidSessions(sessionID string) bool {
	var count int
	query := `SELECT COUNT(*) FROM sessions WHERE session_id=? AND expires_at > CURRENT_TIMESTAMP`
	_ = s.DB.QueryRow(query, sessionID).Scan(&count)
	return count > 0
}

// Get user by session ID
func (s *Service) Get_UserBySession(sessionID string) *model.User {
	var u model.User
	query := `
		SELECT users.id, users.username, users.email, users.password
		FROM users
		JOIN sessions ON users.id = sessions.user_id
		WHERE sessions.session_id=? AND sessions.expires_at > CURRENT_TIMESTAMP
	`
	err := s.DB.QueryRow(query, sessionID).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil
	}
	return &u
}

// Get posts for a user including like and comment counts
func (s *Service) Get_PostsByUser(userID int) ([]model.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.created_at,
	       (SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type='like') as like_count,
	       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comment_count
	FROM posts p
	WHERE p.user_id = ?
	ORDER BY p.created_at DESC
	`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		var contentJSON string
		if err := rows.Scan(&p.ID, &p.Title, &contentJSON, &p.CreatedAt, &p.LikeCount, &p.CommentCount); err != nil {
			return nil, err
		}
		// Convert JSON content to []Block
		var blocks []model.Block
		if err := json.Unmarshal([]byte(contentJSON), &blocks); err != nil {
			p.Content = nil
		} else {
			p.Content = blocks
		}
		posts = append(posts, p)
	}

	return posts, nil
}

// Get liked posts by user
func (s *Service) Get_LikedPostsByUser(userID int) ([]model.Post, error) {
	return query.GetLikedPostsByUser(s.DB, userID)
}

