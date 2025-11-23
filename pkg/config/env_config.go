package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	c "github.com/easy-comerce/backend/pkg/constants"
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
	return getEnv(c.EnvKeyActiveProfile, c.EnvLocal)
}

func GetUuidType() string {
	uuidType := GetConfig().AppConfig.UuidType
	if strings.EqualFold(uuidType, c.UuidRandom) {
		return c.UuidRandom
	}
	return c.UuidSequence
}

func loadConfig() *Config {
	envVars := []string{
		c.EnvKeyHost,
		c.EnvKeyPort,
		c.EnvKeyDBHost,
		c.EnvKeyDBPort,
		c.EnvKeyDBUser,
		c.EnvKeyDBPassword,
		c.EnvKeyDBName,
		c.EnvKeyDBSSLMode,
		c.EnvKeyDBShowConsoleLog,
		c.EnvKeyUuidType,
		c.EnvKeyDBStatsLogType,
		c.EnvKeyJWTSecret,
		c.EnvKeyDomainURL,
		c.EnvKeyFcmServerKey,
		c.EnvKeyFcmUrl,
		c.EnvKeyGoogleClientID,
		c.EnvKeyFacebookAppID,
		c.EnvKeyObjStoreRegion,
		c.EnvKeyObjStoreBucketName,
		c.EnvKeyObjStoreAccountID,
		c.EnvKeyObjStoreAccessKeyID,
		c.EnvKeyObjStoreAccessKeySecret,
		c.EnvKeyObjStorePublicDomain,
		c.EnvKeySocialLoginDefaultPassword,
		c.EnvKeyAppHealthCheckToken,
		c.EnvSuperAdminEmail,
		c.EnvSuperAdminPassword,
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
	case c.EnvLocal:
		envFileName = ".env.local"
	default:
		envFileName = ".env.dev"
	}

	envPath := filepath.Join(utils.GetProjectRootPath(), envFileName)
	err := godotenv.Load(envPath)
	if err != nil {
		panic(fmt.Sprintf("Environment file not found: file_name=%s, path=%s, error=%v", envFileName, envPath, err))
	}

	config := &Config{
		AppConfig: AppConfig{
			Host:                       getEnvWithPanic(c.EnvKeyHost),
			Port:                       getEnvWithPanic(c.EnvKeyPort),
			DomainURL:                  getEnv(c.EnvKeyDomainURL, ""),
			UuidType:                   getEnv(c.EnvKeyUuidType, ""),
			WriteLogWhen:               getEnv(c.EnvKeyWriteLogWhen, ""),
			DBStatsLogType:             getEnv(c.EnvKeyDBStatsLogType, ""),
			SuperAdminEmail:            getEnvWithPanic(c.EnvSuperAdminEmail),
			SuperAdminPassword:         getEnvWithPanic(c.EnvSuperAdminPassword),
			SocialLoginDefaultPassword: getEnvWithPanic(c.EnvKeySocialLoginDefaultPassword),
		},
		DBConfig: DBConfig{
			DBHost:           getEnvWithPanic(c.EnvKeyDBHost),
			DBPort:           getEnvWithPanic(c.EnvKeyDBPort),
			DBUser:           getEnvWithPanic(c.EnvKeyDBUser),
			DBPassword:       getEnvWithPanic(c.EnvKeyDBPassword),
			DBName:           getEnvWithPanic(c.EnvKeyDBName),
			DBSSLMode:        getEnvWithPanic(c.EnvKeyDBSSLMode),
			DBShowConsoleLog: getEnvWithPanic(c.EnvKeyDBShowConsoleLog),
		},
		SecretConfig: SecretConfig{
			JWTSecret:           getEnvWithPanic(c.EnvKeyJWTSecret),
			AppHealthCheckToken: getEnvWithPanic(c.EnvKeyAppHealthCheckToken),
		},
		ExtServiceConfig: ExternalServiceConfig{
			FcmServerKey:   getEnv(c.EnvKeyFcmServerKey, ""),
			FcmUrl:         getEnv(c.EnvKeyFcmUrl, ""),
			GoogleClientID: getEnv(c.EnvKeyGoogleClientID, ""),
			FacebookAppID:  getEnv(c.EnvKeyFacebookAppID, ""),
		},
		ObjStoreConfig: ObjectStoreConfig{
			Region:          getEnvWithPanic(c.EnvKeyObjStoreRegion),
			BucketName:      getEnvWithPanic(c.EnvKeyObjStoreBucketName),
			AccountID:       getEnvWithPanic(c.EnvKeyObjStoreAccountID),
			AccessKeyID:     getEnvWithPanic(c.EnvKeyObjStoreAccessKeyID),
			AccessKeySecret: getEnvWithPanic(c.EnvKeyObjStoreAccessKeySecret),
			PublicDomain:    getEnvWithPanic(c.EnvKeyObjStorePublicDomain),
		},
	}

	return config
}

func GetStorageDomain() string {
	storageDomain := GetConfig().ObjStoreConfig.PublicDomain
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
