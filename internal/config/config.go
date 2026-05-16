package config

import (
	"log"
	"os"
	"time"
)

const (
	MigrationsDir = "migrations"
	JWTExpiry     = 168 * time.Hour
)

type Config struct {
	Addr               string
	DatabaseURL        string
	JWTSecret          string
	BackendURL         string
	BackendInternalKey string
}

func Load() Config {
	return Config{
		Addr:               ":8081",
		DatabaseURL:        getEnv("DATABASE_URL"),
		JWTSecret:          getEnv("JWT_SECRET"),
		BackendURL:         getEnv("BACKEND_URL"),
		BackendInternalKey: getEnv("BACKEND_INTERNAL_KEY"),
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("переменная окружения %s не найдена", key)
	}
	return value
}
