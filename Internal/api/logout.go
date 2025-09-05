package api

import (
	"net/http"
)

func (server *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get session_id cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sessionID := cookie.Value

	// Clear temp blocks for this session
	delete(server.TempBlocks, sessionID)

	// Remove session from DB
	server.Service.LogoutUser(sessionID)

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
