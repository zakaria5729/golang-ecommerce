package main

import (
	"net/http"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/feature/auth"
	"github.com/easy-comerce/backend/internal/route"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
)

func main() {
	cfg := config.Load()

	db.InitDB(cfg)
	// if err := migrations.RunMigrations(); err != nil {
	// 	log.Printf("Warning: Failed to run migrations: %v", err)
	// }

	// if config.GetActiveProfile() == constants.EnvDev {
	// 	if err := db.GetDB().AutoMigrate(
	// 		&user.User{},
	// 		&role.Role{},
	// 		&permission.Permission{},
	// 		&address.Address{},
	// 		&category.Category{},
	// 		&review.Review{},
	// 		&wishlist.Wishlist{},
	// 		&browsing_history.BrowsingHistory{},
	// 	); err != nil {
	// 		logger.Logger.Error("Failed to auto migrate", "error", err)
	// 	}
	// }

	if err := auth.InitializeDefaultSuperAdmin(); err != nil {
		logger.Logger.Error("Failed to initialize super admin", "error", err)
	}

	mux := http.NewServeMux()
	route.RegisterAllRoutes(mux, cfg)

	// handler := middleware.ChainMiddleware(
	// 	middleware.RecoveryMiddleware,
	// 	middleware.CORSMiddleware,
	// )(mux)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	logger.Logger.Info("Server starting", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Logger.Error("Failed to start server", "error", err, "port", cfg.Port)
	}
}
