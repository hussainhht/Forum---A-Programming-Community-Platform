package main

import (
	"forum/database"
)

func main() {
	database.InitDB()
	database.Migrate()
}