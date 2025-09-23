//go:build notuse
// +build notuse

package api

import (
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
)

type BrainfuckPageData struct {
	SourceURL  string
	Posts      []model.Post
	IsLoggedIn bool
	CountPosts int
}

func (server *Server) Get_BrainfuckHandler(w http.ResponseWriter, r *http.Request) {
	countpost := 0
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
	}

	var posts []model.Post
	var catID int
	var err error

	// Try primary category name, then fallback to a general 'esoteric' if present
	catID, err = server.Service.GetCategoryIDByName("brainfuck")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	if catID == 0 {
		catID, err = server.Service.GetCategoryIDByName("esoteric")
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

	data := BrainfuckPageData{
		SourceURL:  "https://esolangs.org/wiki/Brainfuck",
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

	tpl, err := base.ParseGlob("./web/templates/category/Brainfuck.html")
	if err != nil {
		log.Printf("brainfuck: template parse error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "Brainfuck.html", data); err != nil {
		log.Printf("brainfuck: execute error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}
