package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr                   string
	DatabasePath           string
	JWTSecret              string
	JWTExpiry              time.Duration
	VerificationCodeTTL    time.Duration
	SkipEmailVerification  bool
}

func Load() Config {
	jwtExp := 168 * time.Hour // 7d
	if s := os.Getenv("JWT_EXPIRY_HOURS"); s != "" {
		if h, err := strconv.Atoi(s); err == nil && h > 0 {
			jwtExp = time.Duration(h) * time.Hour
		}
	}
	codeTTL := 15 * time.Minute
	if s := os.Getenv("VERIFICATION_CODE_TTL_MINUTES"); s != "" {
		if m, err := strconv.Atoi(s); err == nil && m > 0 {
			codeTTL = time.Duration(m) * time.Minute
		}
	}
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "account-cabinet.db"
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-insecure-change-me"
	}
	return Config{
		Addr:                  addr,
		DatabasePath:          dbPath,
		JWTSecret:             secret,
		JWTExpiry:             jwtExp,
		VerificationCodeTTL:   codeTTL,
		SkipEmailVerification: os.Getenv("SKIP_EMAIL_VERIFICATION") == "true",
	}
}
