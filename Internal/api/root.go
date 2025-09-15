package api

import (
	"fmt"
	"forum/Internal/model"
	"html/template"
	"net/http"
)

func (server *Server) Get_HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		server.Service.HandleError(w, http.StatusNotFound)
		return
	}

	sessionIDCookie, _ := r.Cookie("session_id")

	// Create template with custom functions
	tmpl, tmplErr := template.ParseFiles("./web/templates/root.html", "./web/templates/home.html")

	if tmplErr != nil {
		server.Service.HandleError(w, http.StatusInternalServerError)
		fmt.Println("Error creating template:", tmplErr)
		return
	}

	// Get categories for the filter
	categories, err := server.Service.GetCategories()
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	// Get posts - either filtered by categories or all posts
	var posts []model.Post
	categoryIDs := r.URL.Query()["category"]

	if len(categoryIDs) > 0 {
		posts, err = server.Service.GetPostsByCategories(categoryIDs)
	} else {
		posts, err = server.Service.GetAllPosts()
	}

	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	//if the user logged in this will be true
	isLoggedIn := false
	if sessionIDCookie != nil {
		isLoggedIn = server.Service.IsValidSession(sessionIDCookie.Value)
	}

	pageData := model.PageData{
		IsLoggedIn:         isLoggedIn,
		Posts:              posts,
		Categories:         categories,
		SelectedCategories: categoryIDs,
		CSSFile:            "./web/static/css/newtyles.css",
		Cells: 		   generateCells(),
	}

	// Pass posts to the template
	tmpl.Execute(w, pageData)
}

func generateCells() []model.Cell {
	cells := []model.Cell{}
	for r := 0; r < 20; r++ {
		for c := 0; c < 20; c++ {
			cells = append(cells, model.Cell{Row: r, Col: c})
		}
	}
	return cells
}
