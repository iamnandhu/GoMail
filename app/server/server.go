package server

import (
	"context"
	"net/http"
	"time"

	"GoMail/app/config"
	"GoMail/app/handler"
	"GoMail/app/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *config.Config
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	router := gin.Default()

	// Setup middlewares
	corsConfig := configureCORS(cfg)
	router.Use(cors.New(corsConfig))
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Setup routes
	mw := middleware.New()
	handler.InitPublicRoutes(router)
	handler.InitProtectedRoutes(router, mw)

	// Initialize the server
	server := &Server{
		router: router,
		config: cfg,
		httpServer: &http.Server{
			Addr:    ":" + cfg.Server.Port,
			Handler: router,
		},
	}

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

// Run starts the HTTP server
func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}