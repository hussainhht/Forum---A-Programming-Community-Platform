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
	if err != nil || cookie.Value == "" || !server.Service.IsValidSession(cookie.Value) {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}
	sessionID := cookie.Value

	// Parse form data
	if err := r.ParseForm(); err != nil {
		server.Service.HandleError(w, http.StatusBadRequest)
		return
	}

	action := r.FormValue("action") // "add-block", "remove-block", or "submit-post"

	// Load temp blocks for this session
	blocks := tempPosts[sessionID]

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
		title := r.FormValue("title")
		selectedCats := r.Form["category"]
		if title == "" || len(blocks) == 0 || len(selectedCats) == 0 {
			allCats, _ := server.Service.GetCategories()
			renderCreatePost(w, CreatePostPageData{
				Error:      "Title, at least one block, and at least one category are required",
				Categories: allCats,
				TempBlocks: blocks,
				Title:      title,
			}, "")
			return
		}

		// Convert selectedCats (string IDs) to []int
		var catIDs []int
		for _, s := range selectedCats {
			if id, err := strconv.Atoi(s); err == nil {
				catIDs = append(catIDs, id)
			}
		}

		if err := server.Service.CreatePost(sessionID, title, catIDs, blocks); err != nil {
			allCats, _ := server.Service.GetCategories()
			renderCreatePost(w, CreatePostPageData{
				Error:      "Failed to create post",
				Categories: allCats,
				TempBlocks: blocks,
				Title:      title,
			}, "")
			return
		}

		// Clear temp blocks on success
		tempPosts[sessionID] = nil
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Update temp blocks for the session
	tempPosts[sessionID] = blocks

	// Render page with current temp blocks
	allCats, _ := server.Service.GetCategories()
	renderCreatePost(w, CreatePostPageData{
		Categories: allCats,
		TempBlocks: blocks,
		Title:      r.FormValue("title"),
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
