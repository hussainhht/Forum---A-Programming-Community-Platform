package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"net/http"
)

// Common handler function for both Help and Privacy pages
func (server *Server) handleHelpPrivacy(w http.ResponseWriter, r *http.Request, showPrivacy bool) {
	// Parse the templates
	tmpl, tmplErr := template.ParseFiles("./web/templates/root.html", "./web/templates/help-privacy.html")

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

	pageData := struct {
		model.PageData
		ShowPrivacy bool
	}{
		PageData: model.PageData{
			IsLoggedIn: isLoggedIn,
			User:       user,
			Categories: categories,
			CSSFile:    "/web/static/css/help-privacy.css",
			ExtraCSS:   nil,
			Theme:      nil,
		},
		ShowPrivacy: showPrivacy,
	}

	// Execute the template
	tmpl.Execute(w, pageData)
}

func (server *Server) Get_HelpHandler(w http.ResponseWriter, r *http.Request) {
	server.handleHelpPrivacy(w, r, false) // Show help panel by default
}

func (server *Server) Get_PrivacyTermsHandler(w http.ResponseWriter, r *http.Request) {
	server.handleHelpPrivacy(w, r, true) // Show privacy panel by default
}
