package api

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

func (server *Server) Get_LoginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, tmplErr := template.ParseFiles("./web/templates/login.html")
	if tmplErr != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (server *Server) Post_LoginHandler(w http.ResponseWriter, r *http.Request) {
	emailOrUsername := strings.TrimSpace(strings.ToLower(r.FormValue("emailORUsername")))
	password := r.FormValue("password")

	newSessionID, err := server.Service.LoginUser(emailOrUsername, password)
	if err != nil {
		//! remove this log
		log.Printf("\nLogin failed for user %q: %v", emailOrUsername, err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400, // 1 day
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
