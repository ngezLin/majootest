package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"majootest/case2-go/internal/config"
	"majootest/case2-go/internal/controllers"
	"majootest/case2-go/internal/repositories"
	"majootest/case2-go/internal/routes"
	"majootest/case2-go/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.JWTSecret == "" {
		log.Fatal("[FATAL] JWT_SECRET environment variable is required. Please set it in your .env file or environment.")
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Printf("Starting Blog API Server in %s mode...", cfg.Environment)

	// Database Connection (MySQL)
	db, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Connection Pool Settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping database
	if err := db.Ping(); err != nil {
		log.Printf("[WARNING] Database ping failed (ensure MySQL is running): %v", err)
	} else {
		log.Println("[INFO] Successfully connected to MySQL database")
	}

	// Initialize Repositories & Transaction Manager
	txManager := repositories.NewTxManager(db)
	userRepo := repositories.NewUserRepository(db)
	postRepo := repositories.NewPostRepository(db)
	commentRepo := repositories.NewCommentRepository(db)

	// Initialize Services
	authService := services.NewAuthService(userRepo, cfg)
	postService := services.NewPostService(postRepo, commentRepo, txManager)
	commentService := services.NewCommentService(commentRepo, postRepo, txManager)

	// Initialize Controllers
	authController := controllers.NewAuthController(authService)
	postController := controllers.NewPostController(postService)
	commentController := controllers.NewCommentController(commentService)

	// Initialize Gin Engine
	r := gin.Default()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(r, cfg, authController, postController, commentController)

	// Server Configuration & Graceful Shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background goroutine
	go func() {
		log.Printf("[INFO] Server listening on http://localhost:%s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	// Graceful shutdown listener
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server stopped gracefully")
}
