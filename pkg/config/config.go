package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv          string
	AppPort         string
	JWTSecret       string
	JWTExpireMinute int

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func Load() Config {
	expire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_MINUTES", "60"))
	return Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		AppPort:         getEnv("APP_PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpireMinute: expire,
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "halolight"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
	}
}
