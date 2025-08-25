package handlers

type Post struct {
	ID      int
	Title   string
	Content string
	Author  string
}

type User struct {
	ID       int
	Username string
	Password string
}
