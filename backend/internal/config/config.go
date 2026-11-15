package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string

	CORSAllowedOrigins []string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using environment variables")
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8090"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "soilcycle"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "soilcycle"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", ""),

		CORSAllowedOrigins: getOrigins(
			getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getOrigins(value string) []string {
	var origins []string

	for _, origin := range strings.Split(value, ",") {
		origin = strings.TrimSpace(origin)

		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}
