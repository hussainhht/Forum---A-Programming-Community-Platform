package backend

import (
	"fmt"
	"forum/database"
	"log"
	"time"
)

func InsertPost(title, content string, userID int, categoryID int, createdAt time.Time) {
	insertPost := `INSERT INTO posts (title, content, user_id, category_id, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(insertPost, title, content, userID, categoryID, createdAt)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadUser() {
	fmt.Println("Reading users from database:")
	query := `SELECT id, username, email, password, created_at FROM users`
	rows, err := database.DB.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email, password, createdAt string
		err := rows.Scan(&id, &username, &email, &password, &createdAt)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("User: ID=%d, Username=%s, Email=%s, CreatedAt=%s\n", id, username, email, createdAt)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
}
