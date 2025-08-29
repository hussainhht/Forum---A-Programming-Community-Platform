package api

import (
	"net/http"
	"strconv"
)

func (server *Server) CommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	commentID, err := strconv.Atoi(r.FormValue("comment_id"))
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Get the session ID from the session
	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		http.Error(w, "Failed to get session ID", http.StatusUnauthorized)
		return
	}

	userID, err := server.Service.GetUserIDFromSessionID(sessionID)
	if err != nil {
		http.Error(w, "Failed to get user ID", http.StatusUnauthorized)
		return
	}

	reactionType := r.FormValue("reaction_type")

	// Call the service layer to handle the like action
	err = server.Service.CommentReaction(commentID, userID, reactionType)
	if err != nil {
		http.Error(w, "Failed to react to comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
