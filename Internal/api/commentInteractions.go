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
		http.Redirect(w, r, r.Header.Get("Referer")+"?error=Invalid+comment+ID", http.StatusSeeOther)
		return
	}

	postID := r.FormValue("post_id")
	if postID == "" {
		http.Redirect(w, r, r.Header.Get("Referer")+"?error=Missing+post+ID", http.StatusSeeOther)
		return
	}

	// Get the session ID from the session
	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/post/%s?error=Please+log+in+to+react", postID), http.StatusSeeOther)
		return
	}

	user, err := server.Service.GetUserFromSessionID(sessionID)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/post/%s?error=Session+expired", postID), http.StatusSeeOther)
		return
	}

	reactionType := r.FormValue("reaction_type")

	// Call the service layer to handle the like action
	err = server.Service.CommentReaction(commentID, user.ID, reactionType)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/post/%s?error=Failed+to+update+reaction", postID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%s", postID), http.StatusSeeOther)
}

func (server *Server) Post_CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postID := r.FormValue("post_id")
	if postID == "" {
		http.Redirect(w, r, r.Header.Get("Referer")+"?error=Missing+post+ID", http.StatusSeeOther)
		return
	}

	// Get the session ID from the session
	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/post/%s?error=Please+log+in+to+comment", postID), http.StatusSeeOther)
		return
	}

	user, err := server.Service.GetUserFromSessionID(sessionID)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/post/%s?error=Session+expired", postID), http.StatusSeeOther)
		return
	}

	content := r.FormValue("content")
	if len(content) > 1000 {
		http.Redirect(w, r, fmt.Sprintf("/post/%s?error=Comment+too+long+(max+1000+characters)", postID), http.StatusSeeOther)
		return
	}

	// Call the service layer to create the comment
	err = server.Service.CreateComment(postID, user.ID, content)
	if err != nil {
		fmt.Println("Error creating comment:", err)
		http.Redirect(w, r, fmt.Sprintf("/post/%s?error=Failed+to+create+comment", postID), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%s?success=Comment+added+successfully", postID), http.StatusSeeOther)
}
