package backend

import (
	"forum/database"
	"net/http"
)

func HandleRoutes() {
	http.HandleFunc("/", HomeHandler)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, database.DB)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, r, database.DB)

	})

	http.HandleFunc("/CreatePost", func(w http.ResponseWriter, r *http.Request) {
		CreatePostHandler(w, r, database.DB)
	})
}
