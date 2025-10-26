package main

import (
	"net/http"

	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/config"
	dl "github.com/easy-comerce/backend/pkg/data_loader"
	l "github.com/easy-comerce/backend/pkg/logger"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
	"github.com/easy-comerce/backend/route"
)

func main() {
	cfg := c.InitConfig()
	db.InitializeDB()

	if err := dl.InitRoleAndSuperAdmin(); err != nil {
		panic("❌ Failed to create initial role and user: " + err.Error())
	}

	baseMux := http.NewServeMux()
	router := router.New(baseMux)
	route.RegisterAllRoutes(router)

	handler := router.Use(
		m.CorsMiddleware,
		m.RecoveryMiddleware,
		m.LoggingMiddleware,
	)

	server := &http.Server{
		Handler: handler,
		Addr:    ":" + cfg.Port,
	}

	l.Logger.Info("Server starting", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		l.Logger.Error("❌ Failed to start server", "error", err, "port", cfg.Port)
	}
}
