package main

import (
	"log"

	"github.com/sherwin-yu/go/config"
	"github.com/sherwin-yu/go/database"
	"github.com/sherwin-yu/go/server"
)

func main() {
	cfg := config.Load()

	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	srv := server.NewServer(db)
	log.Fatal(srv.Start(cfg.Port))
}
