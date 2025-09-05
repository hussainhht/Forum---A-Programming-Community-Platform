package api

import (
	"encoding/json"
	"forum/Internal/model"
	"html/template"
	"net/http"
	"strconv"
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
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	sessionID := cookie.Value

	// Clean up expired sessions periodically
	server.cleanupExpiredSessions()

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

	// Use a template with FuncMap providing "contains" like in renderCreatePost
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

	tmpl, err = tmpl.ParseFiles("./web/templates/create-post.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	_ = tmpl.Execute(w, data)
}

// POST handler for creating a post
func (server *Server) Post_CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" || !server.Service.IsValidSession(cookie.Value) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
		if content != "" && (blockType == "code" || blockType == "text" || blockType == "link") {
			block := model.Block{
				Type:    blockType,
				Content: content,
			}

			// Handle link blocks
			if blockType == "link" {
				text, url, isValid := server.Service.ParseMarkdownLink(content)
				if isValid {
					block.Link = &model.Link{
						Text: text,
						URL:  url,
					}
				} else {
					// If not valid markdown format, treat as regular content
					block.Type = "text"
				}
			}

			tempBlocks = append(tempBlocks, block)
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
			renderCreatePost(w, server, title, catIDs, tempBlocks, "Failed to create post: "+err.Error())
			return
		}

		// Clear temp blocks
		server.TempBlocks[sessionID] = []model.Block{}

		// Redirect to home page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	case "clear-session":
		// Clear temp blocks when user wants to start fresh
		server.TempBlocks[sessionID] = []model.Block{}
		renderCreatePost(w, server, "", []int{}, []model.Block{}, "")
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
	tempBlocks []model.Block,
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

// cleanupExpiredSessions removes temp blocks for invalid sessions
func (server *Server) cleanupExpiredSessions() {
	for sessionID := range server.TempBlocks {
		if !server.Service.IsValidSession(sessionID) {
			delete(server.TempBlocks, sessionID)
		}
	}
}
