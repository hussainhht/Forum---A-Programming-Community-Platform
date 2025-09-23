//go:build notuse
// +build notuse

package api

import (
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
)

type DartPageData struct {
	SourceURL  string
	FlutterURL string
	Posts      []model.Post
	IsLoggedIn bool
	CountPosts int
}

// Get_DartHandler serves the Dart / Flutter page. Primary category: dart; fallback: flutter.
func (server *Server) Get_DartHandler(w http.ResponseWriter, r *http.Request) {
	countpost := 0
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
	}

	// Try dart category first; if not found, try flutter.
	catID, err := server.Service.GetCategoryIDByName("dart")
	if err != nil || catID == 0 {
		// fallback
		catID, err = server.Service.GetCategoryIDByName("flutter")
		if err != nil {
			server.Service.HandleError(w, http.StatusInternalServerError)
			return
		}
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

	data := DartPageData{
		SourceURL:  "https://dart.dev/guides",
		FlutterURL: "https://docs.flutter.dev/",
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

	tpl, err := base.ParseGlob("./web/templates/category/dart.html")
	if err != nil {
		log.Printf("dart: template parse error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "Dart.html", data); err != nil {
		log.Printf("dart: execute error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}
