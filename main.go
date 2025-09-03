package main

import (
	"fmt"
	"forum/Internal/api"
	"forum/database"
	"log"
	//m "forum/Internal/model"
)

func main() {
	// Initialize DB
	DB, err := database.InitDB()
	if err != nil {
		log.Fatal("DB init error:", err)
	}

	// Initialize service
	service := api.NewService(DB)

	// Initialize server
	server := api.NewServer(7777, service)

	fmt.Printf("Server running at http://localhost:%d\n", server.Port)
	log.Fatal(server.Start())
}
