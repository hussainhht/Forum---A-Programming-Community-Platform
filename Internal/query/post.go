package query

import "database/sql"

var (
	InsertPostQuery         = "INSERT INTO posts (title, content, user_id, created_at) VALUES (?, ?, ?, datetime('now'))"
	InsertPostCategoryQuery = "INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)"
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
