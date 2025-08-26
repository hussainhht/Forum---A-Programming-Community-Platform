package handlers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	template, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		template.Execute(w, nil)
		return
	}

	if r.Method == http.MethodPost {
		email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
		password := r.FormValue("password")

		var storedPassword string
		err := db.QueryRow(
			"SELECT password FROM users WHERE lower(email) = lower(?) LIMIT 1",
			email,
		).Scan(&storedPassword)

		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// compare plain password with the hash
		compareErr := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if compareErr != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "✅ Login successful! Welcome %s", email)

	}

}
