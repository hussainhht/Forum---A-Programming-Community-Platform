package api

import (
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
)

type JavaScriptPageData struct {
	SourceURL  string
	Posts      []model.Post
	IsLoggedIn bool
	CountPosts int
}

// Get_JavaScriptHandler serves the JavaScript zone page.
func (server *Server) Get_JavaScriptHandler(w http.ResponseWriter, r *http.Request) {
	countpost := 0
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
	}

	// Try category names in order: "javascript", then fallback to "js"
	var posts []model.Post
	var catID int
	var err error

	catID, err = server.Service.GetCategoryIDByName("javascript")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	if catID == 0 { // fallback
		catID, err = server.Service.GetCategoryIDByName("js")
		if err != nil {
			server.Service.HandleError(w, http.StatusInternalServerError)
			return
		}
	}

	if catID != 0 {
		posts, err = server.Service.GetPostsByCategories([]string{intToStr(catID)})
		if err != nil {
			server.Service.HandleError(w, http.StatusInternalServerError)
			return
		}
		countpost = len(posts)
	} else {
		posts = []model.Post{}
	}

	data := JavaScriptPageData{
		SourceURL:  "https://developer.mozilla.org/docs/Web/JavaScript",
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
		CountPosts: countpost,
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

	tpl, err := base.ParseGlob("./web/templates/*.html")
	if err != nil {
		log.Printf("javascript: template parse error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "js.html", data); err != nil {
		log.Printf("javascript: execute error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}
