package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/easy-comerce/backend/pkg/config"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/timeutil"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitializeDB() *gorm.DB {
	once.Do(func() {
		db = loadDB()
		db.Exec("SET search_path TO public")
	})
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("database not initialized. Call InitDB() first")
	}
	return db
}

func CloseDB() {
	if db == nil {
		return
	}

	sqlDb, err := db.DB()
	if err != nil {
		logger.Logger.Error("❌ Failed to get database instance", "error", err)
		return
	}

	if err := sqlDb.Close(); err != nil {
		logger.Logger.Error("❌ Database connection closing failed", "error", err)
		return
	}

	logger.Logger.Info("Database connection closed successfully")
	db = nil
}

func loadDB() *gorm.DB {
	cfg := config.GetConfig().DBConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	logMode := gormLogger.Error
	if cfg.ShowLog == "true" {
		logMode = gormLogger.Info
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logMode),
		NowFunc: func() time.Time {
			return timeutil.NowUTC()
		},
	})

	if err != nil {
		logger.Logger.Error("❌ Database connection failed", "error", err, "host", cfg.Host, "port", cfg.Port, "username", cfg.User, "dbname", cfg.Name, "env", config.GetActiveProfile(), "show_log", cfg.ShowLog)
		panic(err)
	}

	logger.Logger.Info("Database connection successful", "host", cfg.Host, "port", cfg.Port, "username", cfg.User, "dbname", cfg.Name, "env", config.GetActiveProfile(), "show_log", cfg.ShowLog)
	return gormDB
}
