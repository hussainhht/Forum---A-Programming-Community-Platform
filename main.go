package main

import (
	"fmt"
	"forum/Internal/api"
	"forum/database"
	"log"
)

func main() {
	// Initialize DB
	DB, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	service := api.NewService(DB)
	server := api.NewServer(7777, service)
	fmt.Printf("Server running at http://localhost:%d", server.Port)
	log.Fatal(server.Start())
}
