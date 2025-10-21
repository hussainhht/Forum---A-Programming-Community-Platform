package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func (server *Server) Get_LoginHandler(w http.ResponseWriter, r *http.Request) {
	//check if user is already logged in
	cookie, err := r.Cookie("session_id")
	if err == nil && server.Service.IsValidSession(cookie.Value) {
		// User is already logged in
		server.Service.LogoutUser(cookie.Value)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, tmplErr := template.ParseFiles("./web/templates/login.html")
	if tmplErr != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func (server *Server) Post_LoginHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		renderLogin(w, r, "Failed to parse form")
		return
	}

	emailOrUsername := strings.TrimSpace(r.FormValue("emailORUsername"))
	password := r.FormValue("password")

	newSessionID, err := server.Service.LoginUser(emailOrUsername, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		renderLogin(w, r, "Invalid email/username or password")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour), // 1 day
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func renderLogin(w http.ResponseWriter, r *http.Request, errMsg string) {
	tmpl, _ := template.ParseFiles("./web/templates/login.html")
	data := model.LoginPageData{Error: errMsg}
	nameorEmail := r.FormValue("emailORUsername")
	fmt.Println(nameorEmail)
	data.Form.EmailOrUsername = nameorEmail
	_ = tmpl.Execute(w, data)
}
