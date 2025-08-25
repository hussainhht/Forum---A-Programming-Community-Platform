package main

import (
	"forum/database"
	"forum/handlers"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Initialize DB
	database.InitDB()

	tmpl := template.Must(template.ParseGlob("templates/*.html"))

	// Routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, tmpl)
	})

	// Start server
	log.Println("Server running at http://localhost:7777")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
