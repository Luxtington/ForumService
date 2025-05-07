package config

import (
	_ "log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL          string
	JWTSecret      string
	Port           int
	AuthServiceURL string
}

func Load() *Config {
	// Загружаем .env файл (если есть)
	_ = godotenv.Load()

	port, _ := strconv.Atoi(getEnv("PORT", "8080"))

	return &Config{
		DBURL:          getEnv("DB_URL", "postgres://user:pass@localhost:5432/forum?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "default-secret-key"),
		Port:           port,
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
