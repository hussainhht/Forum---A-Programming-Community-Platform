package model

import "time"

type Post struct {
	ID        int
	UserID    int // FK to User
	Title     string
	Content   string
	CreatedAt time.Time
	Username  string // For display purposes (not stored in DB)
	LikeCount int
}

type PostReaction struct {
	PostID       int // FK to Post
	UserID       int // FK to User
	ReactionType string
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

type PageData struct {
	IsLoggedIn bool       `json:"is_logged_in"`
	CSRFToken  string     `json:"csrf_token"`
	Post       *Post      `json:"post"`
	Comments   []Comment  `json:"comments"`
	LikeCount  int        `json:"like_count"`
	UserLiked  bool       `json:"user_liked"`
	Categories []Category `json:"categories"`
}
