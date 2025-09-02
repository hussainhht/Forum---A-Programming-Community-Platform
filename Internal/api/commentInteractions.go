package api

import (
	"fmt"
	"net/http"
	"strconv"
)


func (server *Server) CommentReactionHandler(w http.ResponseWriter, r *http.Request) {

	// Check if user is authorized


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
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

func (server *Server) Post_CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postID := r.FormValue("post_id")
	if postID == "" {
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
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")
	if len(content) > 1000 {
		http.Error(w, "Comment content is too long", http.StatusBadRequest)
		return
	}

	// Call the service layer to create the comment
	err = server.Service.CreateComment(postID, userID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postID), http.StatusSeeOther)
	w.WriteHeader(http.StatusCreated)
}

