package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Database struct {
	db *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) CreateUser(user User) (*User, error) {
	query := `INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := d.db.QueryRow(query, user.Name, user.Email, time.Now()).Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return &user, nil
}

func (d *Database) GetUser(id int) (*User, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	var user User
	err := d.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Database) GetAllUsers() ([]User, error) {
	query := `SELECT id, name, email, created_at FROM users ORDER BY id`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (d *Database) UpdateUser(id int, user User) (*User, error) {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email, created_at`
	var updatedUser User
	err := d.db.QueryRow(query, user.Name, user.Email, id).Scan(
		&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (d *Database) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := d.db.Exec(query, id)
	return err
}

type Server struct {
	db *Database
}

func NewServer(db *Database) *Server {
	return &Server{db: db}
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	createdUser, err := s.db.CreateUser(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (s *Server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := s.db.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *Server) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := s.db.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updatedUser, err := s.db.UpdateUser(id, user)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = s.db.DeleteUser(id)
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) setupRoutes() *mux.Router {
	r := mux.NewRouter()
	
	r.HandleFunc("/users", s.createUserHandler).Methods("POST")
	r.HandleFunc("/users", s.getAllUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", s.getUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", s.updateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{id}", s.deleteUserHandler).Methods("DELETE")

	return r
}

func main() {
	connectionString := "postgres://username:password@localhost/dbname?sslmode=disable"
	
	db, err := NewDatabase(connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.db.Close()

	server := NewServer(db)
	router := server.setupRoutes()

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}