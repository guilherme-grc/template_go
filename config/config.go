package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

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

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	accessExpiry, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRY_MINUTES", "15"))
	refreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_DAYS", "7"))

	return &Config{
		AppPort:              getEnv("APP_PORT", "8080"),
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres"),
		DBName:               getEnv("DB_NAME", "reembolso"),
		JWTSecret:            getEnv("JWT_SECRET", "changeme"),
		JWTAccessExpiryMin:   accessExpiry,
		JWTRefreshExpiryDays: refreshExpiry,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
