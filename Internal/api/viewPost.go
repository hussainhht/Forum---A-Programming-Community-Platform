package api

import (
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
)

func (server *Server) Get_PostHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for getting a single post

	postID := r.PathValue("id") // Extract ID from URL pattern /post/{id}

	isLoggedIn := false
	var user *model.User

	if cookie, err := r.Cookie("session_id"); err == nil {
		if server.Service.IsValidSession(cookie.Value) {
			isLoggedIn = true
			user, _ = server.Service.GetUserFromSessionID(cookie.Value)
		}
	}

	post, err := server.Service.GetPostByID(postID)
	if err != nil {
		log.Printf("viewPost: failed to get post id=%s: %v", postID, err)
		server.Service.HandleError(w, r, http.StatusNotFound)
		return
	}

	// Set user reaction flags if user is logged in
	if user != nil {
		reaction, err := server.Service.GetPostReaction(post.ID, user.ID)
		if err == nil && reaction != nil {
			post.UserLiked = reaction.ReactionType == "like"
			post.UserDisliked = reaction.ReactionType == "dislike"
		}
	}

	pageData := model.PageData{
		IsLoggedIn: isLoggedIn,
		Post:       post,
		User:       user,
	}
	// Render the post using a template
	// Parse with sidebar (if sidebar defines template name 'sidebar')
	tmpl, tmplErr := template.ParseFiles("./web/templates/root.html", "./web/templates/view-post.html")
	if tmplErr != nil {
		log.Printf("viewPost: template parse error: %v", tmplErr)
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
	// fmt.Println(post).
	if err := tmpl.Execute(w, pageData); err != nil {
		log.Printf("viewPost: execute error: %v", err)
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
}
