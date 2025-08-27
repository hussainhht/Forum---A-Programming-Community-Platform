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
	tmpl , err:= template.ParseFiles("templates/home.html")
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

func RegisterHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./templates/register.html")
		return
	}

	if r.Method == http.MethodPost {
		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get form values
		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmPassword")

		// Validate input
		if username == "" || email == "" || password == "" || confirmPassword == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Check username length
		if len(username) < 3 || len(username) > 20 {
			http.Error(w, "Username must be between 3 and 20 characters", http.StatusBadRequest)
			return
		}

		// Validate email format
		if !IsValidEmail(email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		// Validate password
		if !IsValidPassword(password) {
			http.Error(w, "Password must be at least 8 characters and contain uppercase, lowercase, digit, and special character", http.StatusBadRequest)
			return
		}

		// Check if passwords match
		if password != confirmPassword {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// Check if username already exists
		var existingUsername string
		err := db.QueryRow("SELECT username FROM users WHERE username = ? LIMIT 1", username).Scan(&existingUsername)
		if err == nil {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		} else if err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Check if email already exists
		var existingEmail string
		err = db.QueryRow("SELECT email FROM users WHERE email = ? LIMIT 1", email).Scan(&existingEmail)
		if err == nil {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		} else if err != sql.ErrNoRows {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		// Insert new user into database
		_, err = db.Exec(
			"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			username, email, password,
		)
		if err != nil {
			http.Error(w, "Failed to create account", http.StatusInternalServerError)
			return
		}

		// Success response
		fmt.Fprintf(w, "✅ Registration successful! Welcome %s. You can now <a href='/login'>login here</a>.", username)
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	template, err := template.ParseFiles("./templates/createpost.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		template.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		// Handle form submission
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		// categories := r.Form["category"]
		createdAt := time.Now()

		// Validate input
		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Insert post into database
		_, err := db.Exec(
			"INSERT INTO posts (title, content, created_at) VALUES (?, ?, ?)",
			title, content, createdAt,
		)
		if err != nil {
			http.Error(w, "Failed to create post", http.StatusInternalServerError)
			return
		}
		fmt.Println("Post created successfully")
		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
