package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	ServerPort string
	MongoDB    MongoDBConfig
	SMTPConfig SMTPConfig
}

// MongoDBConfig holds MongoDB connection details
type MongoDBConfig struct {
	URI        string
	Database   string
	Collection string
}

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// Load loads configuration from environment variables or defaults
func Load() (*Config, error) {
	// Default config values
	cfg := &Config{
		ServerPort: getEnv("PORT", "8080"),
		MongoDB: MongoDBConfig{
			URI:        getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:   getEnv("MONGODB_DB", "gomail"),
			Collection: getEnv("MONGODB_COLLECTION", "emails"),
		},
		SMTPConfig: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.example.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "no-reply@example.com"),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 