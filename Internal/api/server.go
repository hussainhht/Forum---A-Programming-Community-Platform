package api

import (
	"fmt"
	"forum/Internal/service"

	"database/sql"
	"net/http"
)

type Server struct {
	Service *service.Service
	Port    int
}

func (server *Server) Start() error {
	router := http.NewServeMux()

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("web/static/"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.HandleFunc("/", server.Get_RootHandler)

	router.HandleFunc("GET /register", server.Get_RegisterHandler)
	router.HandleFunc("POST /register", server.Post_RegisterHandler)

	router.HandleFunc("GET /login", server.Get_LoginHandler)
	router.HandleFunc("POST /login", server.Post_LoginHandler)

	router.HandleFunc("GET /create-post", server.Get_CreatePostHandler)
	router.HandleFunc("POST /create-post", server.Post_CreatePostHandler)

	router.HandleFunc("POST /create-comment", server.Post_CreateCommentHandler)

	router.HandleFunc("GET /post/{id}", server.Get_PostHandler)

	router.HandleFunc("POST /post-reaction", server.Post_ReactionHandler)

	router.HandleFunc("POST /comment-reaction", server.CommentReactionHandler)
	// router.HandleFunc("POST /comment", server.PostCommentHandler)

	router.HandleFunc("GET /logout", server.LogoutHandler)


	return http.ListenAndServe(fmt.Sprintf(":%d", server.Port), router)
}

func NewService(db *sql.DB) *service.Service {
	return &service.Service{DB: db}
}

func NewServer(port int, service *service.Service) *Server {
	return &Server{
		Service: service,
		Port:    port,
	}
}
