//go:build notuse
// +build notuse

package api

// Moved to Internal/api/notuse/lua.go

//go:build notuse
// +build notuse

package api

import (
	"forum/Internal/model"
	"html/template"
	"log"
	"net/http"
)

type LuaPageData struct {
	SourceURL  string
	Posts      []model.Post
	IsLoggedIn bool
	CountPosts int
}

// Get_LuaHandler serves the Lua page.
func (server *Server) Get_LuaHandler(w http.ResponseWriter, r *http.Request) {
	countpost := 0
	isLoggedIn := false
	if c, err := r.Cookie("session_id"); err == nil && server.Service.IsValidSession(c.Value) {
		isLoggedIn = true
	}

	catID, err := server.Service.GetCategoryIDByName("lua")
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

	data := LuaPageData{
		SourceURL:  "https://www.lua.org/manual/5.4/",
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

	tpl, err := base.ParseGlob("./web/templates/category/lua.html")
	if err != nil {
		log.Printf("lua: template parse error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tpl.ExecuteTemplate(w, "lua.html", data); err != nil {
		log.Printf("lua: execute error: %v", err)
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}
