package api

import (
	"html/template"
	"net/http"
	"strings"
)

func (server *Server) Get_RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("./web/templates/register.html")
	if tmplErr != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (server *Server) Post_RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if err := server.Service.RegisterUser(username, email, password, confirmPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Success response
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
