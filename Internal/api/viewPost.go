package api

import (
	"forum/Internal/model"
	"html/template"
	"net/http"
)

func (server *Server) Get_PostHandler(w http.ResponseWriter, r *http.Request) {
	// Handler logic for getting a single post
	postID := r.PathValue("id") // Extract ID from URL pattern /post/{id}

	post, err := server.Service.GetPostByID(postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	pageData := model.PageData{
		IsLoggedIn: true, // Replace with actual login check
		Post:       post,
	}
	// Render the post using a template
	tmpl, tmplErr := template.ParseFiles("./web/templates/view-post.html")
	if tmplErr != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}
	// fmt.Println(post)
	if err := tmpl.Execute(w, pageData); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}
