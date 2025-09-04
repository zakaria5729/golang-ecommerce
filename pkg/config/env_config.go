package config

import "os"

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
	return Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "ep-nameless-boat-a1ha9si8-pooler.ap-southeast-1.aws.neon.tech"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "neondb_owner"),
		DBPassword: getEnv("DB_PASSWORD", "npg_bUkBA2NOm3SC"),
		DBName:     getEnv("DB_NAME", "neondb"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "require"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
