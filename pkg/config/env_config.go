package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/joho/godotenv"
)

var cfg Config

func GetActiveProfile() string {
	return getEnv(c.EnvActiveProfile, c.EnvDev)
}

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBShowLog  string
	JWTSecret  string
	ObjStore   ObjectStoreConfig
}

type ObjectStoreConfig struct {
	Region          string
	BucketName      string
	AccountID       string
	AccessKeyID     string
	AccessKeySecret string
	PublicDomain    string
}

func Load() Config {
	// Clear existing environment variables to ensure fresh load
	envVars := []string{
		c.EnvKeyPort,
		c.EnvKeyDBHost,
		c.EnvKeyDBPort,
		c.EnvKeyDBUser,
		c.EnvKeyDBPassword,
		c.EnvKeyDBName,
		c.EnvKeyDBSSLMode,
		c.EnvKeyDBShowLog,
		c.EnvKeyJWTSecret,
		c.EnvKeyObjStoreRegion,
		c.EnvKeyObjStoreBucketName,
		c.EnvKeyObjStoreAccountID,
		c.EnvKeyObjStoreAccessKeyID,
		c.EnvKeyObjStoreAccessKeySecret,
		c.EnvKeyObjStorePublicDomain,
	}
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	// Get the project root directory
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	var envFileName string
	switch GetActiveProfile() {
	case c.EnvStage:
		envFileName = ".env.stage"
	case c.EnvProd:
		envFileName = ".env.prod"
	default:
		envFileName = ".env.dev"
	}

	envPath := filepath.Join(projectRoot, envFileName)
	err := godotenv.Load(envPath)
	if err != nil {
		logger.Logger.Warn("Environment file not found", "file", envFileName, "path", envPath, "error", err)
	}

	cfg = Config{
		Port:       getEnvWithPanic(c.EnvKeyPort),
		DBHost:     getEnvWithPanic(c.EnvKeyDBHost),
		DBPort:     getEnvWithPanic(c.EnvKeyDBPort),
		DBUser:     getEnvWithPanic(c.EnvKeyDBUser),
		DBPassword: getEnvWithPanic(c.EnvKeyDBPassword),
		DBName:     getEnvWithPanic(c.EnvKeyDBName),
		DBSSLMode:  getEnvWithPanic(c.EnvKeyDBSSLMode),
		DBShowLog:  getEnvWithPanic(c.EnvKeyDBShowLog),
		JWTSecret:  getEnvWithPanic(c.EnvKeyJWTSecret),
		ObjStore: ObjectStoreConfig{
			Region:          getEnvWithPanic(c.EnvKeyObjStoreRegion),
			BucketName:      getEnvWithPanic(c.EnvKeyObjStoreBucketName),
			AccountID:       getEnvWithPanic(c.EnvKeyObjStoreAccountID),
			AccessKeyID:     getEnvWithPanic(c.EnvKeyObjStoreAccessKeyID),
			AccessKeySecret: getEnvWithPanic(c.EnvKeyObjStoreAccessKeySecret),
			PublicDomain:    getEnvWithPanic(c.EnvKeyObjStorePublicDomain),
		},
	}
	return cfg
}

func GetPublicDomain() string {
	publicDomain := cfg.ObjStore.PublicDomain
	if publicDomain == "" {
		publicDomain = getEnv(c.EnvKeyObjStorePublicDomain, "")
	}
	return publicDomain
}

func getEnvWithPanic(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	panic(fmt.Sprintf("Environment variable %s is not set", key))
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
