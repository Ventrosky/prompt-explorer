package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	ServerPort   string
	DatabaseURL  string
	Environment  string
	GeminiAPIKey string
}

// Load loads configuration from environment variables and .env file
func Load() *Config {
	// Try to load .env file (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or could not be loaded: %v", err)
	}

	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8081"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/prompt_explorer?sslmode=disable"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
	}
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// HasGeminiAPIKey returns true if Gemini API key is configured
func (c *Config) HasGeminiAPIKey() bool {
	return c.GeminiAPIKey != ""
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
