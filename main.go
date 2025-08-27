package main

import (
	"forum/backend"
	"forum/database"
	"log"
	"net/http"
)

func main() {
	// Initialize DB
	database.InitDB()

	backend.ServeFiles()
	// Routes
	backend.HandleRoutes()
	// Start server
	log.Println("Server running at http://localhost:7777")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
