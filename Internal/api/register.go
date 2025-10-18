package api

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type RegisterPageData struct {
	Error        string
	ShowRegister bool // Flag to show register form by default
	Form         struct {
		Username string
		Email    string
		FullName string
	}
}

func (server *Server) Get_RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// When accessing /register, show the register form (checkbox checked)
	tmpl, tmplErr := template.ParseFiles("./web/templates/login.html")
	if tmplErr != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	data := RegisterPageData{
		ShowRegister: true, // This will be used to check the checkbox
	}
	tmpl.Execute(w, data)
}

func (server *Server) Post_RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form for potential file upload
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
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

	fullname := strings.TrimSpace(r.FormValue("fullname"))
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

	// Handle optional image upload
	photoFilename := ""
	file, header, errFile := r.FormFile("profile-img")
	if errFile == nil && header != nil {
		defer file.Close()
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" {
			photoFilename = fmt.Sprintf("user_%s%s", username, ext)
			outPath := filepath.Join("web", "static", "img", photoFilename)
			out, err := os.Create(outPath)
			if err == nil {
				defer out.Close()
				_, _ = io.Copy(out, file)
			} else {
				photoFilename = "" // fallback to default later
			}
		}
	}

	if err := server.Service.RegisterUser(username, email, password, confirmPassword, fullname, photoFilename); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		renderRegister(w, err.Error(), r)
		return
	}
	// Success response
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func renderRegister(w http.ResponseWriter, errMsg string, r *http.Request) {

	tmpl, _ := template.ParseFiles("./web/templates/login.html")

	data := RegisterPageData{
		Error:        errMsg,
		ShowRegister: true, // Keep register form visible on error
	}

	data.Form.Username = r.FormValue("username")
	data.Form.Email = r.FormValue("email")
	data.Form.FullName = r.FormValue("fullname")
	_ = tmpl.Execute(w, data)
}
