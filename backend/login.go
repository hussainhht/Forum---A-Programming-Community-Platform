package backend

import (
	"database/sql"
	"forum/database"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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

		var (
			userID         int
			storedPassword string
		)
		err := db.QueryRow(
			"SELECT id, password FROM users WHERE lower(email) = lower(?) LIMIT 1",
			email,
		).Scan(&userID, &storedPassword)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// compare plain password with stored hash
		if bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		// Delete any existing sessions for this user
<<<<<<< HEAD
   	 	_, _ = db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
		// Remove all records from the table called sessions where the column user_id matches a certain value
		// table called sessions // column user_id // ? is a placeholder for the value userID
=======
		_, _ = db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
>>>>>>> 1902c390f92b8c688644bdadce8f80c01009f05a

		sessionID := GenerateSessionID()
		// Pass an absolute expiration time instead of a duration
		if err = database.CreateSession(db, userID, sessionID, time.Now().Add(24*time.Hour)); err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   86400, // 1 day
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}

func GenerateSessionID() string {
	return uuid.NewString() // returns a string like: "550e8400-e29b-41d4-a716-446655440000"
}
