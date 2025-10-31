package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DatabasePath  string
	Port          string
	SessionSecret string
	AdminUsername string
	AdminEmail    string
	AdminPassword string
}

var E EnvConfig

func LoadEnv() {
	// Load .env file if present (does not fail if missing)
	_ = godotenv.Load()

	E = EnvConfig{
		DatabasePath:  getEnv("DATABASE_PATH", "forum.db"),
		Port:          getEnv("PORT", "8080"),
		SessionSecret: getEnv("SESSION_SECRET", "dev-secret"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "changeme"),
	}

	log.Printf("Loaded .env: DB=%s PORT=%s", E.DatabasePath, E.Port)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
