package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName  string
	AppPort  string
	DBUser   string
	DBPass   string
	DBHost   string
	DBPort   string
	DBName   string
	DBUrl    string
	RedisUrl string
	LogLevel string
}

// Load loads configuration from a specified .env file
func Load(envFile string) *Config {
	// Try to load the specified .env file (e.g. ".env.auth")
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️ No %s file found, using system environment variables", envFile)
	}

	cfg := &Config{
		AppName:  getEnv("APP_NAME", "social-network"),
		AppPort:  getEnv("APP_PORT", "8000"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASSWORD", "postgres"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBName:   getEnv("DB_NAME", "socialnet"),
		DBUrl:    getEnv("DATABASE_URL", ""),
		RedisUrl: getEnv("REDIS_URL", "localhost:6379"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Build DBUrl if not given
	if cfg.DBUrl == "" {
		cfg.DBUrl = "postgres://" + cfg.DBUser + ":" + cfg.DBPass +
			"@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=disable"
	}

	log.Printf("✅ Config loaded for %s on port %s", cfg.AppName, cfg.AppPort)
	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
