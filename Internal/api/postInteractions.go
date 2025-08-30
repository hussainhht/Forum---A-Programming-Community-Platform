package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func (server *Server) PostReactionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
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
	err = server.Service.PostReaction(postID, userID, reactionType)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to react to post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
