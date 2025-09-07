package main

import (
	"net/http"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/internal/route"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/middleware"
)

func main() {
	cfg := config.Load()

	db.InitDB(cfg)
	// if err := migrations.RunMigrations(); err != nil {
	// 	log.Printf("Warning: Failed to run migrations: %v", err)
	// }

	if err := auth.InitializeDefaultSuperAdmin(); err != nil {
		logger.Logger.Error("Failed to initialize super admin", "error", err)
	}

	mux := http.NewServeMux()
	route.RegisterRoutes(mux)

	handler := middleware.ChainMiddleware(
		middleware.RecoveryMiddleware,
		middleware.CORSMiddleware,
	)(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	logger.Logger.Info("Server starting", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Logger.Error("Failed to start server", "error", err, "port", cfg.Port)
	}
}
