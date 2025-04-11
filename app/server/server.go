package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"GoMail/app/config"
	"GoMail/app/handler"
	"GoMail/app/handler/email"
	emailLogic "GoMail/app/logic/email"
	"GoMail/app/middleware"
	"GoMail/app/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Server represents the HTTP server
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
	config     *config.Config
	db         *mongo.Database
	repo       repository.Repository
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	router := gin.Default()

	// Setup middlewares
	corsConfig := configureCORS(cfg)
	router.Use(cors.New(corsConfig))
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.SecurityHeaders()) // Add security headers to all responses

	// Get MongoDB client from config and create database instance
	mongoClient := cfg.MongoDB.Client
	if mongoClient == nil {
		panic("MongoDB client not initialized. Make sure to call config.Init() first")
	}
	db := mongoClient.Database(cfg.MongoDB.Database)
	
	// Initialize repository
	repo := repository.New(&repository.DB{MongoDB: db})
	
	// Initialize email service
	emailService := emailLogic.NewEmailService(cfg, repo)
	
	// Create email handler
	emailHandler := email.NewHandler(emailService)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Setup routes with the emailHandler instance
	handler.InitPublicRoutes(router, emailHandler, repo)
	handler.InitProtectedRoutes(router, emailHandler)

	// Initialize token indexes for token revocation support
	if cfg.JWT.EnableTokenRevoking {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		if err := repo.InitTokenIndexes(ctx); err != nil {
			log.Printf("WARNING: Failed to initialize token indexes: %v", err)
		}
	}

	// Initialize the server
	server := &Server{
		router: router,
		config: cfg,
		db:     db,
		repo:   repo,
		httpServer: &http.Server{
			Addr:         ":" + cfg.Server.Port,
			Handler:      router,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
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