package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/easy-comerce/backend/db"
	c "github.com/easy-comerce/backend/pkg/config"
	dl "github.com/easy-comerce/backend/pkg/data_loader"
	"github.com/easy-comerce/backend/pkg/logger"
	l "github.com/easy-comerce/backend/pkg/logger"
	m "github.com/easy-comerce/backend/pkg/middleware"
	"github.com/easy-comerce/backend/pkg/router"
	"github.com/easy-comerce/backend/route"
)

func main() {
	cfg := c.InitConfig()
	db.InitializeDB()

	defer logger.CloseLogFile(true)
	defer db.CloseDB()

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
		m.TrailingSlashMiddleware,
	)

	server := &http.Server{
		Handler: handler,
		Addr:    ":" + cfg.AppConfig.Port,
	}

	go func() {
		l.Info("Server starting", "port", cfg.AppConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("❌ Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	l.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		l.Error("❌ Server forced to shutdown", "error", err)
	}
}
