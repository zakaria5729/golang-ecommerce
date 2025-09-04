package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/logger"
	"github.com/joho/godotenv"
)

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

	return Config{
		Port:       getEnv(constants.EnvKeyPort, "8080"),
		DBHost:     getEnv(constants.EnvKeyDBHost, "localhost"),
		DBPort:     getEnv(constants.EnvKeyDBPort, "5432"),
		DBUser:     getEnv(constants.EnvKeyDBUser, "postgres"),
		DBPassword: getEnv(constants.EnvKeyDBPassword, "password"),
		DBName:     getEnv(constants.EnvKeyDBName, "easy_commerce"),
		DBSSLMode:  getEnv(constants.EnvKeyDBSSLMode, "disable"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
