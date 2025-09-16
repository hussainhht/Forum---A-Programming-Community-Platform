package model

import "time"

type Block struct {
	Type    string `json:"type"` // "text", "code", or "link"
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

type Link struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

type TempPost struct {
	Title              string
	Blocks             []Block
	SelectedCategories []int
}

type Post struct {
	ID            int
	UserID        int // FK to User
	Title         string
	Content       []Block `json:"content"`
	CreatedAt     time.Time
	Username      string // For display purposes (not stored in DB)
	LikeCount     int
	Comments      []Comment
	Preview       string
	CommentCount  int
	FeaturedScore float64
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
	FeaturedPosts      []Post
	Post               *Post
	LikeCount          int
	UserLiked          bool
	Categories         []Category
	SelectedCategories []string
	User               *User
	ErrorMsg           string
	ErrorCode          int
	CSSFile            string
}

type Cell struct {
	Row int
	Col int
}
