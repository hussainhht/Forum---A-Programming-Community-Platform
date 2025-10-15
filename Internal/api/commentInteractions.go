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
		server.Service.HandleError(w, r, http.StatusBadRequest)
		return
	}

	// Get the session ID from the session
	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusUnauthorized)
		return
	}

	user, err := server.Service.GetUserFromSessionID(sessionID)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusUnauthorized)
		return
	}

	reactionType := r.FormValue("reaction_type")

	// Call the service layer to handle the like action
	err = server.Service.CommentReaction(commentID, user.ID, reactionType)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (server *Server) Post_CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postID := r.FormValue("post_id")
	if postID == "" {
		server.Service.HandleError(w, r, http.StatusBadRequest)
		return
	}

	// Get the session ID from the session
	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusUnauthorized)
		return
	}

	user, err := server.Service.GetUserFromSessionID(sessionID)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusUnauthorized)
		return
	}

	content := r.FormValue("content")
	if len(content) > 1000 {
		server.Service.HandleError(w, r, http.StatusBadRequest)
		return
	}

	// Call the service layer to create the comment
	err = server.Service.CreateComment(postID, user.ID, content)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s", postID), http.StatusSeeOther)
	w.WriteHeader(http.StatusCreated)
}
