package main

import (
	"github.com/easy-comerce/backend/db"
	"github.com/easy-comerce/backend/db/migrations"
	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
)

func main() {
	cfg := config.Load()

	db.InitDB(cfg)
	logger.Logger.Info("Running SQL migrations...")
	if err := migrations.RunMigrations(); err != nil {
		logger.Logger.Error("Failed to run migrations", "error", err)
		panic(err)
	}
	logger.Logger.Info("SQL migrations completed successfully")
}
