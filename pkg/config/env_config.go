package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	c "github.com/easy-comerce/backend/pkg/constants"
	l "github.com/easy-comerce/backend/pkg/logger"
	"github.com/easy-comerce/backend/pkg/utils"
	"github.com/joho/godotenv"
)

var (
	cfg  *Config
	once sync.Once
)

func InitConfig() *Config {
	once.Do(func() {
		cfg = loadConfig()
	})
	return cfg
}

func GetConfig() *Config {
	if cfg == nil {
		panic("config not initialized. Call InitConfig() first")
	}
	return cfg
}

func GetActiveProfile() string {
	return getEnv(c.EnvKeyActiveProfile, c.EnvStage)
}

func loadConfig() *Config {
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
		c.EnvKeyDomainURL,
		c.EnvKeyFcmServerKey,
		c.EnvKeyFcmUrl,
		c.EnvKeyGoogleClientID,
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
	var envFileName string
	switch GetActiveProfile() {
	case c.EnvStage:
		envFileName = ".env.stage"
	case c.EnvProd:
		envFileName = ".env.prod"
	default:
		envFileName = ".env.dev"
	}

	envPath := filepath.Join(utils.GetProjectRootPath(), envFileName)
	err := godotenv.Load(envPath)
	if err != nil {
		l.Logger.Warn("Environment file not found", "file", envFileName, "path", envPath, "error", err)
	}

	config := &Config{
		Port:               getEnvWithPanic(c.EnvKeyPort),
		DBHost:             getEnvWithPanic(c.EnvKeyDBHost),
		DBPort:             getEnvWithPanic(c.EnvKeyDBPort),
		DBUser:             getEnvWithPanic(c.EnvKeyDBUser),
		DBPassword:         getEnvWithPanic(c.EnvKeyDBPassword),
		DBName:             getEnvWithPanic(c.EnvKeyDBName),
		DBSSLMode:          getEnvWithPanic(c.EnvKeyDBSSLMode),
		DBShowLog:          getEnvWithPanic(c.EnvKeyDBShowLog),
		JWTSecret:          getEnvWithPanic(c.EnvKeyJWTSecret),
		DomainURL:          getEnv(c.EnvKeyDomainURL, ""),
		FcmServerKey:       getEnv(c.EnvKeyJWTSecret, ""),
		FcmUrl:             getEnv(c.EnvKeyJWTSecret, ""),
		GoogleClientID:     getEnv(c.EnvKeyGoogleClientID, ""),
		FacebookAppID:      getEnv(c.EnvKeyFacebookAppID, ""),
		SuperAdminEmail:    getEnvWithPanic(c.EnvSuperAdminEmail),
		SuperAdminPassword: getEnvWithPanic(c.EnvSuperAdminPassword),
		ObjStore: ObjectStoreConfig{
			Region:          getEnvWithPanic(c.EnvKeyObjStoreRegion),
			BucketName:      getEnvWithPanic(c.EnvKeyObjStoreBucketName),
			AccountID:       getEnvWithPanic(c.EnvKeyObjStoreAccountID),
			AccessKeyID:     getEnvWithPanic(c.EnvKeyObjStoreAccessKeyID),
			AccessKeySecret: getEnvWithPanic(c.EnvKeyObjStoreAccessKeySecret),
			PublicDomain:    getEnvWithPanic(c.EnvKeyObjStorePublicDomain),
		},
	}

	l.Logger.Info("Config initialized successfully", "env", GetActiveProfile())
	return config
}

func GetStorageDomain() string {
	storageDomain := GetConfig().ObjStore.PublicDomain
	if storageDomain == "" {
		storageDomain = getEnv(c.EnvKeyObjStorePublicDomain, "")
	}
	return storageDomain
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
