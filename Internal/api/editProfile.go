package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EditProfilePageData holds data for edit profile template
type EditProfilePageData struct {
	model.PageData
	Error   string
	Success string
}

// Get_EditProfileHandler renders the edit profile page
func (server *Server) Get_EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || !server.Service.ValidSessions(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user from session
	user := server.Service.Get_UserBySession(cookie.Value)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	categories, _ := server.Service.GetCategories()
	data := EditProfilePageData{
		PageData: model.PageData{
			IsLoggedIn: true,
			User:       user,
			Categories: categories,
			CSSFile:    "/assets/edit-profile.css",
		},
		Error:   r.URL.Query().Get("error"),
		Success: r.URL.Query().Get("success"),
	}

	tmpl, err := template.ParseFiles("./web/templates/root.html", "./web/templates/edit-profile.html")
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
}

// Post_UpdatePasswordHandler handles changing password
func (server *Server) Post_UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || !server.Service.ValidSessions(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user from session
	user := server.Service.Get_UserBySession(cookie.Value)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/edit-profile?error=Invalid+form+data", http.StatusSeeOther)
		return
	}

	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	// Validate passwords match
	if newPassword != confirmPassword {
		http.Redirect(w, r, "/edit-profile?error=New+passwords+do+not+match", http.StatusSeeOther)
		return
	}

	// Update password
	if err := server.Service.UpdateUserPassword(user.ID, currentPassword, newPassword); err != nil {
		http.Redirect(w, r, "/edit-profile?error="+err.Error(), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/edit-profile?success=Password+changed+successfully", http.StatusSeeOther)
}

// Post_UpdatePhotoHandler handles uploading profile photo
func (server *Server) Post_UpdatePhotoHandler(w http.ResponseWriter, r *http.Request) {
	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || !server.Service.ValidSessions(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get user from session
	user := server.Service.Get_UserBySession(cookie.Value)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse multipart form (max 5MB)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		http.Redirect(w, r, "/edit-profile?error=File+too+large", http.StatusSeeOther)
		return
	}

	// Get file from form
	file, handler, err := r.FormFile("photo")
	if err != nil {
		http.Redirect(w, r, "/edit-profile?error=Please+select+a+photo", http.StatusSeeOther)
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		http.Redirect(w, r, "/edit-profile?error=Invalid+file+type.+Please+upload+JPG,+PNG,+or+GIF", http.StatusSeeOther)
		return
	}

	// Create unique filename
	filename := fmt.Sprintf("profile_%d_%d%s", user.ID, time.Now().Unix(), ext)
	uploadPath := filepath.Join(".", "web", "static", "img", filename)

	// Create destination file
	dst, err := os.Create(uploadPath)
	if err != nil {
		http.Redirect(w, r, "/edit-profile?error=Failed+to+save+photo", http.StatusSeeOther)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		http.Redirect(w, r, "/edit-profile?error=Failed+to+save+photo", http.StatusSeeOther)
		return
	}

	// Update user photo in database
	if err := server.Service.UpdateUserPhoto(user.ID, filename); err != nil {
		// Delete uploaded file if database update fails
		os.Remove(uploadPath)
		http.Redirect(w, r, "/edit-profile?error=Failed+to+update+profile+photo", http.StatusSeeOther)
		return
	}

	// Delete old photo if it's not the default
	if user.Photo != "default.png" && user.Photo != "" {
		oldPhotoPath := filepath.Join(".", "web", "static", "img", user.Photo)
		os.Remove(oldPhotoPath) // Ignore error if file doesn't exist
	}

	http.Redirect(w, r, "/edit-profile?success=Profile+photo+updated+successfully", http.StatusSeeOther)
}
