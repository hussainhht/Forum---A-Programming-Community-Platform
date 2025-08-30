package query

import (
	"database/sql"
	"forum/Internal/model"
)

var (
	InsertPostQuery         = "INSERT INTO posts (title, content, user_id, created_at) VALUES (?, ?, ?, datetime('now'))"
	InsertPostCategoryQuery = "INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)"
	GetAllPostsQuery        = `
		SELECT p.id, p.title, p.content, p.created_at, u.username, p.user_id
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`
	InsertPostReactionQuery = "INSERT INTO post_reactions (post_id, user_id, reaction_type) VALUES (?, ?, ?)"
	GetPostByIDQuery        = "SELECT id, title, content, created_at, user_id FROM posts WHERE id = ?"
)

func InsertPost(tx *sql.Tx, title, content string, categories []string, userID int) (int64, error) {

	// Step 1: Insert the post (without categories)
	result, err := tx.Exec(
		InsertPostQuery,
		title, content, userID,
	)
	if err != nil {
		return 0, err
	}

	// Step 2: Get the newly created post ID
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Step 3: Insert categories into the junction table
	for _, categoryID := range categories {
		_, err := tx.Exec(
			InsertPostCategoryQuery,
			postID, categoryID,
		)
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

// GetAllPosts retrieves all posts with their authors
func GetAllPosts(db *sql.DB) ([]model.Post, error) {
	rows, err := db.Query(GetAllPostsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Username,
			&post.UserID,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetPostByID(db *sql.DB, postID string) (*model.Post, error) {
	row := db.QueryRow(GetPostByIDQuery, postID)

	var post model.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &post, nil
}

// InsertPostReaction inserts a reaction into the post table
func InsertPostReaction(db *sql.DB, postID int, userID int, reactionType string) error {
	_, err := db.Exec(InsertPostReactionQuery, postID, userID, reactionType)
	return err
}

func GetPostReaction(db *sql.DB, postID int, userID int) (*model.PostReaction, error) {
	row := db.QueryRow("SELECT reaction_type FROM post_reactions WHERE post_id = ? AND user_id = ?", postID, userID)

	var reaction model.PostReaction
	reaction.UserID = userID
	reaction.PostID = postID
	if err := row.Scan(&reaction.ReactionType); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reaction, nil
}

func UpdatePostReaction(db *sql.DB, postID int, userID int, reactionType string) error {
	_, err := db.Exec("UPDATE post_reactions SET reaction_type = ? WHERE post_id = ? AND user_id = ?", reactionType, postID, userID)
	return err
}

func GetPostLikeCount(db *sql.DB, postID int) (int, error) {
	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND reaction_type = ?", postID, "like").Scan(&likeCount)
	if err != nil {
		return 0, err
	}
	return likeCount, nil
}