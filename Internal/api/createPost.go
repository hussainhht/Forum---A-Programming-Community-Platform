package api

import (
	"encoding/json"
	"forum/Internal/model"
	"html/template"
	"net/http"
	"strconv"
	m "forum/Internal/model"
)

// CreatePostPageData holds the data for rendering the create post page
type CreatePostPageData struct {
	Error              string
	Categories         []model.Category
	SelectedCategories []int
	Title              string
	TempBlocks         []model.Block
}

// GET handler for creating a post
func (server *Server) Get_CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || !server.Service.IsValidSession(cookie.Value) {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}
	sessionID := cookie.Value

	// Initialize TempBlocks for this session
	if _, ok := server.TempBlocks[sessionID]; !ok {
		server.TempBlocks[sessionID] = []model.Block{}
	}

	// Get categories
	categories, err := server.Service.GetCategories()
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	data := CreatePostPageData{
		Title:              "",
		Categories:         categories,
		SelectedCategories: []int{},
		TempBlocks:         server.TempBlocks[sessionID],
		Error:              "",
	}

	tmpl := template.Must(template.ParseFiles("./web/templates/create-post.html"))
	_ = tmpl.Execute(w, data)
}

// POST handler for creating a post
func (server *Server) Post_CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" || !server.Service.IsValidSession(cookie.Value) {
		server.Service.HandleError(w, http.StatusUnauthorized)
		return
	}
	sessionID := cookie.Value

	// Initialize TempBlocks if missing
	if _, ok := server.TempBlocks[sessionID]; !ok {
		server.TempBlocks[sessionID] = []model.Block{}
	}

	if err := r.ParseForm(); err != nil {
		server.Service.HandleError(w, http.StatusBadRequest)
		return
	}

	action := r.FormValue("action")
	title := r.FormValue("title")
	blockType := r.FormValue("type")
	content := r.FormValue("content")
	categories := r.Form["category"]

	// Convert categories to []int
	var catIDs []int
	for _, c := range categories {
		id, _ := strconv.Atoi(c)
		catIDs = append(catIDs, id)
	}

	tempBlocks := server.TempBlocks[sessionID]

	switch action {
	case "add-block":
		if content != "" && (blockType == "code" || blockType == "text") {
			tempBlocks = append(tempBlocks, model.Block{
				Type:    blockType,
				Content: content,
			})
			server.TempBlocks[sessionID] = tempBlocks
		}
		renderCreatePost(w, server, title, catIDs, tempBlocks, "")
		return

	case "remove-block":
		if len(tempBlocks) > 0 {
			tempBlocks = tempBlocks[:len(tempBlocks)-1]
			server.TempBlocks[sessionID] = tempBlocks
		}
		renderCreatePost(w, server, title, catIDs, tempBlocks, "")
		return

	case "submit-post":
		if title == "" || len(tempBlocks) == 0 || len(catIDs) == 0 {
			renderCreatePost(w, server, title, catIDs, tempBlocks, "Title, categories, and blocks are required")
			return
		}

		// Convert blocks to JSON
		blocksJSON, err := json.Marshal(tempBlocks)
		if err != nil {
			renderCreatePost(w, server, title, catIDs, tempBlocks, "Failed to encode blocks")
			return
		}

		// Convert catIDs to []string
		catIDsStr := []string{}
		for _, id := range catIDs {
			catIDsStr = append(catIDsStr, strconv.Itoa(id))
		}

		// Call service to create post
		if err := server.Service.CreatePost(sessionID, title, string(blocksJSON), catIDsStr); err != nil {
			renderCreatePost(w, server, title, catIDs, tempBlocks, "Failed to create post")
			return
		}

		// Clear temp blocks
		server.TempBlocks[sessionID] = []model.Block{}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	default:
		renderCreatePost(w, server, title, catIDs, tempBlocks, "")
		return
	}
}


// Common renderer for Create Post page
func renderCreatePost(
	w http.ResponseWriter,
	server *Server,
	title string,
	selectedCats []int,
	tempBlocks []m.Block,
	errMsg string,
) {
	// Fetch all categories from service
	categories, err := server.Service.GetCategories()
	if err != nil {
		http.Error(w, "Error loading categories", http.StatusInternalServerError)
		return
	}

	// Prepare page data for template
	data := CreatePostPageData{
		Title:              title,
		Categories:         categories,
		SelectedCategories: selectedCats,
		TempBlocks:         tempBlocks,
		Error:              errMsg,
	}

	// Create template with FuncMap for "contains"
	tmpl := template.New("create-post.html").Funcs(template.FuncMap{
		"contains": func(slice []int, val int) bool {
			for _, s := range slice {
				if s == val {
					return true
				}
			}
			return false
		},
	})

	// Parse the template file
	tmpl, err = tmpl.ParseFiles("./web/templates/create-post.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	// Execute template
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
		return
	}
}


