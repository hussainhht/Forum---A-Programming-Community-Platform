package main

import (
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

	log.Fatal(server.Start())
}
