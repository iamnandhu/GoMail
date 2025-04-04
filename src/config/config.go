package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
	SMTP     SMTPConfig     `yaml:"smtp"`
	JWT      JWTConfig      `yaml:"jwt"`
	LogLevel string         `yaml:"logLevel"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

// MongoDBConfig holds MongoDB connection details
type MongoDBConfig struct {
	URI               string        `yaml:"uri"`
	Database          string        `yaml:"database"`
	Timeout           time.Duration `yaml:"timeout"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
	Client            *mongo.Client `yaml:"-"`
}

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	From      string `yaml:"from"`
	TLSEnable bool   `yaml:"tlsEnable"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret    string        `yaml:"secret"`
	ExpiresIn time.Duration `yaml:"expiresIn"`
}

// Load loads configuration from environment variables and config file
func Load() (*Config, error) {
	// Initialize default config
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		MongoDB: MongoDBConfig{
			URI:               getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:          getEnv("MONGODB_DB", "gomail"),
			Timeout:           10 * time.Second,
			ConnectionTimeout: 10 * time.Second,
		},
		SMTP: SMTPConfig{
			Host:      getEnv("SMTP_HOST", "smtp.example.com"),
			Port:      getEnv("SMTP_PORT", "587"),
			Username:  getEnv("SMTP_USERNAME", ""),
			Password:  getEnv("SMTP_PASSWORD", ""),
			From:      getEnv("SMTP_FROM", "no-reply@example.com"),
			TLSEnable: true,
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "default_jwt_secret_change_in_production"),
			ExpiresIn: 24 * time.Hour,
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Try to load and override with YAML config if available
	if err := loadYAMLConfig(cfg); err != nil {
		// Log warning but continue with environment variables
		fmt.Printf("Warning: Could not load config.yaml: %v\n", err)
	}

	// Override with environment variables
	overrideWithEnv(cfg)

	// Initialize MongoDB connection
	if err := initMongoDBConnection(cfg); err != nil {
		return nil, fmt.Errorf("failed to initialize MongoDB connection: %w", err)
	}

	return cfg, nil
}

// loadYAMLConfig loads configuration from config.yaml if available
func loadYAMLConfig(cfg *Config) error {
	configPath := getEnv("CONFIG_PATH", "config.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	
	return yaml.Unmarshal(data, cfg)
}

// overrideWithEnv overrides config with environment variables
func overrideWithEnv(cfg *Config) {
	// Server config
	if port := os.Getenv("PORT"); port != "" {
		cfg.Server.Port = port
	}
	
	// MongoDB config
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		cfg.MongoDB.URI = uri
	}
	if db := os.Getenv("MONGODB_DB"); db != "" {
		cfg.MongoDB.Database = db
	}
	
	// SMTP config
	if host := os.Getenv("SMTP_HOST"); host != "" {
		cfg.SMTP.Host = host
	}
	if port := os.Getenv("SMTP_PORT"); port != "" {
		cfg.SMTP.Port = port
	}
	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		cfg.SMTP.Username = username
	}
	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		cfg.SMTP.Password = password
	}
	if from := os.Getenv("SMTP_FROM"); from != "" {
		cfg.SMTP.From = from
	}
	
	// JWT config
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		cfg.JWT.Secret = secret
	}
}

// initMongoDBConnection initializes the MongoDB connection
func initMongoDBConnection(cfg *Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.ConnectionTimeout)
	defer cancel()
	
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	
	// Ping the database to verify connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer pingCancel()
	
	if err := client.Ping(pingCtx, nil); err != nil {
		return err
	}
	
	cfg.MongoDB.Client = client
	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Disconnect closes the MongoDB connection
func (cfg *Config) Disconnect(ctx context.Context) error {
	if cfg.MongoDB.Client != nil {
		return cfg.MongoDB.Client.Disconnect(ctx)
	}
	return nil
} 