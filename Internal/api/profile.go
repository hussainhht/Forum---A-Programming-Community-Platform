package api

import (
	"forum/Internal/model"
	"html/template"
	"net/http"
)

// Get_ProfileHandler renders the user's profile page
func (server *Server) Get_ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// 1️⃣ Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil || !server.Service.ValidSessions(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 2️⃣ Get user from DB
	user, err := server.Service.GetUserFromSessionID(cookie.Value)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// 3️⃣ Get user's posts with like/comment counts
	posts, err := server.Service.Get_PostsByUser(user.ID)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	// 4️⃣ Get user's liked posts
	likedPosts, err := server.Service.Get_LikedPostsByUser(user.ID)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	categories, _ := server.Service.GetCategories()
	data := model.ProfileViewData{
		PageData: model.PageData{
			IsLoggedIn: true,
			User:       user,
			Categories: categories,
			CSSFile:    "/assets/profile.css", // served dynamically by backend
		},
		UserPosts:  posts,
		LikedPosts: likedPosts,
	}
	tmpl, err := template.ParseFiles("./web/templates/root.html", "./web/templates/profile.html")
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}
}
