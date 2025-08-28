package api

import (
	"html/template"
	"net/http"
)

func (server *Server) GetCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("./web/templates/create-post.html")
	if tmplErr != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Get categories to display in the form
	categories, err := server.Service.GetCategories()
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, categories)
}

func (server *Server) PostCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// check user authentication
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" || !server.Service.IsValidSession(cookie.Value) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["category"]

	// validate input
	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	if len(categories) == 0 {
		http.Error(w, "At least one category is required", http.StatusBadRequest)
		return
	}

	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		http.Error(w, "Error getting session ID", http.StatusInternalServerError)
		return
	}

	if err := server.Service.CreatePost(sessionID, title, content, categories); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Success response
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
