package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Host string
}

func LoadConfig() *Config {
	godotenv.Load()

	return &Config{
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
