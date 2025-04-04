package config

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Set up test environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	os.Setenv("JWT_SECRET", "test_secret")
	
	// Point to the test config file
	os.Setenv("CONFIG_PATH", "../../config.yaml")
	
	// Load the configuration
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}
	
	// Make sure we clean up the MongoDB connection after the test
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := cfg.Disconnect(ctx); err != nil {
			t.Logf("Failed to disconnect MongoDB: %v", err)
		}
	}()
	
	// Verify that environment variables override config file values
	if cfg.Server.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got '%s'", cfg.Server.Port)
	}
	
	// Verify that config file values are loaded
	if cfg.SMTP.Host != "smtp.example.com" {
		t.Errorf("Expected SMTP Host to be 'smtp.example.com', got '%s'", cfg.SMTP.Host)
	}
	
	// Verify that JWT configuration is correct
	if cfg.JWT.Secret != "test_secret" {
		t.Errorf("Expected JWT Secret to be 'test_secret', got '%s'", cfg.JWT.Secret)
	}
	
	// Verify MongoDB connection was established
	if cfg.MongoDB.Client == nil {
		t.Error("MongoDB client was not initialized")
	}
} 