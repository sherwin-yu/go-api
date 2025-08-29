# Go PostgreSQL REST API

A simple REST API built with Go that interacts with a PostgreSQL database.

## Features

- CRUD operations for users
- RESTful endpoints
- JSON request/response handling
- PostgreSQL database integration

## Prerequisites

- Go 1.25 or higher
- PostgreSQL database

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Create a PostgreSQL database and run the schema:
```bash
psql -U your_username -d your_database -f schema.sql
```

3. Update the connection string in `main.go`:
```go
connectionString := "postgres://username:password@localhost/dbname?sslmode=disable"
```

4. Run the application:
```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

- `GET /users` - Get all users
- `GET /users/{id}` - Get user by ID
- `POST /users` - Create a new user
- `PUT /users/{id}` - Update user by ID
- `DELETE /users/{id}` - Delete user by ID

## Example Usage

### Create a user:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

### Get all users:
```bash
curl http://localhost:8080/users
```

### Get user by ID:
```bash
curl http://localhost:8080/users/1
```

### Update a user:
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Smith","email":"johnsmith@example.com"}'
```

### Delete a user:
```bash
curl -X DELETE http://localhost:8080/users/1
```