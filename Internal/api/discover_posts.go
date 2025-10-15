package api

import (
	"forum/Internal/model"
	"html/template"
	"net/http"
	"strconv"
)

// Get_DiscoverPostsHandler serves the Discover Posts page at /posts
func (server *Server) Get_DiscoverPostsHandler(w http.ResponseWriter, r *http.Request) {
	sessionIDCookie, _ := r.Cookie("session_id")

	// Parse query params
	q := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	sort := r.URL.Query().Get("sort")

	// Normalize empty strings to empty (to show all posts initially)
	if q == "" {
		q = ""
	}
	if category == "" || category == "All" {
		category = ""
	}
	if sort == "" {
		sort = "latest"
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}
	const pageSize = 12
	offset := (page - 1) * pageSize

	// template
	tmpl, tmplErr := template.ParseFiles("./web/templates/root.html", "./web/templates/DiscoverPosts.html")
	if tmplErr != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	// categories for select
	categories, err := server.Service.GetCategories()
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	// fetch posts
	posts, hasNext, err := server.Service.GetDiscoverPosts(q, category, sort, pageSize, offset)
	if err != nil {
		server.Service.HandleError(w, r, http.StatusInternalServerError)
		return
	}

	isLoggedIn := false
	if sessionIDCookie != nil {
		isLoggedIn = server.Service.IsValidSession(sessionIDCookie.Value)
	}
	var user *model.User
	if isLoggedIn {
		user, _ = server.Service.GetUserFromSessionID(sessionIDCookie.Value)
	}

	pageData := model.PageData{
		IsLoggedIn:       isLoggedIn,
		User:             user,
		Categories:       categories,
		Posts:            posts,
		SelectedCategory: category,
		SearchQuery:      q,
		Sort:             sort,
		HasNextPage:      hasNext,
		HasPrevPage:      page > 1,
		NextPage:         page + 1,
		PrevPage:         page - 1,
		CSSFile:          "/web/static/css/newtyles.css",
		ExtraCSS:         []string{"/web/static/css/discover.css"},
	}

	tmpl.Execute(w, pageData)
}
