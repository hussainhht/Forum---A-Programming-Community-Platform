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

	router.HandleFunc("/", server.GetRootHandler)

	router.HandleFunc("GET /register", server.GetRegisterHandler)
	router.HandleFunc("POST /register", server.PostRegisterHandler)

	router.HandleFunc("GET /login", server.GetLoginHandler)
	router.HandleFunc("POST /login", server.PostLoginHandler)

	router.HandleFunc("GET /create-post", server.GetCreatePostHandler)
	router.HandleFunc("POST /create-post", server.PostCreatePostHandler)

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
