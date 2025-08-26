package handlers

import (
	_ "forum/backend"
	"forum/database"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// backend.InsertUser()
	if r.URL.Path != "/" {
		return
	}
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	rows, err := database.DB.Query(`
            SELECT posts.id, posts.title, posts.content, users.username
            FROM posts
            JOIN users ON posts.user_id = users.id
            ORDER BY posts.id DESC
        `)
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
			http.Error(w, "Error reading posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	err = tmpl.ExecuteTemplate(w, "home.html", posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	
}

func ServeFiles() {
	cssFS := http.FileServer(http.Dir("./static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", cssFS))
}
