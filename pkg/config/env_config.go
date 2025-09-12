package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/joho/godotenv"
)

var cfg Config

func GetActiveProfile() string {
	return getEnv(constants.EnvActiveProfile, constants.EnvDev)
}

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
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
		constants.EnvKeyPort,
		constants.EnvKeyDBHost,
		constants.EnvKeyDBPort,
		constants.EnvKeyDBUser,
		constants.EnvKeyDBPassword,
		constants.EnvKeyDBName,
		constants.EnvKeyDBSSLMode,
		constants.EnvKeyJWTSecret,
		constants.EnvKeyObjStoreRegion,
		constants.EnvKeyObjStoreBucketName,
		constants.EnvKeyObjStoreAccountID,
		constants.EnvKeyObjStoreAccessKeyID,
		constants.EnvKeyObjStoreAccessKeySecret,
		constants.EnvKeyObjStorePublicDomain,
	}
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	// Get the project root directory
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")

	var envFileName string
	switch GetActiveProfile() {
	case constants.EnvStage:
		envFileName = ".env.stage"
	case constants.EnvProd:
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
		Port:       getEnvWithPanic(constants.EnvKeyPort),
		DBHost:     getEnvWithPanic(constants.EnvKeyDBHost),
		DBPort:     getEnvWithPanic(constants.EnvKeyDBPort),
		DBUser:     getEnvWithPanic(constants.EnvKeyDBUser),
		DBPassword: getEnvWithPanic(constants.EnvKeyDBPassword),
		DBName:     getEnvWithPanic(constants.EnvKeyDBName),
		DBSSLMode:  getEnvWithPanic(constants.EnvKeyDBSSLMode),
		JWTSecret:  getEnvWithPanic(constants.EnvKeyJWTSecret),
		ObjStore: ObjectStoreConfig{
			Region:          getEnvWithPanic(constants.EnvKeyObjStoreRegion),
			BucketName:      getEnvWithPanic(constants.EnvKeyObjStoreBucketName),
			AccountID:       getEnvWithPanic(constants.EnvKeyObjStoreAccountID),
			AccessKeyID:     getEnvWithPanic(constants.EnvKeyObjStoreAccessKeyID),
			AccessKeySecret: getEnvWithPanic(constants.EnvKeyObjStoreAccessKeySecret),
			PublicDomain:    getEnvWithPanic(constants.EnvKeyObjStorePublicDomain),
		},
	}
	return cfg
}

func GetPublicDomain() string {
	publicDomain := cfg.ObjStore.PublicDomain
	if publicDomain == "" {
		publicDomain = getEnv(constants.EnvKeyObjStorePublicDomain, "")
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
