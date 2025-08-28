package model

import "time"

type Post struct {
	ID         int
	UserID     int // FK to User
	CategoryID int // FK to Category
	Title      string
	Content    string
	CreatedAt  time.Time
	Username   string // For display purposes (not stored in DB)
}

type User struct {
	ID        int
	SessionID string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Comment struct {
	ID        int
	PostID    int // FK to Post
	UserID    int // FK to User
	Content   string
	CreatedAt time.Time
}

type Reaction struct {
	ID     int
	PostID int  // FK to Post
	UserID int  // FK to User
	Type   bool // 0: dislike, 1: like
}

type Category struct {
	ID   int
	Name string
}
