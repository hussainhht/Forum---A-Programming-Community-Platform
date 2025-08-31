package api

import (
	"net/http"
	query "forum/Internal/query"
)

func (server *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get session_id cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Remove session from DB
	sessionID := cookie.Value
	userID, _ := server.Service.GetUserIDFromSessionID(sessionID)
	if userID != 0 {
		_ = query.RemoveSession(server.Service.DB, userID)
	}

	// Delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // expire now
	})

	// Redirect to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
