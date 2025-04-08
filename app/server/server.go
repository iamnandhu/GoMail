package server

import (
	"context"
	"net/http"
	"time"

	"GoMail/app/config"
	"GoMail/app/handler"

	// "GoMail/app/logic"
	"GoMail/app/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *config.Config
	handler    *handler.Handler
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	router := gin.Default()

	// Create services
	// TODO: Uncomment and implement EmailService
	// emailService := logic.NewEmailService(cfg)

	// Create handlers
	// TODO: Update handler initialization with appropriate services
	h := &handler.Handler{} // Temporary placeholder

	// Initialize the server
	server := &Server{
		router: router,
		config: cfg,
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
		handler: h,
	}

	// Setup middlewares
	server.setupMiddlewares()

	// Setup routes
	server.setupRoutes()

	return server
}

// configureCORS returns CORS configuration
func configureCORS(cfg *config.Config) cors.Config {
	corsConfig := cors.DefaultConfig()
	
	// Use the CORS configuration from config if available
	if len(cfg.Cors.AllowedOrigins) > 0 {
		corsConfig.AllowOrigins = cfg.Cors.AllowedOrigins
	} else {
		corsConfig.AllowAllOrigins = true
	}
	
	corsConfig.AllowCredentials = true
	
	if len(cfg.Cors.AllowedHeaders) > 0 {
		corsConfig.AllowHeaders = cfg.Cors.AllowedHeaders
	} else {
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
	}
	
	if len(cfg.Cors.AllowedMethods) > 0 {
		corsConfig.AllowMethods = cfg.Cors.AllowedMethods
	} else {
		corsConfig.AddAllowMethods("POST", "OPTIONS", "GET", "PUT", "DELETE")
	}
	
	if len(cfg.Cors.ExposeHeaders) > 0 {
		corsConfig.ExposeHeaders = cfg.Cors.ExposeHeaders
	}
	
	if cfg.Cors.MaxAge > 0 {
		corsConfig.MaxAge = time.Duration(cfg.Cors.MaxAge) * time.Second
	}

	return corsConfig
}

// setupMiddlewares adds middleware to the router
func (s *Server) setupMiddlewares() {
	// Add CORS middleware
	s.router.Use(cors.New(configureCORS(s.config)))

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
		s.handler.RegisterRoutes(api)
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