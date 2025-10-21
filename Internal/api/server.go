package api

import (
	"database/sql"
	"fmt"
	m "forum/Internal/model"
	"forum/Internal/service"
	"net/http"
	"os"
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

	// Dynamic CSS assets (example: profile.css served only when needed)
	router.HandleFunc("GET /assets/profile.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
		data, err := os.ReadFile("./web/static/css/profile.css")
		if err != nil {
			server.Service.HandleError(w, r, http.StatusNotFound)
			return
		}
		w.Write(data)
	})

	router.HandleFunc("GET /assets/edit-profile.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes
		data, err := os.ReadFile("./web/static/css/edit-profile.css")
		if err != nil {
			server.Service.HandleError(w, r, http.StatusNotFound)
			return
		}
		w.Write(data)
	})

	// Discover Posts
	router.HandleFunc("GET /posts", server.Get_DiscoverPostsHandler)

	// Register
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.Get_RegisterHandler(w, r)
		case http.MethodPost:
			server.Post_RegisterHandler(w, r)
		default:
			server.Service.HandleError(w, r, http.StatusMethodNotAllowed)
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
			server.Service.HandleError(w, r, http.StatusMethodNotAllowed)
		}
	})

	// Edit Profile
	router.HandleFunc("GET /edit-profile", server.Get_EditProfileHandler)
	router.HandleFunc("POST /edit-profile/password", server.Post_UpdatePasswordHandler)
	router.HandleFunc("POST /edit-profile/photo", server.Post_UpdatePhotoHandler)

	// Login
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			server.Get_LoginHandler(w, r)
		case http.MethodPost:
			server.Post_LoginHandler(w, r)
		default:
			server.Service.HandleError(w, r, http.StatusMethodNotAllowed)
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
			server.Service.HandleError(w, r, http.StatusMethodNotAllowed)
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

	// Help page
	router.HandleFunc("GET /help", server.Get_HelpHandler)

	// Privacy & Terms page
	router.HandleFunc("GET /privacy-terms", server.Get_PrivacyTermsHandler)

	// Generic category route: /category/{slug}
	router.HandleFunc("GET /category/{slug}", server.Get_CategoryHandler)
	fmt.Printf("Server running at http://localhost:%d\n", server.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", server.Port), router)
}
