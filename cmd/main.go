package main

import (
	"log"
	"net/http"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/route"
	"github.com/easy-comerce/backend/pkg/config"
)

func main() {
	cfg := config.Load()

	db.InitDB(cfg)
	// if err := migrations.RunMigrations(); err != nil {
	// 	log.Printf("Warning: Failed to run migrations: %v", err)
	// }

	mux := http.NewServeMux()
	route.RegisterRoutes(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
