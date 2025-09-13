package db

import (
	"fmt"
	"time"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.Config) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
		NowFunc: func() time.Time {
			return timeutil.NowUTC()
		},
	})

	if err != nil {
		logger.Logger.Error("Failed to connect to database", "error", err, "host", cfg.DBHost, "port", cfg.DBPort, "username", cfg.DBUser, "dbname", cfg.DBName, "env", config.GetActiveProfile())
		panic(err)
	}

	logger.Logger.Info("Database connected successfully", "host", cfg.DBHost, "port", cfg.DBPort, "username", cfg.DBUser, "dbname", cfg.DBName, "env", config.GetActiveProfile())
}

func GetDB() *gorm.DB {
	return DB
}
