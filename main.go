package main

import (
	"forum/database"
	"forum/handlers"
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

	handlers.ServeFiles()
	// Routes
	http.HandleFunc("/", handlers.HomeHandler)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		handlers.LoginHandler(w, r, database.DB)

	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, database.DB)
		if err := database.ExportUsersCSV(database.DB, "users.csv"); err != nil {
			log.Println("export users.csv error:", err)
		}

	})

	// Start server
	log.Println("Server running at http://localhost:7777")
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		log.Fatal(err)
	}
}
