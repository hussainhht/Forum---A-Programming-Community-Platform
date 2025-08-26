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

	// Start server
	log.Println("Server running at http://localhost:7777")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
