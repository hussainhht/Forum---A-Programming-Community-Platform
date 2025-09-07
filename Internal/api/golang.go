package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"net/http"
)

type GolangPageData struct {
	SourceURL string
	Posts     []model.Post
}

func (server *Server) Get_GolangHandler(w http.ResponseWriter, r *http.Request) {
	catID, err := server.Service.GetCategoryIDByName("golang")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	if catID == 0 {
		data := GolangPageData{
			SourceURL: "https://go.dev/doc/",
			Posts:     []model.Post{},
		}
		tpl, err := template.ParseFiles("./web/templates/golang.html", "./web/templates/sidebar.html")
		if err != nil {
			server.Service.HandleError(w, http.StatusInternalServerError)
			return
		}
		if err := tpl.Execute(w, data); err != nil {
			server.Service.HandleError(w, http.StatusInternalServerError)
		}
		return
	}

	posts, err := server.Service.GetPostsByCategories([]string{intToStr(catID)})
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}

	data := GolangPageData{
		SourceURL: "https://go.dev/doc/",
		Posts:     posts,
	}

	tpl, err := template.ParseFiles("./web/templates/golang.html", "./web/templates/sidebar.html")
	if err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, data); err != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
	}
}

func intToStr(n int) string {
	return fmt.Sprintf("%d", n)
}
