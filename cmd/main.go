package main

import (
	"net/http"

	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/internal/route"
	c "github.com/easy-comerce/backend/pkg/config"
	dl "github.com/easy-comerce/backend/pkg/data_loader"
	"github.com/easy-comerce/backend/pkg/logger"
)

func main() {
	cfg := c.InitConfig()
	db.InitializeDB()

	if err := dl.InitRoleAndSuperAdmin(); err != nil {
		panic("Failed to create initial role and user: " + err.Error())
	}

	mux := http.NewServeMux()
	route.RegisterAllRoutes(mux)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + cfg.Port,
	}

	logger.Logger.Info("Server starting", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Logger.Error("Failed to start server", "error", err, "port", cfg.Port)
	}
}
