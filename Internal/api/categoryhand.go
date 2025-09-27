package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// CategoryPageData is a unified view model for all category pages
type CategoryPageData struct {
	model.PageData
	Slug        string
	DisplayName string
	SourceURL   string
	Posts       []model.Post
	CountPosts  int
	Theme       *CategoryTheme // dynamic color theme injected into template
}

// CategoryTheme holds CSS variable values for category theming.
// These map onto the variables consumed by category-base.css (legacy --go-* kept for compatibility).
type CategoryTheme struct {
	Accent        string
	AccentDark    string
	AccentLight   string
	Secondary     string
	BgPrimary     string
	BgSecondary   string
	BgCard        string
	BgElevated    string
	TextPrimary   string
	TextSecondary string
	TextMuted     string
	Border        string
	BorderLight   string
	Shadow        string
	ShadowStrong  string
	Radius        string
	RadiusSmall   string
	Spacing       string
	BoxShadow     string
}

// themeFor returns a CategoryTheme with sensible defaults per language slug.
// If a slug is unknown, a generic purple/blue theme is returned.
func themeFor(slug string) *CategoryTheme {
	base := &CategoryTheme{
		Accent:        "#00d4ff",
		AccentDark:    "#0099cc",
		AccentLight:   "#33ddff",
		Secondary:     "#7c3aed",
		BgPrimary:     "#0a0e1a",
		BgSecondary:   "#1a1f2e",
		BgCard:        "#242938",
		BgElevated:    "#2d3748",
		TextPrimary:   "#f7fafc",
		TextSecondary: "#cbd5e0",
		TextMuted:     "#a0aec0",
		Border:        "#4a5568",
		BorderLight:   "#2d3748",
		Shadow:        "0 10px 25px rgba(0, 212, 255, 0.1)",
		ShadowStrong:  "0 20px 40px rgba(0, 212, 255, 0.2)",
		BoxShadow:     "0 4px 12px rgba(255, 255, 255, 0)", // added for lang buttons and floating tags
		Radius:        "16px",
		RadiusSmall:   "8px",
		Spacing:       "24px",
	}

	switch strings.ToLower(slug) {
	case "golang": // keep defaults (Go cyan + purple)
		return base
	case "python":
		base.Accent = "#3776ab"
		base.AccentDark = "#27567e"
		base.AccentLight = "#4d8ec4"
		base.Secondary = "#ffdf5a"
		base.Shadow = "0 10px 25px rgba(55,118,171,0.18)"
		base.ShadowStrong = "0 20px 40px rgba(55,118,171,0.30)"
		base.BoxShadow = "0 4px 12px rgba(55,118,171,0.3)"
	case "javascript", "js":
		base.Accent = "#f7df1e"
		base.AccentDark = "#d4b400"
		base.AccentLight = "#ffe955"
		base.Secondary = "#323330"
		base.Shadow = "0 10px 25px rgba(247,223,30,0.15)"
		base.ShadowStrong = "0 20px 40px rgba(247,223,30,0.28)"
		base.BoxShadow = "0 4px 12px rgba(255, 233, 0, 0.3)"
		// base.AccentLight = "#ffe955"
		
	case "typescript", "ts":
		base.Accent = "#3178c6"
		base.AccentDark = "#255a92"
		base.AccentLight = "#4b8fd2"
		base.Secondary = "#2d2d30"
		base.Shadow = "0 10px 25px rgba(49,120,198,0.18)"
		base.ShadowStrong = "0 20px 40px rgba(49,120,198,0.30)"
		base.BoxShadow = "0 4px 12px rgba(49,120,198,0.3)"
	case "rust":
		base.Accent = "#dea584"
		base.AccentDark = "#b06a44"
		base.AccentLight = "#efbfa8"
		base.Secondary = "#ce422b"
		base.Shadow = "0 10px 25px rgba(206,66,43,0.15)"
		base.ShadowStrong = "0 20px 40px rgba(206,66,43,0.28)"
		base.BoxShadow = "0 4px 12px rgba(206,66,43,0.3)"
	case "java":
		base.Accent = "#d34949ff"
		base.AccentDark = "#ed5555ff"
		base.AccentLight = "#ff7777"
		base.Secondary = "#007396"
		base.Shadow = "0 10px 25px rgba(255,79,79,0.15)"
		base.ShadowStrong = "0 20px 40px rgba(255,79,79,0.30)"
		base.BoxShadow = "0 4px 12px rgba(255,79,79,0.3)"

	case "ruby":
		base.Accent = "#d70037"
		base.AccentDark = "#a40026"
		base.AccentLight = "#ff5f7f"
		base.Secondary = "#701c3f"
		base.Shadow = "0 10px 25px rgba(215,0,55,0.15)"
		base.ShadowStrong = "0 20px 40px rgba(215,0,55,0.30)"
		base.BoxShadow = "0 4px 12px rgba(215,0,55,0.3)"
	case "cpp", "c++":
		base.Accent = "#00599c"
		base.AccentDark = "#00426f"
		base.AccentLight = "#1a6fb0"
		base.Secondary = "#9c033a"
	case "csharp", "c#":
		base.Accent = "#68217a"
		base.AccentDark = "#4d1859"
		base.AccentLight = "#8630a0"
		base.Secondary = "#239120"
	case "php":
		base.Accent = "#777bb3"
		base.AccentDark = "#5a5d85"
		base.AccentLight = "#8f93c3"
		base.Secondary = "#4f5b93"
	case "swift":
		base.Accent = "#fa7343"
		base.AccentDark = "#cc5930"
		base.AccentLight = "#ff8a60"
		base.Secondary = "#ffaf43"
	case "kotlin":
		base.Accent = "#7f52ff"
		base.AccentDark = "#623fcc"
		base.AccentLight = "#9875ff"
		base.Secondary = "#ff8a00"
	case "dart", "flutter":
		base.Accent = "#0175c2"
		base.AccentDark = "#015a95"
		base.AccentLight = "#2990d4"
		base.Secondary = "#13b9fd"
	case "bash":
		base.Accent = "#3eaf2c"
		base.AccentDark = "#2d7f20"
		base.AccentLight = "#56c645"
		base.Secondary = "#5c5c5c"
	}
	return base
}

