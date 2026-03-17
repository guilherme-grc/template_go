package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppPort              string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	JWTSecret            string
	JWTAccessExpiryMin   int
	JWTRefreshExpiryDays int
}

// Load reads configuration from environment variables, with optional .env fallback.
// Aborts the process if any required variable is missing or malformed.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env file not found, using system environment variables")
	}

	// [C] #1 — JWT_SECRET is mandatory; no fallback to avoid predictable secrets in production
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		slog.Error("JWT_SECRET is required and has not been defined")
		os.Exit(1)
	}

	// [M] #2 — Explicitly handle conversion errors
	accessExpiry, err := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRY_MINUTES", "15"))
	if err != nil {
		slog.Error("invalid JWT_ACCESS_EXPIRY_MINUTES", "error", err)
		os.Exit(1)
	}

	refreshExpiry, err := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_DAYS", "7"))
	if err != nil {
		slog.Error("invalid JWT_REFRESH_EXPIRY_DAYS", "error", err)
		os.Exit(1)
	}

	cfg := &Config{
		AppPort:              getEnv("APP_PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "template_db"),
		JWTSecret:            jwtSecret,
		JWTAccessExpiryMin:   accessExpiry,
		JWTRefreshExpiryDays: refreshExpiry,
	}

	if err := cfg.validate(); err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	return cfg
}

// validate checks that all required fields are present.
func (c *Config) validate() error {
	var errs []error
	if c.DBPassword == "" {
		errs = append(errs, errors.New("DB_PASSWORD is required"))
	}
	if c.DBUser == "" {
		errs = append(errs, errors.New("DB_USER is required"))
	}
	return errors.Join(errs...)
}

// getEnv returns the value of key from the environment, or fallback if not set.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
