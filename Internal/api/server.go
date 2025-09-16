package api

import (
	"database/sql"
	"fmt"
	m "forum/Internal/model"
	"forum/Internal/service"
	"net/http"
)

// Server represents the HTTP server and its dependencies
type Server struct {
	Port       int
	Service    *service.Service
	TempBlocks map[string][]m.Block // sessionID -> temporary blocks for create post
}

// NewServer initializes a Server
func NewServer(port int, service *service.Service) *Server {
	return &Server{
		Port:       port,
		Service:    service,
		TempBlocks: make(map[string][]m.Block),
	}
}

// NewService initializes Service from a database
func NewService(db *sql.DB) *service.Service {
	return &service.Service{DB: db}
}

// Start begins listening and routing HTTP requests
func (server *Server) Start() error {
	router := http.NewServeMux()

	// Serve static files (CSS, images)
	fs := http.FileServer(http.Dir("./web"))
	router.Handle("/web/", http.StripPrefix("/web/", fs))

	// Root
	router.HandleFunc("/", server.Get_HomeHandler)

	// Register
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.Get_RegisterHandler(w, r)
		case http.MethodPost:
			server.Post_RegisterHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Login
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.Get_LoginHandler(w, r)
		case http.MethodPost:
			server.Post_LoginHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Logout
	router.HandleFunc("/logout", server.LogoutHandler)

	// Create Post
	router.HandleFunc("/create-post", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.Get_CreatePostHandler(w, r)
		case http.MethodPost:
			server.Post_CreatePostHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Create Comment
	router.HandleFunc("POST /create-comment", server.Post_CreateCommentHandler)

	// Post page (view post)
	router.HandleFunc("GET /post/{id}", server.Get_PostHandler)

	// Reactions
	router.HandleFunc("POST /post-reaction", server.Post_ReactionHandler)
	router.HandleFunc("POST /comment-reaction", server.CommentReactionHandler)

	router.HandleFunc("/playground", server.Get_PlaygroundHandler)
	router.HandleFunc("/playground/preview", server.Post_PlaygroundPreviewHandler)
	router.HandleFunc("/download", server.Post_DownloadHandler)

	router.HandleFunc("/golang", server.Get_GolangHandler)
	// Optional: redirect uppercase variant to canonical path
	router.HandleFunc("/Golang", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/golang", http.StatusMovedPermanently)
	})

	// Rust zone
	router.HandleFunc("/rust", server.Get_RustHandler)
	router.HandleFunc("/Rust", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/rust", http.StatusMovedPermanently)
	})

	// Ruby zone
	router.HandleFunc("/ruby", server.Get_RubyHandler)
	router.HandleFunc("/Ruby", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ruby", http.StatusMovedPermanently)
	})

	// Java zone
	router.HandleFunc("/java", server.Get_JavaHandler)
	router.HandleFunc("/Java", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/java", http.StatusMovedPermanently)
	})

	// C zone
	router.HandleFunc("/c", server.Get_CHandler)
	router.HandleFunc("/C", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/c", http.StatusMovedPermanently)
	})

	// CSS zone
	router.HandleFunc("/css", server.Get_CSSHandler)
	router.HandleFunc("/CSS", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/css", http.StatusMovedPermanently)
	})

	// C# zone
	router.HandleFunc("/csharp", server.Get_CSharpHandler)
	router.HandleFunc("/CSharp", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/csharp", http.StatusMovedPermanently)
	})

	// HTML zone
	router.HandleFunc("/html", server.Get_HTMLHandler)
	router.HandleFunc("/HTML", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/html", http.StatusMovedPermanently)
	})
	// JavaScript zone
	router.HandleFunc("/javascript", server.Get_JavaScriptHandler)
	router.HandleFunc("/js", func(w http.ResponseWriter, r *http.Request) { // short alias
		http.Redirect(w, r, "/javascript", http.StatusMovedPermanently)
	})
	fmt.Printf("Server running at http://localhost:%d\n", server.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", server.Port), router)
}
