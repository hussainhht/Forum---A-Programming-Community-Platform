package main

import (
	"forum/backend"
	"forum/database"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize DB
	database.InitDB()
	if err := database.ImportUsersCSV(database.DB, "users.csv"); err != nil {
		log.Println("initial import error:", err)
	}
	database.AutoSyncCSVToDB("users.csv", 2*time.Second)

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
