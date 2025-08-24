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
    http.HandleFunc("/", handlers.HomeHandler(tmpl))


	// Start server
	log.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
