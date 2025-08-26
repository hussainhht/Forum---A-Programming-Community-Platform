package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request , db *sql.DB) {

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "./temblats/login.html")
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

		if err != nil || password != storedPassword {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "✅ Login successful! Welcome %s", email)

	}

}
