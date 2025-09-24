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

	//profile
	router.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.Get_ProfileHandler(w, r)
		// Uncomment below if you plan to allow POST to update profile
		// case http.MethodPost:
		// 	server.Post_ProfileHandler(w, r)
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


	// Generic category route: /category/{slug}
	router.HandleFunc("GET /category/{slug}", server.Get_CategoryHandler)
	fmt.Printf("Server running at http://localhost:%d\n", server.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", server.Port), router)
}
