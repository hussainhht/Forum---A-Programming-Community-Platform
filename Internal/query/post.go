package query

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/Internal/model"
	"strings"
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
	GetPostsByCategoriesQuery = `
		SELECT DISTINCT p.id, p.title, p.content, p.created_at, u.username, p.user_id
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		WHERE pc.category_id IN (`
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
		var (
			post        model.Post
			contentJSON string
		)
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&contentJSON,
			&post.CreatedAt,
			&post.Username,
			&post.UserID,
		)
		if err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(contentJSON), &post.Content)
		// derive preview: first text block
		for _, b := range post.Content {
			if b.Type == "text" && post.Preview == "" {
				post.Preview = b.Content
				break
			}
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostsByCategories retrieves posts filtered by selected category IDs
func GetPostsByCategories(db *sql.DB, categoryIDs []string) ([]model.Post, error) {
	if len(categoryIDs) == 0 {
		return GetAllPosts(db)
	}

	// Build the query with placeholders
	placeholders := ""
	for i := range categoryIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	query := GetPostsByCategoriesQuery + placeholders + ") ORDER BY p.created_at DESC"

	// Convert string slice to interface slice for query args
	args := make([]interface{}, len(categoryIDs))
	for i, id := range categoryIDs {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var (
			post        model.Post
			contentJSON string
		)
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&contentJSON,
			&post.CreatedAt,
			&post.Username,
			&post.UserID,
		)
		if err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(contentJSON), &post.Content)
		for _, b := range post.Content {
			if b.Type == "text" && post.Preview == "" {
				post.Preview = b.Content
				break
			}
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func SelectCategoryByID(db *sql.DB, categoryID int) (*model.Category, error) {
	row := db.QueryRow("SELECT id, name FROM categories WHERE id = ?", categoryID)

	var category model.Category
	if err := row.Scan(&category.ID, &category.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func GetPostByID(db *sql.DB, postID string) (*model.Post, error) {
	row := db.QueryRow(GetPostByIDQuery, postID)

	var (
		post        model.Post
		contentJSON string
	)
	if err := row.Scan(&post.ID, &post.Title, &contentJSON, &post.CreatedAt, &post.UserID); err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, nil
		// }
		return nil, err
	}

	_ = json.Unmarshal([]byte(contentJSON), &post.Content)

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

func DeletePostReaction(db *sql.DB, postID int, userID int) error {
	_, err := db.Exec("DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?", postID, userID)
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

func GetPostDislikeCount(db *sql.DB, postID int) (int, error) {
	var dislikeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM post_reactions WHERE post_id = ? AND reaction_type = ?", postID, "dislike").Scan(&dislikeCount)
	if err != nil {
		return 0, err
	}
	return dislikeCount, nil
}

func GetFeaturedPosts(db *sql.DB) ([]model.Post, error) {
	rows, err := db.Query(`
WITH engagement AS (
    SELECT
        p.id AS post_id,
        COUNT(DISTINCT CASE WHEN pr.reaction_type = 'like' THEN pr.user_id END) AS likes,
        COUNT(DISTINCT c.id) AS comments
    FROM posts p
    LEFT JOIN post_reactions pr ON p.id = pr.post_id
    LEFT JOIN comments c ON p.id = c.post_id
    GROUP BY p.id
)
SELECT
    p.id,
    p.title,
    p.content,
    p.created_at,
    u.username,
	p.user_id,
    e.likes,
    e.comments,
    COALESCE(
		(
	    	(e.likes * 3 + e.comments * 5) * 1.0 /
        	(strftime('%s','now') - strftime('%s', p.created_at)) / 3600.0 + 2 
		), 0
    ) AS featured_score
FROM posts p
JOIN users u ON p.user_id = u.id
JOIN engagement e ON p.id = e.post_id
ORDER BY featured_score DESC
LIMIT 4;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var contentJSON string
		if err := rows.Scan(&post.ID, &post.Title, &contentJSON, &post.CreatedAt, &post.Username, &post.UserID, &post.LikeCount, &post.CommentCount, &post.FeaturedScore); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(contentJSON), &post.Content)
		post.Preview = getPostPreview(post.Content)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func GetLikedPostsByUser(db *sql.DB, userID int) ([]model.Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.created_at,
	       (SELECT COUNT(*) FROM post_reactions WHERE post_id = p.id AND reaction_type='like') as like_count,
	       (SELECT COUNT(*) FROM comments WHERE post_id = p.id) as comment_count
	FROM posts p
	JOIN post_reactions pr ON p.id = pr.post_id
	WHERE pr.user_id = ? AND pr.reaction_type = 'like'
	ORDER BY p.created_at DESC
	`
	rows, err := db.Query(query, userID)
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

func GetLatestPosts(db *sql.DB) ([]model.Post, error) {
	rows, err := db.Query(`
	SELECT p.id, p.title, p.content, p.created_at, u.username, p.user_id
	FROM posts p
	JOIN users u ON p.user_id = u.id
	ORDER BY p.created_at DESC
	LIMIT 4;
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []model.Post
	for rows.Next() {
		var post model.Post
		var contentJSON string
		if err := rows.Scan(&post.ID, &post.Title, &contentJSON, &post.CreatedAt, &post.Username, &post.UserID); err != nil {
			return nil, err
		}

		_ = json.Unmarshal([]byte(contentJSON), &post.Content)
		post.Preview = getPostPreview(post.Content)
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// GetDiscoverPosts provides filtered, searched, sorted, paginated posts
func GetDiscoverPosts(db *sql.DB, search, category, sort string, limit, offset int) ([]model.Post, bool, error) {
	// base select - now includes category name
	base := `SELECT DISTINCT p.id, p.title, p.content, p.created_at, u.username, p.user_id,
		(SELECT COUNT(*) FROM post_reactions pr WHERE pr.post_id = p.id AND pr.reaction_type='like') as like_count,
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comment_count,
		COALESCE(c.name, '') as category_name
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_categories pc ON p.id = pc.post_id
		LEFT JOIN categories c ON pc.category_id = c.id`
	var where []string
	var args []interface{}

	// Only apply search filter if search term is not empty
	if search != "" && strings.TrimSpace(search) != "" {
		where = append(where, "(p.title LIKE ? OR p.content LIKE ?)")
		like := fmt.Sprintf("%%%s%%", search)
		args = append(args, like, like)
	}

	// Only apply category filter if category is not empty
	if category != "" && strings.TrimSpace(category) != "" {
		where = append(where, "c.name = ?")
		args = append(args, category)
	}

	if len(where) > 0 {
		base += " WHERE " + strings.Join(where, " AND ")
	}
	// sorting
	switch sort {
	case "popular":
		base += " ORDER BY like_count DESC, p.created_at DESC"
	case "comments":
		base += " ORDER BY comment_count DESC, p.created_at DESC"
	default:
		base += " ORDER BY p.created_at DESC"
	}
	// pagination fetch limit+1 to detect next page
	base += " LIMIT ? OFFSET ?"
	args = append(args, limit+1, offset)
	rows, err := db.Query(base, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	posts := []model.Post{}
	for rows.Next() {
		var post model.Post
		var contentJSON string
		if err := rows.Scan(&post.ID, &post.Title, &contentJSON, &post.CreatedAt, &post.Username, &post.UserID, &post.LikeCount, &post.CommentCount, &post.Category); err != nil {
			return nil, false, err
		}
		_ = json.Unmarshal([]byte(contentJSON), &post.Content)
		// derive preview
		for _, b := range post.Content {
			if b.Type == "text" {
				post.Preview = b.Content
				break
			}
		}
		posts = append(posts, post)
	}
	hasNext := false
	if len(posts) > limit {
		hasNext = true
		posts = posts[:limit]
	}
	return posts, hasNext, rows.Err()
}

// GetDiscoverPostsMultiCategory provides filtered, searched, sorted, paginated posts with multiple category support
func GetDiscoverPostsMultiCategory(db *sql.DB, search string, categories []string, sort string, limit, offset int) ([]model.Post, bool, error) {
	// Base query without category filter
	baseQuery := `SELECT DISTINCT p.id, p.title, p.content, p.created_at, u.username, p.user_id,
		(SELECT COUNT(*) FROM post_reactions pr WHERE pr.post_id = p.id AND pr.reaction_type='like') as like_count,
		(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id`

	var where []string
	var args []interface{}

	// Add category filter if categories are provided
	if len(categories) > 0 {
		// Filter out empty strings
		validCategories := []string{}
		for _, cat := range categories {
			if strings.TrimSpace(cat) != "" {
				validCategories = append(validCategories, cat)
			}
		}

		if len(validCategories) > 0 {
			baseQuery += `
			JOIN post_categories pc ON p.id = pc.post_id
			JOIN categories c ON pc.category_id = c.id`

			// Build OR condition for multiple categories
			catPlaceholders := make([]string, len(validCategories))
			for i := range validCategories {
				catPlaceholders[i] = "c.name = ?"
				args = append(args, validCategories[i])
			}
			where = append(where, "("+strings.Join(catPlaceholders, " OR ")+")")
		}
	}

	// Apply search filter if search term is not empty
	if search != "" && strings.TrimSpace(search) != "" {
		where = append(where, "(p.title LIKE ? OR p.content LIKE ?)")
		like := fmt.Sprintf("%%%s%%", search)
		args = append(args, like, like)
	}

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	// Sorting
	switch sort {
	case "popular":
		baseQuery += " ORDER BY like_count DESC, p.created_at DESC"
	case "comments":
		baseQuery += " ORDER BY comment_count DESC, p.created_at DESC"
	default:
		baseQuery += " ORDER BY p.created_at DESC"
	}

	// Pagination: fetch limit+1 to detect next page
	baseQuery += " LIMIT ? OFFSET ?"
	args = append(args, limit+1, offset)

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var post model.Post
		var contentJSON string
		if err := rows.Scan(&post.ID, &post.Title, &contentJSON, &post.CreatedAt, &post.Username, &post.UserID, &post.LikeCount, &post.CommentCount); err != nil {
			return nil, false, err
		}
		_ = json.Unmarshal([]byte(contentJSON), &post.Content)
		// Derive preview
		for _, b := range post.Content {
			if b.Type == "text" {
				post.Preview = b.Content
				break
			}
		}
		posts = append(posts, post)
	}

	hasNext := false
	if len(posts) > limit {
		hasNext = true
		posts = posts[:limit]
	}
	return posts, hasNext, rows.Err()
}

func getPostPreview(content []model.Block) string {
	var preview string
	wordCount := 0
	for _, block := range content {
		wordCount += len(strings.Fields(block.Content))
		preview += block.Content + " "
		if wordCount >= 30 {
			break
		}
	}
	return preview
}
