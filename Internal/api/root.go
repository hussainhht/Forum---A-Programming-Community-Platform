package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"net/http"
)

func (server *Server) Get_HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" && r.URL.Path != "" {
		fmt.Println("Path not found:", r.URL.Path)
		server.Service.HandleError(w, r, http.StatusNotFound)
		return
	}

	sessionIDCookie, _ := r.Cookie("session_id")

	// Create template with custom functions
	tmpl, tmplErr := template.ParseFiles("./web/templates/root.html", "./web/templates/home.html")

	if tmplErr != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		fmt.Println("Error creating template:", tmplErr)
		return
	}

	// Get categories for the filter
	categories, err := server.Service.GetCategories()
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		fmt.Println("Error getting categories:", err)
		return
	}

	// Get posts - either filtered by categories or all posts
	var featuredPosts []model.Post

	featuredPosts, err = server.Service.GetFeaturedPosts()

	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		fmt.Println("Error getting featured posts:", err)
		return
	}
	var LatestPosts []model.Post

	LatestPosts, err = server.Service.GetLatestPosts()
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		fmt.Println("Error getting latest posts:", err)
		return
	}

	//if the user logged in this will be true
	isLoggedIn := false
	if sessionIDCookie != nil {
		isLoggedIn = server.Service.IsValidSession(sessionIDCookie.Value)
	}
	// Get user info if logged in
	var user *model.User
	if isLoggedIn {
		user, err = server.Service.GetUserFromSessionID(sessionIDCookie.Value)
		if err != nil {
			server.Service.HandleError(w, r, http.StatusInternalServerError)
			fmt.Println("Error getting user from session ID:", err)
			return
		}
	}

	pageData := model.PageData{
		IsLoggedIn:    isLoggedIn,
		User:          user,
		FeaturedPosts: featuredPosts,
		LatestPosts:   LatestPosts,
		Categories:    categories,
		CSSFile:       "/web/static/css/home.css",
		ExtraCSS:      nil,
		Theme:         nil, // home uses fallback variables defined in category-base
	}

	// Pass posts to the template
	tmpl.Execute(w, pageData)
}
