// Package config loads runtime configuration from environment variables.
// Secrets are never committed; see docs/02-system/deployment.md §4.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration for the backend.
type Config struct {
	AppEnv string // "dev" | "uat" | "prod"
	Port   string

	// Database: pooled URL for the app runtime (see deployment.md §3 — must be the
	// Neon -pooler connection string on free tier). Migrations use the direct URL.
	DatabaseURL       string
	DatabaseURLDirect string

	// Pool sizing kept small for Neon free-tier connection limits.
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	// Auth
	JWTSecret string
	JWTTTL    time.Duration
}

// Load reads configuration from the environment, applying sane defaults for local dev.
func Load() (*Config, error) {
	c := &Config{
		AppEnv:            getenv("APP_ENV", "dev"),
		Port:              getenv("PORT", "8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		DatabaseURLDirect: getenv("DATABASE_URL_DIRECT", os.Getenv("DATABASE_URL")),
		DBMaxOpenConns:    getenvInt("DB_MAX_OPEN_CONNS", 8),
		DBMaxIdleConns:    getenvInt("DB_MAX_IDLE_CONNS", 2),
		DBConnMaxLifetime: time.Duration(getenvInt("DB_CONN_MAX_LIFETIME_MIN", 5)) * time.Minute,
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTTTL:            time.Duration(getenvInt("JWT_TTL_HOURS", 24)) * time.Hour,
	}

	if c.DatabaseURL == "" {
		return nil, fmt.Errorf("config: DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return nil, fmt.Errorf("config: JWT_SECRET is required")
	}
	return c, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}
