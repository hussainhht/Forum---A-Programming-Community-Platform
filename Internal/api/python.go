package api

import (
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
)

type PythonPageData struct {
	SourceURL  string
	Posts      []model.Post
	IsLoggedIn bool
	CountPosts int
}

// Get_PythonHandler serves the Python zone page.
func (server *Server) Get_PythonHandler(w http.ResponseWriter, r *http.Request) {
	countpost := 0
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
	}

	catID, err := server.Service.GetCategoryIDByName("python")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	var posts []model.Post
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

	data := PythonPageData{
		SourceURL:  "https://docs.python.org/",
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
		log.Printf("python: template parse error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "python.html", data); err != nil {
		log.Printf("python: execute error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}
