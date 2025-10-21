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
		http.Redirect(w, r, r.Header.Get("Referer")+"?error=Invalid+post+ID", http.StatusSeeOther)
		return
	}

	// Get the session ID from the session
	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/post/%d?error=Please+log+in+to+react", postID), http.StatusSeeOther)
		return
	}
	user, err := server.Service.GetUserFromSessionID(sessionID)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/post/%d?error=Session+expired", postID), http.StatusSeeOther)
		return
	}

	reactionType := r.FormValue("reaction_type")

	// Call the service layer to handle the reaction action
	err = server.Service.PostReaction(postID, user.ID, reactionType)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, fmt.Sprintf("/post/%d?error=Failed+to+update+reaction", postID), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/%d", postID), http.StatusSeeOther)
}
