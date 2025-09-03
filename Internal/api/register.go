package api

import (
	"html/template"
	"net/http"
	"strings"
)

type RegisterPageData struct {
	Error string
	Form  struct {
		Username string
		Email    string
	}
}

func (server *Server) Get_RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("./web/templates/register.html")
	if tmplErr != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (server *Server) Post_RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		renderRegister(w, "Failed to parse form", r)
		return
	}

	// Get form values
	username := strings.TrimSpace(r.FormValue("username"))
	if len(username) > 50 {
		renderRegister(w, "Username is too long,(the maximum length is 50 characters)", r)
		return
	}

	if strings.Contains(username, " ") {
		renderRegister(w, "Username must not contain spaces", r)
		return
	}

	email := strings.TrimSpace(strings.ToLower(r.FormValue("email")))
	if len(email) > 100 {
		renderRegister(w, "Email is too long,(the maximum length is 50 characters)", r)
		return
	}
	password := r.FormValue("password")
	if len(password) > 80 {
		renderRegister(w, "the password must be max 80 characters long", r)
		return
	}
	confirmPassword := r.FormValue("confirmPassword")

	if err := server.Service.RegisterUser(username, email, password, confirmPassword); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		renderRegister(w, err.Error(), r)
		return
	}
	// Success response
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func renderRegister(w http.ResponseWriter, errMsg string, r *http.Request) {

	tmpl, _ := template.ParseFiles("./web/templates/register.html")

	data := RegisterPageData{
		Error: errMsg,
	}
	
	data.Form.Username = r.FormValue("username")
	data.Form.Email = r.FormValue("email")
	_ = tmpl.Execute(w, data)
}
