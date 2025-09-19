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
	Slug        string
	DisplayName string
	SourceURL   string
	Posts       []model.Post
	IsLoggedIn  bool
	CountPosts  int
}

// Get_CategoryHandler serves any category at /category/{slug} using a single template.
func (server *Server) Get_CategoryHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	// session check
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
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

	data := CategoryPageData{
		Slug:        strings.ToLower(slug),
		DisplayName: display,
		SourceURL:   sourceURLFor(slug),
		Posts:       posts,
		IsLoggedIn:  isLoggedIn,
		CountPosts:  len(posts),
	}

	base := template.New("all").Funcs(template.FuncMap{
		"contains": func(slice []int, val int) bool {
			for _, s := range slice {
				if s == val {
					return true
				}
			}
			return false
		},
	})

	tpl, err := base.ParseFiles("./web/templates/category.html")
	if err != nil {
		log.Printf("category[%s]: template parse error: %v", slug, err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "category.html", data); err != nil {
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
