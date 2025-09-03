package db

import (
	"fmt"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Shanghai",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		logger.Logger.Error("Failed to connect to database", "error", err, "host", cfg.DBHost, "port", cfg.DBPort, "dbname", cfg.DBName)
		panic(err)
	}

	logger.Logger.Info("Database connected successfully", "host", cfg.DBHost, "port", cfg.DBPort, "dbname", cfg.DBName)
}

func GetDB() *gorm.DB {
	return DB
}
