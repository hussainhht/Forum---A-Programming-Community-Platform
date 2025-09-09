// Internal/api/golang.go
package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
)

type GolangPageData struct {
	SourceURL  string
	Posts      []model.Post
	IsLoggedIn bool
	CountPosts int
}

func (server *Server) Get_GolangHandler(w http.ResponseWriter, r *http.Request) {

	countpost := 0
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
	}

	catID, err := server.Service.GetCategoryIDByName("golang")
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

	data := GolangPageData{
		SourceURL:  "https://go.dev/doc/",
		Posts:      posts,
		IsLoggedIn: isLoggedIn,
		CountPosts: countpost,
	}

	// Parse all templates with a FuncMap that includes helpers used by other pages (e.g., create-post.html)
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
		log.Printf("golang: template parse error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "golang.html", data); err != nil {
		log.Printf("golang: execute error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}

func intToStr(n int) string {
	return fmt.Sprintf("%d", n)
}
