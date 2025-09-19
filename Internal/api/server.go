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

	router.HandleFunc("/golang", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/golang", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Golang", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/golang", http.StatusMovedPermanently)
	})

	// Rust zone
	router.HandleFunc("/rust", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/rust", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Rust", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/rust", http.StatusMovedPermanently)
	})

	// Generic category route (new): /category/{slug}
	router.HandleFunc("GET /category/{slug}", server.Get_CategoryHandler)

	// Ruby zone
	router.HandleFunc("/ruby", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/ruby", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Ruby", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/ruby", http.StatusMovedPermanently)
	})

	// Java zone
	router.HandleFunc("/java", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/java", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Java", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/java", http.StatusMovedPermanently)
	})

	// C zone
	router.HandleFunc("/c", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/c", http.StatusMovedPermanently)
	})
	router.HandleFunc("/C", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/c", http.StatusMovedPermanently)
	})

	// CSS zone
	router.HandleFunc("/css", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/css", http.StatusMovedPermanently)
	})
	router.HandleFunc("/CSS", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/css", http.StatusMovedPermanently)
	})

	// C# zone
	router.HandleFunc("/csharp", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/csharp", http.StatusMovedPermanently)
	})
	router.HandleFunc("/CSharp", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/csharp", http.StatusMovedPermanently)
	})

	// HTML zone
	router.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/html", http.StatusMovedPermanently)
	})
	router.HandleFunc("/HTML", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/html", http.StatusMovedPermanently)
	})

	// C++ zone
	router.HandleFunc("/cpp", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/cpp", http.StatusMovedPermanently)
	})
	router.HandleFunc("/c++", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/cpp", http.StatusMovedPermanently)
	})

	// MATLAB zone
	router.HandleFunc("/matlab", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/matlab", http.StatusMovedPermanently)
	})
	router.HandleFunc("/MATLAB", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/matlab", http.StatusMovedPermanently)
	})

	// Bash zone
	router.HandleFunc("/bash", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/bash", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Bash", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/bash", http.StatusMovedPermanently)
	})

	// Assembly zone
	router.HandleFunc("/assembly", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/assembly", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Assembly", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/assembly", http.StatusMovedPermanently)
	})

	// Python zone
	router.HandleFunc("/python", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/python", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Python", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/python", http.StatusMovedPermanently)
	})
	// JavaScript zone
	router.HandleFunc("/javascript", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/javascript", http.StatusMovedPermanently)
	})
	router.HandleFunc("/js", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/javascript", http.StatusMovedPermanently)
	})

	// Brainfuck (Esoteric) zone
	router.HandleFunc("/brainfuck", server.Get_BrainfuckHandler)
	router.HandleFunc("/Brainfuck", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/brainfuck", http.StatusMovedPermanently)
	})

	// PHP zone
	router.HandleFunc("/php", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/php", http.StatusMovedPermanently)
	})
	router.HandleFunc("/PHP", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/php", http.StatusMovedPermanently)
	})

	// R zone
	router.HandleFunc("/r", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/r", http.StatusMovedPermanently)
	})
	router.HandleFunc("/R", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/r", http.StatusMovedPermanently)
	})

	// Lua zone
	router.HandleFunc("/lua", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/lua", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Lua", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/lua", http.StatusMovedPermanently)
	})

	// TypeScript zone
	router.HandleFunc("/typescript", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/typescript", http.StatusMovedPermanently)
	})
	router.HandleFunc("/TypeScript", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/typescript", http.StatusMovedPermanently)
	})

	// Swift zone
	router.HandleFunc("/swift", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/swift", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Swift", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/swift", http.StatusMovedPermanently)
	})

	// Dart & Flutter zone
	router.HandleFunc("/dart", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/dart", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Dart", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/dart", http.StatusMovedPermanently)
	})

	// Kotlin zone
	router.HandleFunc("/kotlin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/kotlin", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Kotlin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/kotlin", http.StatusMovedPermanently)
	})
	router.HandleFunc("/flutter", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/dart", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Flutter", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/dart", http.StatusMovedPermanently)
	})

	// SQL page
	router.HandleFunc("/sql", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/sql", http.StatusMovedPermanently)
	})
	router.HandleFunc("/SQL", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/sql", http.StatusMovedPermanently)
	})

	// Fortran zone
	router.HandleFunc("/fortran", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/fortran", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Fortran", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/fortran", http.StatusMovedPermanently)
	})

	// Julia zone
	router.HandleFunc("/julia", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/julia", http.StatusMovedPermanently)
	})
	router.HandleFunc("/Julia", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/category/julia", http.StatusMovedPermanently)
	})
	fmt.Printf("Server running at http://localhost:%d\n", server.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", server.Port), router)
}
