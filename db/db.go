package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/easy-comerce/backend/pkg/config"
	c "github.com/easy-comerce/backend/pkg/constants"
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
		logger.Error("❌ Failed to get database instance", "error", err)
		return
	}

	if err := sqlDb.Close(); err != nil {
		logger.Error("❌ Database connection closing failed", "error", err)
		return
	}

	logger.Info("Database connection closed successfully")
	db = nil
}

func loadDB() *gorm.DB {
	cfg := config.GetConfig().DBConfig
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	logMode := gormLogger.Error
	if cfg.DBShowLog == "true" {
		logMode = gormLogger.Info
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logMode),
		NowFunc: func() time.Time {
			return timeutil.NowUTC()
		},
	})

	if err != nil {
		logger.Error("❌ Database connection failed", "error", err, "host", cfg.DBHost, "port", cfg.DBPort, "username", cfg.DBUser, "dbname", cfg.DBName, "env", config.GetActiveProfile(), "show_log", cfg.DBShowLog)
		panic(err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil || sqlDB == nil {
		logger.Error("❌ Failed to get database instance", "error", err)
		panic(err)
	}

	sqlDB.SetMaxIdleConns(c.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(c.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.DBConnMaxLifeTimeHour) * time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Duration(c.DBConnMaxIdleTimeMinute) * time.Minute)

	logger.Info("Database connection successful", "host", cfg.DBHost, "port", cfg.DBPort, "username", cfg.DBUser, "dbname", cfg.DBName, "env", config.GetActiveProfile(), "show_log", cfg.DBShowLog)
	return gormDB
}
