package server

import (
	"context"
	"net/http"

	"GoMail/src/config"
	"GoMail/src/handler"
	"GoMail/src/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router *gin.Engine
	httpServer *http.Server
	config *config.Config
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	router := gin.Default()
	
	// Initialize the server
	server := &Server{
		router: router,
		config: cfg,
		httpServer: &http.Server{
			Addr:    ":" + cfg.ServerPort,
			Handler: router,
		},
	}

	// Setup middlewares
	server.setupMiddlewares()
	
	// Setup routes
	server.setupRoutes()

	return server
}

// configureCORS returns CORS configuration
func configureCORS() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders(
		"Content-Type", 
		"Content-Length", 
		"Accept-Encoding", 
		"X-CSRF-Token", 
		"Authorization", 
		"accept", 
		"origin", 
		"Cache-Control", 
		"X-Requested-With",
	)
	corsConfig.AddAllowMethods("POST", "OPTIONS", "GET", "PUT", "DELETE")
	
	return corsConfig
}

// setupMiddlewares adds middleware to the router
func (s *Server) setupMiddlewares() {
	// Add CORS middleware
	s.router.Use(cors.New(configureCORS()))
	
	// Add other global middlewares
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recovery())
}

// setupRoutes configures all the routes
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Register handlers
		handler.RegisterRoutes(api)
	}
}

// Run starts the HTTP server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
} 