// Get_CategoryHandler serves any category at /category/{slug} using a single template.
func (server *Server) Get_CategoryHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	// session & user
	var user *model.User
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
		user, _ = server.Service.GetUserFromSessionID(c.Value)
	}

	// Resolve category by slug
	catID, err := server.Service.GetCategoryIDByName(strings.ToLower(slug))
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	var posts []model.Post
	if catID != 0 {
		posts, err = server.Service.GetPostsByCategories([]string{fmt.Sprintf("%d", catID)})
		if err != nil {
			server.Service.HandleError(w, http.StatusInternalServerError)
			return
		}
	} else {
		posts = []model.Post{}
	}

	// best-effort display name
	display := prettifyName(slug)

	allCats, err := server.Service.GetCategories()
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	data := CategoryPageData{
		Slug:        strings.ToLower(slug),
		DisplayName: display,
		SourceURL:   sourceURLFor(slug),
		Posts:       posts,
		CountPosts:  len(posts),
		Theme:       themeFor(slug),
		PageData: model.PageData{
			IsLoggedIn: isLoggedIn,
			User:       user,
			Categories: allCats,
			CSSFile:    "/web/static/css/newtyles.css",                  // (unused old base)
			ExtraCSS:   []string{"/web/static/css/" + cssFileFor(slug)}, // still load category specific overrides if exist
		},
	}

	base := template.New("root.html").Funcs(template.FuncMap{
		"contains": func(slice []int, val int) bool {
			for _, s := range slice {
				if s == val {
					return true
				}
			}
			return false
		},
	})

	tpl, err := base.ParseFiles("./web/templates/root.html", "./web/templates/category.html")
	if err != nil {
		log.Printf("category[%s]: template parse error: %v", slug, err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.Execute(w, data); err != nil {
		log.Printf("category[%s]: execute error: %v", slug, err)
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}

func prettifyName(slug string) string {
	s := strings.ToLower(slug)
	// special cases
	switch s {
	case "cpp", "c++":
		return "C++"
	case "csharp", "c#":
		return "C#"
	case "golang":
		return "Golang"
	case "sql":
		return "SQL"
	case "html":
		return "HTML"
	case "css":
		return "CSS"
	case "js", "javascript":
		return "JavaScript"
	case "ts", "typescript":
		return "TypeScript"
	case "r":
		return "R"
	default:
		// Title-case the first letter only
		return strings.ToUpper(s[:1]) + s[1:]
	}
}

func sourceURLFor(slug string) string {
	switch strings.ToLower(slug) {
	case "golang":
		return "https://go.dev/doc/"
	case "rust":
		return "https://www.rust-lang.org/learn"
	case "python":
		return "https://docs.python.org/"
	case "java":
		return "https://docs.oracle.com/javase/"
	case "javascript", "js":
		return "https://developer.mozilla.org/docs/Web/JavaScript"
	case "typescript":
		return "https://www.typescriptlang.org/docs/"
	case "c":
		return "https://en.cppreference.com/w/c"
	case "cpp":
		return "https://en.cppreference.com/w/"
	case "csharp":
		return "https://learn.microsoft.com/dotnet/csharp/"
	case "html":
		return "https://developer.mozilla.org/docs/Web/HTML"
	case "css":
		return "https://developer.mozilla.org/docs/Web/CSS"
	case "sql":
		return "https://dev.mysql.com/doc/"
	case "php":
		return "https://www.php.net/docs.php"
	case "kotlin":
		return "https://kotlinlang.org/docs/home.html"
	case "dart", "flutter":
		return "https://dart.dev/guides"
	case "swift":
		return "https://www.swift.org/documentation/"
	case "fortran":
		return "https://fortran-lang.org/learn/"
	case "lua":
		return "https://www.lua.org/docs.html"
	case "julia":
		return "https://docs.julialang.org/"
	case "r":
		return "https://cran.r-project.org/manuals.html"
	case "bash":
		return "https://www.gnu.org/software/bash/manual/bash.html"
	case "assembly":
		return "https://www.felixcloutier.com/x86/"
	case "brainfuck":
		return "https://esolangs.org/wiki/Brainfuck"
	case "matlab":
		return "https://www.mathworks.com/help/matlab/"
	default:
		return "#"
	}
}

// cssFileFor returns the css filename for a given slug.
// It normalizes popular aliases to existing css files under web/static/css.
func cssFileFor(slug string) string {
	s := strings.ToLower(slug)
	switch s {
	case "js", "javascript":
		return "javascript.css"
	case "ts", "typescript":
		return "typescript.css"
	case "c++", "cpp":
		return "cpp.css"
	case "c#", "csharp":
		return "csharp.css"
	default:
		// assume we have {slug}.css
		return s + ".css"
	}
}
