package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/sherwin-yu/go-api/models"
)

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

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) CreateUser(user models.User) (*models.User, error) {
	query := `INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := d.db.QueryRow(query, user.Name, user.Email, time.Now()).Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID = id
	return &user, nil
}

func (d *Database) GetUser(id int) (*models.User, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	var user models.User
	err := d.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Database) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, name, email, created_at FROM users ORDER BY id`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (d *Database) UpdateUser(id int, user models.User) (*models.User, error) {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id, name, email, created_at`
	var updatedUser models.User
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
