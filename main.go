package main

import (
	"forum/database"
	"forum/handlers"
	"log"
	"net/http"
)

func main() {
	// Initialize DB
	database.InitDB()

	handlers.ServeFiles()
	// Routes
	http.HandleFunc("/", handlers.HomeHandler)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		handlers.LoginHandler(w, r, database.DB)

	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, database.DB)
	})

	// Start server
	log.Println("Server running at http://localhost:7777")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
