package api

import (
	"forum/Internal/model"
	"html/template"
	"net/http"
)

// ProfilePageData holds data passed to profile template
type ProfilePageData struct {
	User      *model.User
	UserPosts []model.Post
}

// Get_ProfileHandler renders the user's profile page
func (server *Server) Get_ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || !server.Service.ValidSessions(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 2️⃣ Get user from DB
	user := server.Service.Get_UserBySession(cookie.Value)
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 3️⃣ Get user's posts with like/comment counts
	posts, err := server.Service.Get_PostsByUser(user.ID)
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	// 4️⃣ Render template
	data := ProfilePageData{
		User:      user,
		UserPosts: posts,
	}

	tmpl, err := template.ParseFiles("./web/templates/profile.html")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
