package api

import (
	"forum/Internal/model"
	"html/template"
	"net/http"
)

func (server *Server) Get_RootHandler(w http.ResponseWriter, r *http.Request) {

	sessionIDCookie, _ := r.Cookie("session_id")
	

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

	//if the user logged in this will be true
	isLoggedIn := false
	if sessionIDCookie != nil {
		isLoggedIn = server.Service.IsValidSession(sessionIDCookie.Value)
	}

	pageData := model.PageData{
		IsLoggedIn: isLoggedIn,
		Posts:      posts,
	}

	// Pass posts to the template
	tmpl.Execute(w, pageData)
}