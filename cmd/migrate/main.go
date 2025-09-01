package main

import (
	"log"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/db/migrations"
	"github.com/easy-comerce/backend/pkg/config"
)

func main() {
	cfg := config.Load()

	// Initialize database
	db.InitDB(cfg)

	// Run SQL migrations
	log.Println("Running SQL migrations...")
	if err := migrations.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("SQL migrations completed successfully")
}
