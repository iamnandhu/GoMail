package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"GoMail/app/config"
	"GoMail/app/server"
)

func main() {
	// Initialize configuration
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Setup server
	srv := server.New(cfg)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s (environment: %s)", cfg.Server.Port, cfg.Env)
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Disconnect MongoDB client
	if err := cfg.Disconnect(ctx); err != nil {
		log.Printf("Warning: Error disconnecting from MongoDB: %v", err)
	}

	log.Println("Server exiting")
} 