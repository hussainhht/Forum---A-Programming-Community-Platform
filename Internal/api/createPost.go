package api

import (
	m "forum/Internal/model"
	"html/template"
	"net/http"
)

type CreatePostPageData struct {
	Error              string
	Categories         []m.Category
	Title              string
	SelectedCategories []int
	TempBlocks         []m.Block  // stores blocks before submission
}

var tempPosts = map[string][]m.Block{} // map[sessionID][]Block


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
	// Check user authentication
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
	action := r.FormValue("action") // "add-block", "remove-block", "submit-post"

	// Initialize blocks in request scope (temporary)
	var blocks []m.Block

	switch action {
	case "add-block":
		blockType := r.FormValue("type")
		content := r.FormValue("content")
		if content != "" && (blockType == "text" || blockType == "code") {
			blocks = append(blocks, m.Block{Type: blockType, Content: content})
		}

	case "remove-block":
		if len(blocks) > 0 {
			blocks = blocks[:len(blocks)-1]
		}

	case "submit-post":
		if len(title) > 200 {
			renderCreatePost(w, CreatePostPageData{
				Error: "Title is too long (maximum 200 characters)",
			}, "Title is too long")
			return
		}

		categories := r.Form["category"]
		if title == "" || len(blocks) == 0 {
			allCategories, _ := server.Service.GetCategories()
			renderCreatePost(w, CreatePostPageData{
				Categories: allCategories,
				Error:      "Title and at least one block are required",
			}, "Missing fields")
			return
		}

		if len(categories) == 0 {
			allCategories, _ := server.Service.GetCategories()
			renderCreatePost(w, CreatePostPageData{
				Categories: allCategories,
				Error:      "At least one category is required",
			}, "Missing category")
			return
		}

		sessionID, err := server.Service.GetSessionIDFromCookie(r)
		if err != nil {
			allCategories, _ := server.Service.GetCategories()
			renderCreatePost(w, CreatePostPageData{
				Categories: allCategories,
				Error:      "You must be logged in to create a post",
			}, "Not authenticated")
			return
		}

		if err := server.Service.CreatePost(sessionID, title, blocks, categories); err != nil {
			allCategories, _ := server.Service.GetCategories()
			renderCreatePost(w, CreatePostPageData{
				Categories: allCategories,
				Error:      "The post could not be created",
			}, "Post failed")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Always render the form with current state of blocks
	allCategories, _ := server.Service.GetCategories()
	renderCreatePost(w, CreatePostPageData{
		Categories: allCategories,
		TempBlocks: blocks,
		Title:      title,
	}, "")
}


// Common renderer
func renderCreatePost(w http.ResponseWriter, data CreatePostPageData, errMessage string) {
	tmpl, _ := template.ParseFiles("./web/templates/create-post.html")
	if errMessage != "" {
		data.Error = errMessage
	}
	_ = tmpl.Execute(w, data)
}
