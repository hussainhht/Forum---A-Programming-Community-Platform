package query

import (
	"database/sql"
)

var (
	InsertCommentReactionQuery = "INSERT INTO comment_reactions (comment_id, user_id, reaction_type) VALUES (?, ?, ?)"
)

func InsertCommentReaction(db *sql.DB, commentID int, userID int, reactionType string) error {
	_, err := db.Exec(InsertCommentReactionQuery, commentID, userID, reactionType)
	return err
}
