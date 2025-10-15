package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"net/http"
)

func (server *Server) Get_HelpHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the templates
	tmpl, tmplErr := template.ParseFiles("./web/templates/root.html", "./web/templates/help.html")

	if tmplErr != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		fmt.Println("Error creating template:", tmplErr)
		return
	}

	// Check if user is logged in
	sessionIDCookie, _ := r.Cookie("session_id")
	isLoggedIn := false
	if sessionIDCookie != nil {
		isLoggedIn = server.Service.IsValidSession(sessionIDCookie.Value)
	}

	// Get user info if logged in
	var user *model.User
	var err error
	if isLoggedIn {
		user, err = server.Service.GetUserFromSessionID(sessionIDCookie.Value)
		if err != nil {
			server.Service.HandleError(w, r, http.StatusInternalServerError)
			return
		}
	}

	// Get categories for sidebar
	categories, err := server.Service.GetCategories()
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	pageData := model.PageData{
		IsLoggedIn: isLoggedIn,
		User:       user,
		Categories: categories,
		CSSFile:    "/web/static/css/help.css",
		ExtraCSS:   nil,
		Theme:      nil,
	}

	// Execute the template
	tmpl.Execute(w, pageData)
}
