package model

import "time"

type Block struct {
	Type    string `json:"type"` // "text" or "code"
	Content string `json:"content"`
}

type TempPost struct {
    Title              string
    Blocks             []Block
    SelectedCategories []int
}

type Post struct {
	ID        int
	UserID    int // FK to User
	Title     string
	Content   []Block `json:"content"`
	CreatedAt time.Time
	Username  string // For display purposes (not stored in DB)
	LikeCount int
	Comments  []Comment
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
	Username  string // For display purposes (not stored in DB)

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
	IsLoggedIn         bool
	CSRFToken          string
	Post               *Post
	Posts              []Post
	LikeCount          int
	UserLiked          bool
	Categories         []Category
	SelectedCategories []string
	ErrorMsg           string
	ErrorCode          int
}
