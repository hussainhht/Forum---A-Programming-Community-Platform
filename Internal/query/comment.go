package query

import (
	"database/sql"
	"forum/Internal/model"
)

var (
	InsertCommentReactionQuery = "INSERT INTO comments_reactions (comment_id, user_id, reaction_type) VALUES (?, ?, ?)"
)

func InsertComment(db *sql.DB, postID int, userID int, content string) error {
	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	return err
}

func InsertCommentReaction(db *sql.DB, commentID int, userID int, reactionType string) error {
	_, err := db.Exec(InsertCommentReactionQuery, commentID, userID, reactionType)
	return err
}

func GetCommentReaction(db *sql.DB, commentID int, userID int) (*model.CommentReaction, error) {
	row := db.QueryRow("SELECT reaction_type FROM comments_reactions WHERE comment_id = ? AND user_id = ?", commentID, userID)

	var reaction model.CommentReaction
	reaction.UserID = userID
	reaction.CommentID = commentID
	if err := row.Scan(&reaction.ReactionType); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &reaction, nil
}

func UpdateCommentReaction(db *sql.DB, commentID int, userID int, reactionType string) error {
	_, err := db.Exec("UPDATE comments_reactions SET reaction_type = ? WHERE comment_id = ? AND user_id = ?", reactionType, commentID, userID)
	return err
}

func DeleteCommentReaction(db *sql.DB, commentID int, userID int) error {
	_, err := db.Exec("DELETE FROM comments_reactions WHERE comment_id = ? AND user_id = ?", commentID, userID)
	return err
}

func GetCommentLikeCount(db *sql.DB, commentID int) (int, error) {
	var likeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM comments_reactions WHERE comment_id = ? AND reaction_type = ?", commentID, "like").Scan(&likeCount)
	if err != nil {
		return 0, err
	}
	return likeCount, nil
}

func GetCommentDislikeCount(db *sql.DB, commentID int) (int, error) {
	var dislikeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM comments_reactions WHERE comment_id = ? AND reaction_type = ?", commentID, "dislike").Scan(&dislikeCount)
	if err != nil {
		return 0, err
	}
	return dislikeCount, nil
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
