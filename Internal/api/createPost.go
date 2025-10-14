package api

import (
	"encoding/json"
	"forum/Internal/model"
	"html/template"
	"net/http"
	"strconv"
)

// CreatePostPageData holds the data for rendering the create post page
// CreatePostPageData now embeds model.PageData so it can be rendered inside the root layout.
type CreatePostPageData struct {
	model.PageData
	Error              string
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

	// User data
	user, _ := server.Service.GetUserFromSessionID(sessionID)

	data := CreatePostPageData{
		Title:              "",
		SelectedCategories: []int{},
		TempBlocks:         server.TempBlocks[sessionID],
		Error:              "",
		PageData: model.PageData{
			IsLoggedIn: true,
			User:       user,
			Categories: categories,
			CSSFile:    "./web/static/css/newtyles.css",
		},
	}

	// Use standalone create-post template (no root layout)
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

	// Session & user already validated in caller; user retrieval omitted here for simplicity.

	pageData := model.PageData{
		IsLoggedIn: true,
		Categories: categories,
		CSSFile:    "./web/static/css/newtyles.css",
	}

	data := CreatePostPageData{
		PageData:           pageData,
		Title:              title,
		SelectedCategories: selectedCats,
		TempBlocks:         tempBlocks,
		Error:              errMsg,
	}

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
