package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

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
