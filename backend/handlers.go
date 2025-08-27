package backend

import (
	"database/sql"
	"fmt"
	"forum/database"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	rows, err := database.DB.Query(`
            SELECT posts.id, posts.title, posts.content, users.username
            FROM posts
            JOIN users ON posts.user_id = users.id
            ORDER BY posts.id DESC
        `)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			http.Error(w, "Error reading posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}


func CreatePostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	template, err := template.ParseFiles("./templates/createpost.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		// Get categories to display in the form
		rows, err := db.Query("SELECT id, name FROM categories")
		if err != nil {
			http.Error(w, "Error loading categories", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		categories := []struct {
			ID   int
			Name string
		}{}

		for rows.Next() {
			var cat struct {
				ID   int
				Name string
			}
			if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
				http.Error(w, "Error reading categories", http.StatusInternalServerError)
				return
			}
			categories = append(categories, cat)
		}

		template.Execute(w, categories)
		return
	}

	if r.Method == http.MethodPost {
		// Handle form submission
		// get session_id from cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		sessionID := cookie.Value

		// get user from sessionid in db
		var userID int
		row := db.QueryRow("SELECT user_id FROM sessions WHERE session_id = ?", sessionID).Scan(&userID)
		if row == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// get post data
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		categories := r.Form["category"]
		createdAt := time.Now()

		// Validate input
		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		if len(categories) == 0 {
			http.Error(w, "At least one category is required", http.StatusBadRequest)
			return
		}

		// Start a transaction for consistent data
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Insert post into database
		result, err := tx.Exec(
			"INSERT INTO posts (title, content, created_at, user_id) VALUES (?, ?, ?, ?)",
			title, content, createdAt, userID,
		)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the ID of the newly created post
		postID, err := result.LastInsertId()
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to get post ID", http.StatusInternalServerError)
			return
		}

		// Insert categories for the post
		for _, categoryID := range categories {
			_, err = tx.Exec(
				"INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)",
				postID, categoryID,
			)
			if err != nil {
				tx.Rollback()
				http.Error(w, "Failed to add category to post: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Commit the transaction
		if err = tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		fmt.Println("Post created successfully with", len(categories), "categories")
		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
