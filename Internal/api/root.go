package api

import (
	"html/template"
	"net/http"
)

func (server *Server) GetRootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("./web/templates/home.html")
	if tmplErr != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Get all posts from the service
	posts, err := server.Service.GetAllPosts()
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	// Pass posts to the template
	tmpl.Execute(w, posts)
}
