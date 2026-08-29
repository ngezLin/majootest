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

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Majoo Blog API - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
    .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

// RegisterRoutes attaches the application's HTTP routes to the Gin engine.
func RegisterRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authController *controllers.AuthController,
	postController *controllers.PostController,
	commentController *controllers.CommentController,
) {
	// Swagger UI & OpenAPI Specification routes
	r.StaticFile("/openapi.yaml", "./docs/openapi.yaml")
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerHTML)
	})
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerHTML)
	})

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
