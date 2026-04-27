package config

import (
	"os"
	"time"
)

const (
	DefaultAddr               = ":8081"
	DefaultDatabaseURL        = "postgres://focus:focus@127.0.0.1:5432/focus_account_cabinet?sslmode=disable"
	DefaultJWTSecret          = "dev-insecure-change-me"
	DefaultBackendURL         = "http://localhost:8082"
	DefaultBackendInternalKey = "backend-internal-dev-key"
	MigrationsDir             = "migrations"
	JWTExpiry                 = 168 * time.Hour
	BackendTimeout            = 5 * time.Second
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
		Addr:               envOrDefault("ADDR", DefaultAddr),
		DatabaseURL:        envOrDefault("DATABASE_URL", DefaultDatabaseURL),
		JWTSecret:          envOrDefault("JWT_SECRET", DefaultJWTSecret),
		BackendURL:         envOrDefault("BACKEND_URL", DefaultBackendURL),
		BackendInternalKey: envOrDefault("BACKEND_INTERNAL_KEY", DefaultBackendInternalKey),
	}
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
