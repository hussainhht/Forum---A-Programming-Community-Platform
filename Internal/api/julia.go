package api

import (
	"html/template"
	"net/http"
)

// JuliaPageData holds data for the Julia page template
// CountPosts returns the number of posts for template display
// Posts slice contains posts filtered by Julia category
// SessionValid indicates user authentication state
// Username for greeting if logged in
// Theme reused from other pages for possible future theming
// ActivePage used to highlight nav if implemented later
// IconName can drive dynamic icon usage
// LanguageName for human-readable heading
// CategoryName is the category slug in DB
// Tagline a short descriptor
// MoreLink to official docs or resources
// SourceLink to repository reference
// Badge a short top badge text
// Features array of feature cards (title + description)
// Path canonical route path
// Colors optional map for color theming (not yet used)
// Misc map for any expansions

type JuliaPageData struct {
	CountPosts   int
	Posts        interface{}
	SessionValid bool
	Username     string
	Theme        string
	ActivePage   string
	IconName     string
	LanguageName string
	CategoryName string
	Tagline      string
	MoreLink     string
	SourceLink   string
	Badge        string
	Features     []struct {
		Title string
		Desc  string
	}
	Path   string
	Colors map[string]string
	Misc   map[string]interface{}
}

// Get_JuliaHandler serves the Julia page.
func (server *Server) Get_JuliaHandler(w http.ResponseWriter, r *http.Request) {
	// Session validation (pattern consistent with other handlers)
	sessionValid := false
	username := ""
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		sessionValid = true
		if user, errU := server.Service.GetUserFromSessionID(c.Value); errU == nil && user != nil {
			username = user.Username
		}
	}

	// Retrieve Julia category ID
	catID, err := server.Service.GetCategoryIDByName("julia")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	posts, err := server.Service.GetPostsByCategories([]string{intToStr(catID)})
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	data := JuliaPageData{
		CountPosts:   len(posts),
		Posts:        posts,
		SessionValid: sessionValid,
		Username:     username,
		Theme:        "dark",
		ActivePage:   "julia",
		IconName:     "julia",
		LanguageName: "Julia",
		CategoryName: "julia",
		Tagline:      "High-Performance Dynamic Technical Computing",
		MoreLink:     "https://julialang.org/",
		SourceLink:   "https://github.com/JuliaLang/julia",
		Badge:        "Julia — Scientific & Parallel Power",
		Features: []struct{ Title, Desc string }{
			{Title: "Multiple Dispatch", Desc: "Core paradigm enabling elegant and performant polymorphism."},
			{Title: "High Performance", Desc: "JIT compilation via LLVM rivals C and Fortran speeds."},
			{Title: "Rich Ecosystem", Desc: "Growing scientific, data, and ML package libraries."},
			{Title: "Metaprogramming", Desc: "Powerful macros and generated functions for abstractions."},
		},
		Path: "/julia",
		Colors: map[string]string{
			"primary":   "#9558b2",
			"secondary": "#3f9c9c",
		},
		Misc: map[string]interface{}{},
	}

	tmpl, err := template.ParseFiles("web/templates/Julia.html")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "Julia.html", data); err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
}
