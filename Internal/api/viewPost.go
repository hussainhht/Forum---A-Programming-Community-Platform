package api

import (
	"forum/Internal/model"
	"html/template"
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
		server.Service.HandleError(w, http.StatusNotFound)
		return
	}

	pageData := model.PageData{
		IsLoggedIn: isLoggedIn,
		Post:       post,
		User:       user,
	}
	// Render the post using a template
	tmpl, tmplErr := template.ParseFiles("./web/templates/view-post.html", "./web/templates/sidebar.html")
	if tmplErr != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	// fmt.Println(post).
	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
