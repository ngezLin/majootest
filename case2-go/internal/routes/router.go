package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"majootest/case2-go/internal/config"
	"majootest/case2-go/internal/controllers"
	"majootest/case2-go/internal/middlewares"
	"majootest/case2-go/internal/utils"
)

// RegisterRoutes attaches the application's HTTP routes to the Gin engine.
func RegisterRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authController *controllers.AuthController,
	postController *controllers.PostController,
	commentController *controllers.CommentController,
) {
	r.GET("/health", func(c *gin.Context) {
		utils.SendSuccess(c, http.StatusOK, "API is running smoothly", gin.H{
			"status": "healthy",
			"time":   time.Now().UTC().Format(time.RFC3339),
		}, nil)
	})

	v1 := r.Group("/api/v1")
	auth := v1.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)
	auth.GET("/me", middlewares.AuthMiddleware(cfg), authController.GetProfile)

	posts := v1.Group("/posts")
	posts.GET("", postController.List)
	posts.GET("/:id", postController.GetByID)
	posts.GET("/:id/comments", commentController.ListByPostID)
	posts.GET("/:id/comments/:commentId", commentController.GetByID)

	protectedPosts := posts.Group("")
	protectedPosts.Use(middlewares.AuthMiddleware(cfg))
	protectedPosts.POST("", postController.Create)
	protectedPosts.PUT("/:id", postController.Update)
	protectedPosts.DELETE("/:id", postController.Delete)
	protectedPosts.POST("/:id/comments", commentController.Create)
	protectedPosts.PUT("/:id/comments/:commentId", commentController.Update)
	protectedPosts.DELETE("/:id/comments/:commentId", commentController.Delete)
}
