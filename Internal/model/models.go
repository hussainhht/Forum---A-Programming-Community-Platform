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
	Category      string // Category name for display
	LikeCount     int
	DislikeCount  int
	UserLiked     bool
	UserDisliked  bool
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

type CommentReaction struct {
	CommentID    int // FK to Comment
	UserID       int // FK to User
	ReactionType string
}

type User struct {
	ID        int
	SessionID string
	Username  string
	Email     string
	Password  string
	Photo     string
	CreatedAt time.Time
}

type Comment struct {
	ID           int
	PostID       int // FK to Post
	UserID       int // FK to User
	Content      string
	CreatedAt    time.Time
	Username     string // For display purposes (not stored in DB)
	LikeCount    int
	DislikeCount int
	UserLiked    bool
	UserDisliked bool
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
	LatestPosts        []Post
	Posts              []Post // generic list for discover page
	Post               *Post
	LikeCount          int
	UserLiked          bool
	Categories         []Category
	SelectedCategories []string
	SelectedCategory   string
	SearchQuery        string
	Sort               string
	HasNextPage        bool
	HasPrevPage        bool
	NextPage           int
	PrevPage           int
	User               *User
	ErrorMsg           string
	SuccessMsg         string
	ErrorCode          int
	CSSFile            string
	ExtraCSS           []string
	Theme              any // dynamic theme (struct with exported color fields) or nil
}

type Cell struct {
	Row int
	Col int
}

// CategoryPageData is a unified view model for all category pages
type CategoryPageData struct {
	PageData
	Slug        string   // category slug from URL
	DisplayName string   // user-friendly category name
	SourceURL   string   // URL to the original category source
	Posts       []Post   // list of posts in this category
	CountPosts  int      // total number of posts in this category
	Theme       *CategoryTheme // dynamic color theme injected into template
}

// CategoryTheme holds CSS variable values for category theming.
// These map onto the variables consumed by category-base.css (legacy --go-* kept for compatibility).
type CategoryTheme struct {
	Accent        string
	AccentDark    string
	AccentLight   string
	Secondary     string
	BgPrimary     string
	BgSecondary   string
	BgCard        string
	BgElevated    string
	TextPrimary   string
	TextSecondary string
	TextMuted     string
	Border        string
	BorderLight   string
	Shadow        string
	ShadowStrong  string
	Radius        string
	RadiusSmall   string
	Spacing       string
	BoxShadow     string
}

// ProfileViewData embeds PageData so root layout can access shared fields while
// exposing profile-specific collections.
type ProfileViewData struct {
	PageData
	UserPosts  []Post
	LikedPosts []Post
}

// CreatePostPageData holds the data for rendering the create post page
// CreatePostPageData now embeds model.PageData so it can be rendered inside the root layout.
type CreatePostPageData struct {
	PageData
	Error              string
	SelectedCategories []int
	Title              string
	TempBlocks         []Block
}

// EditProfilePageData holds data for edit profile template
type EditProfilePageData struct {
	PageData
	Error   string
	Success string
}

type LoginPageData struct {
	Error        string
	ShowRegister bool // Flag to show register form instead
	Form         struct {
		EmailOrUsername string
	}
}

type RegisterPageData struct {
	Error        string
	ShowRegister bool // Flag to show register form by default
	Form         struct {
		Username string
		Email    string
	}
}

type PlaygroundPageData struct {
	// sticky form values
	Language string
	Filename string
	LineEnd  string // "lf" | "crlf"
	BOM      string // "nobom" | "bom"
	Code     string

	// server-generated preview
	HighlightedHTML any

	// optional error
	Error string
}
