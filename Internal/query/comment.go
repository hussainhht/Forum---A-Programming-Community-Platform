package query

import (
	"database/sql"
	"forum/Internal/model"
)

var (
	InsertCommentReactionQuery = "INSERT INTO comment_reactions (comment_id, user_id, reaction_type) VALUES (?, ?, ?)"
)

func InsertComment(db *sql.DB, postID int, userID int, content string) error {
	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	return err
}

func InsertCommentReaction(db *sql.DB, commentID int, userID int, reactionType string) error {
	_, err := db.Exec(InsertCommentReactionQuery, commentID, userID, reactionType)
	return err
}

func GetCommentsByPostID(db *sql.DB, postID int) ([]model.Comment, error) {
	rows, err := db.Query("SELECT id, content, user_id, post_id FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
