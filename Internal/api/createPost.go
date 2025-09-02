package api

import (
	m "forum/Internal/model"
	"html/template"
	"net/http"
)

type CreatePostPageData struct {
	Error      string
	Categories []m.Category
}

func (server *Server) Get_CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// check if user is authorized
	cookie, err := r.Cookie("session_id")
	if err != nil {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}
	if !server.Service.IsValidSession(cookie.Value) {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}

	// Get categories to display in the form
	categories, err := server.Service.GetCategories()
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}
	categoriesData := CreatePostPageData{
		Categories: categories,
		Error:      "",
	}

	renderCreatePost(w, categoriesData, "")


}

func (server *Server) Post_CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// check user authentication
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}
	valid := server.Service.IsValidSession(cookie.Value)
	if !valid {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		server.Service.HandleError(w, http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")

	if len(title) > 200 {
		renderCreatePost(w, CreatePostPageData{Error: "Title is too long,(the maximum length is 200 characters)"}, "Title is too long")
		return
	}
	content := r.FormValue("content")
	if len(content) > 1000 {
		renderCreatePost(w, CreatePostPageData{Error: "Content is too long,(the maximum length is 1000 characters)"}, "Content is too long")
		return
	}
	categories := r.Form["category"]


	// Get all categories
	allCategories, _ := server.Service.GetCategories()
	//? selectedCategories := []m.Category{}
	categorySet := make(map[string]struct{})
	for _, c := range categories {
		categorySet[c] = struct{}{}
	}
	// for _, cat := range allCategories {
	// 	if _, ok := categorySet[cat.Name]; ok {
	// 		selectedCategories = append(selectedCategories, cat)
	// 	}
	// }

	// validate input
	if title == "" || content == "" {
		
		renderCreatePost(w, CreatePostPageData{Categories: allCategories}, "Title and content are required")
		return
	}

	if len(categories) == 0 {
		renderCreatePost(w, CreatePostPageData{Categories: allCategories, Error: "At least one category is required"}, "At least one category is required")
		return
	}

	sessionID, err := server.Service.GetSessionIDFromCookie(r)
	if err != nil {
		renderCreatePost(w, CreatePostPageData{Categories: allCategories,
			Error: "you must be logged in to create a post"}, "you must be logged in to create a post")
		return
	}

	if err := server.Service.CreatePost(sessionID, title, content, categories); err != nil {
		renderCreatePost(w, CreatePostPageData{Categories: allCategories,
			Error: "the post could not be created"}, "the post could not be created")
		return
	}

	// Success response
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Common renderer
func renderCreatePost(w http.ResponseWriter, data CreatePostPageData, errMessage string) {
	tmpl, _ := template.ParseFiles("./web/templates/create-post.html")
	if errMessage != "" {
		data.Error = errMessage
	}
	_ = tmpl.Execute(w, data)
}
