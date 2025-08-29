package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sherwin-yu/go-api/database"
	"github.com/sherwin-yu/go-api/handlers"
)

type Server struct {
	db      *database.Database
	handler *handlers.UserHandler
}

func NewServer(db *database.Database) *Server {
	return &Server{
		db:      db,
		handler: handlers.NewUserHandler(db),
	}
}

func (s *Server) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.handler.CreateUser).Methods("POST")
	r.HandleFunc("/users", s.handler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", s.handler.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", s.handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", s.handler.DeleteUser).Methods("DELETE")

	return r
}

func (s *Server) Start(port string) error {
	router := s.SetupRoutes()
	fmt.Printf("Server starting on :%s\n", port)
	return http.ListenAndServe(":"+port, router)
}
