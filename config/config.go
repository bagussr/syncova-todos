package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port              string
	Host              string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	ENV               string
	AuthServiceURL    string
	AuthServiceAPIKey string
	JwtSecret         string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		Port:              getEnv("PORT", "8080"),
		Host:              getEnv("HOST", "localhost"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "user"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		DBName:            getEnv("DB_NAME", "dbname"),
		ENV:               getEnv("ENV", "development"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", ""),
		AuthServiceAPIKey: getEnv("AUTH_SERVICE_API_KEY", ""),
		JwtSecret:         getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
