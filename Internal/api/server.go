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

	// C++ zone
	router.HandleFunc("/cpp", server.Get_CPPHandler)
	router.HandleFunc("/c++", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/cpp", http.StatusMovedPermanently)
	})

	// MATLAB zone
	router.HandleFunc("/matlab", server.Get_MATLABHandler)
	router.HandleFunc("/MATLAB", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/matlab", http.StatusMovedPermanently)
	})

	// Bash zone
	router.HandleFunc("/bash", server.Get_BashHandler)
	router.HandleFunc("/Bash", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bash", http.StatusMovedPermanently)
	})

	// Assembly zone
	router.HandleFunc("/assembly", server.Get_AssemblyHandler)
	router.HandleFunc("/Assembly", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assembly", http.StatusMovedPermanently)
	})

	// Python zone
	router.HandleFunc("/python", server.Get_PythonHandler)
	router.HandleFunc("/Python", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/python", http.StatusMovedPermanently)
	})
	// JavaScript zone
	router.HandleFunc("/javascript", server.Get_JavaScriptHandler)
	router.HandleFunc("/js", func(w http.ResponseWriter, r *http.Request) { // short alias
		http.Redirect(w, r, "/javascript", http.StatusMovedPermanently)
	})

	// Brainfuck (Esoteric) zone
	router.HandleFunc("/brainfuck", server.Get_BrainfuckHandler)
	router.HandleFunc("/Brainfuck", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/brainfuck", http.StatusMovedPermanently)
	})

	// PHP zone
	router.HandleFunc("/php", server.Get_PHPHandler)
	router.HandleFunc("/PHP", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/php", http.StatusMovedPermanently)
	})

	// R zone
	router.HandleFunc("/r", server.Get_RHandler)
	router.HandleFunc("/R", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/r", http.StatusMovedPermanently)
	})

	// Lua zone
	router.HandleFunc("/lua", server.Get_LuaHandler)
	router.HandleFunc("/Lua", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/lua", http.StatusMovedPermanently)
	})

	// TypeScript zone
	router.HandleFunc("/typescript", server.Get_TypeScriptHandler)
	router.HandleFunc("/TypeScript", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/typescript", http.StatusMovedPermanently)
	})

	// Swift zone
	router.HandleFunc("/swift", server.Get_SwiftHandler)
	router.HandleFunc("/Swift", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swift", http.StatusMovedPermanently)
	})

	// Dart & Flutter zone
	router.HandleFunc("/dart", server.Get_DartHandler)
	router.HandleFunc("/Dart", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dart", http.StatusMovedPermanently)
	})

	// Kotlin zone
	router.HandleFunc("/kotlin", server.Get_KotlinHandler)
	router.HandleFunc("/Kotlin", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/kotlin", http.StatusMovedPermanently)
	})
	router.HandleFunc("/flutter", server.Get_DartHandler)
	router.HandleFunc("/Flutter", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dart", http.StatusMovedPermanently)
	})

	// SQL page
	router.HandleFunc("/sql", server.Get_SQLHandler)
	router.HandleFunc("/SQL", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/sql", http.StatusMovedPermanently) })

	// Fortran zone
	router.HandleFunc("/fortran", server.Get_FortranHandler)
	router.HandleFunc("/Fortran", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/fortran", http.StatusMovedPermanently)
	})

	// Julia zone
	router.HandleFunc("/julia", server.Get_JuliaHandler)
	router.HandleFunc("/Julia", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/julia", http.StatusMovedPermanently)
	})
	fmt.Printf("Server running at http://localhost:%d\n", server.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", server.Port), router)
}
