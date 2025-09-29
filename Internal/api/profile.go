package api

import (
	"forum/Internal/model"
	"html/template"
	"net/http"
)

// ProfilePageData holds data passed to profile template
type ProfilePageData struct { // kept for backwards compatibility (unused externally now)
	User       *model.User
	UserPosts  []model.Post
	LikedPosts []model.Post
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

	// 4️⃣ Get user's liked posts
	likedPosts, err := server.Service.Get_LikedPostsByUser(user.ID)
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	categories, _ := server.Service.GetCategories()
	// Build unified PageData
	page := model.PageData{
		IsLoggedIn: true,
		User:       user,
		Categories: categories,
	}
	// Reuse Posts slice for user's own posts; provide liked posts via ExtraCSS hack not appropriate -> embed in User field custom? Simpler: extend PageData via template dot chaining with a small struct
	// We'll execute with a composite map to expose additional fields expected by profile template.
	data := map[string]any{
		"Page":       page,
		"User":       user,
		"UserPosts":  posts,
		"LikedPosts": likedPosts,
		"Categories": categories,
	}
	// Parse root + profile
	tmpl, err := template.ParseFiles("./web/templates/root.html", "./web/templates/profile.html")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
}
