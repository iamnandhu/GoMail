package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

const (
	configFile = "config.yaml"
	envKey     = "GOMAIL_ENV" // env variable name to get deployment environment
)

// Config holds all configuration for the application
type Config struct {
	Env      string         `yaml:"env" json:"env"`
	Server   ServerConfig   `yaml:"server" json:"server"`
	MongoDB  MongoDBConfig  `yaml:"mongodb" json:"mongodb"`
	SMTP     SMTPConfig     `yaml:"smtp" json:"smtp"`
	JWT      JWTConfig      `yaml:"jwt" json:"jwt"`
	LogLevel string         `yaml:"logLevel" json:"logLevel"`
	Cors     CorsConfig     `yaml:"cors" json:"cors"`
	Services ServiceConfigs `yaml:"services" json:"services"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string        `yaml:"port" json:"port"`
	PortGRPC     string        `yaml:"portGRPC" json:"portGRPC"`
	ReadTimeout  time.Duration `yaml:"readTimeout" json:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout" json:"writeTimeout"`
}

// MongoDBConfig holds MongoDB connection details
type MongoDBConfig struct {
	URI               string        `yaml:"uri" json:"uri"`
	Username          string        `yaml:"username" json:"username"`
	Password          string        `yaml:"password" json:"password"`
	Database          string        `yaml:"database" json:"database"`
	Endpoint          string        `yaml:"endpoint" json:"endpoint"`
	Timeout           time.Duration `yaml:"timeout" json:"timeout"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout" json:"connectionTimeout"`
	Client            *mongo.Client `yaml:"-" json:"-"`
}

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	Host          string `yaml:"host" json:"host"`
	Port          string `yaml:"port" json:"port"`
	Username      string `yaml:"username" json:"username"`
	Password      string `yaml:"password" json:"password"`
	From          string `yaml:"from" json:"from"`
	UseStartTLS   bool   `yaml:"useStartTLS" json:"useStartTLS"`
	MaxConcurrent int    `yaml:"maxConcurrent" json:"maxConcurrent"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	Secret              string        `yaml:"secret" json:"secret"`
	ExpiresIn           time.Duration `yaml:"expiresIn" json:"expiresIn"`
	EnableTokenRevoking bool          `yaml:"enableTokenRevoking" json:"enableTokenRevoking"`
	RevokedTokensTTL    time.Duration `yaml:"revokedTokensTTL" json:"revokedTokensTTL"`
}

// CorsConfig holds CORS configuration
type CorsConfig struct {
	AllowedOrigins []string `yaml:"allowedOrigins" json:"allowedOrigins"`
	AllowedMethods []string `yaml:"allowedMethods" json:"allowedMethods"`
	AllowedHeaders []string `yaml:"allowedHeaders" json:"allowedHeaders"`
	ExposeHeaders  []string `yaml:"exposeHeaders" json:"exposeHeaders"`
	MaxAge         int      `yaml:"maxAge" json:"maxAge"`
}

// ServiceConfigs holds configuration for external services
type ServiceConfigs map[string]ServiceConfig

// ServiceConfig holds configuration for a specific service
type ServiceConfig struct {
	URL  string `yaml:"url" json:"url"`
	Port string `yaml:"port" json:"port"`
}

var config *Config

// Init initializes the configuration
func Init() (*Config, error) {
	if err := loadConfig(); err != nil {
		return nil, err
	}

	env := os.Getenv(envKey)
	if env != "" {
		config.Env = env
	}

	// Check if we're in a production or staging environment
	if isProduction() || isStaging() {
		if err := overwriteConfigFromEnv(); err != nil {
			return nil, err
		}
	} else {
		// Always apply environment variables if they exist
		if err := overwriteConfigFromEnv(); err != nil {
			return nil, err
		}
	}

	// Initialize MongoDB connection
	if err := initMongoDB(); err != nil {
		return nil, err
	}
	
	// Set default values for token revocation if not set
	if config.JWT.RevokedTokensTTL == 0 {
		config.JWT.RevokedTokensTTL = 24 * time.Hour
	}
	
	return config, nil
}

// loadConfig loads configuration from config.yaml file
func loadConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = configFile
	}
	
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	
	config = &Config{}
	if err := yaml.Unmarshal(yamlFile, config); err != nil {
		return err
	}
	
	return nil
}

// overwriteConfigFromEnv overwrites configuration with environment variables
func overwriteConfigFromEnv() error {
	// Server config
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	if portGRPC := os.Getenv("PORT_GRPC"); portGRPC != "" {
		config.Server.PortGRPC = portGRPC
	}
	
	// MongoDB config
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		config.MongoDB.URI = uri
	}
	if username := os.Getenv("MONGODB_USERNAME"); username != "" {
		config.MongoDB.Username = username
	}
	if password := os.Getenv("MONGODB_PASSWORD"); password != "" {
		config.MongoDB.Password = password
	}
	if db := os.Getenv("MONGODB_DATABASE"); db != "" {
		config.MongoDB.Database = db
	}
	if endpoint := os.Getenv("MONGODB_ENDPOINT"); endpoint != "" {
		config.MongoDB.Endpoint = endpoint
	}
	
	// SMTP config
	if host := os.Getenv("SMTP_HOST"); host != "" {
		config.SMTP.Host = host
	}
	if port := os.Getenv("SMTP_PORT"); port != "" {
		config.SMTP.Port = port
	}
	if username := os.Getenv("SMTP_USERNAME"); username != "" {
		config.SMTP.Username = username
	}
	if password := os.Getenv("SMTP_PASSWORD"); password != "" {
		config.SMTP.Password = password
	}
	if from := os.Getenv("SMTP_FROM"); from != "" {
		config.SMTP.From = from
	}
	if useStartTLSStr := os.Getenv("SMTP_USE_STARTTLS"); useStartTLSStr != "" {
		config.SMTP.UseStartTLS = useStartTLSStr == "true" || useStartTLSStr == "1" || useStartTLSStr == "yes"
	}
	if maxConcurrentStr := os.Getenv("SMTP_MAX_CONCURRENT"); maxConcurrentStr != "" {
		if maxConcurrent, err := strconv.Atoi(maxConcurrentStr); err == nil {
			config.SMTP.MaxConcurrent = maxConcurrent
		}
	}
	
	// JWT config
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.JWT.Secret = secret
	}
	if enableTokenRevokingStr := os.Getenv("JWT_ENABLE_TOKEN_REVOKING"); enableTokenRevokingStr != "" {
		config.JWT.EnableTokenRevoking = enableTokenRevokingStr == "true" || enableTokenRevokingStr == "1" || enableTokenRevokingStr == "yes"
	}

	// Load JSON configuration from GOMAIL_CONFIG env var if it exists
	// This allows passing complex configuration as a single JSON string
	if configJSON := os.Getenv("GOMAIL_CONFIG"); configJSON != "" {
		if err := json.Unmarshal([]byte(configJSON), config); err != nil {
			return fmt.Errorf("failed to parse GOMAIL_CONFIG environment variable: %w", err)
		}
	}
	
	return nil
}

// initMongoDB initializes the MongoDB connection
func initMongoDB() error {
	// If URI is provided, use it directly
	uri := config.MongoDB.URI
	
	// Otherwise, construct URI from components
	if uri == "" && config.MongoDB.Endpoint != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s/%s",
			config.MongoDB.Username,
			config.MongoDB.Password,
			config.MongoDB.Endpoint,
			config.MongoDB.Database)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), config.MongoDB.ConnectionTimeout)
	defer cancel()
	
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	
	// Ping the database to verify connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), config.MongoDB.Timeout)
	defer pingCancel()
	
	if err := client.Ping(pingCtx, nil); err != nil {
		return err
	}
	
	config.MongoDB.Client = client
	return nil
}

// Get returns the current configuration
func Get() *Config {
	if config == nil {
		panic("Configuration not loaded. Call Init() first!")
	}
	return config
}

// GetServiceConfig returns configuration for a specific service
func GetServiceConfig(identifier string) ServiceConfig {
	return config.Services[identifier]
}

// Disconnect closes the MongoDB connection
func (cfg *Config) Disconnect(ctx context.Context) error {
	if cfg.MongoDB.Client != nil {
		return cfg.MongoDB.Client.Disconnect(ctx)
	}
	return nil
}

// Environment helper functions
func isProduction() bool {
	return config.Env == "prod" || config.Env == "production"
}

func isStaging() bool {
	return config.Env == "stg" || config.Env == "staging"
}