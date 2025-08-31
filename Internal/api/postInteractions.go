package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func (server *Server) Post_ReactionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Get the session ID from the session
	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}
	userID, err := server.Service.GetUserIDFromSessionID(sessionID)
	if err != nil {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}

	reactionType := r.FormValue("reaction_type")

	// Call the service layer to handle the reaction action
	err = server.Service.PostReaction(postID, userID, reactionType)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to react to post", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
	w.WriteHeader(http.StatusOK)
}
