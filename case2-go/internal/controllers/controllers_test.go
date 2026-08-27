package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"majootest/case2-go/internal/config"
	"majootest/case2-go/internal/middlewares"
	"majootest/case2-go/internal/models"
	"majootest/case2-go/internal/repositories"
	"majootest/case2-go/internal/services"
)

func setupTestRouter() (*gin.Engine, string, int64) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		JWTSecret:      "test-secret-key-12345",
		JWTExpiryHours: 24,
	}

	userRepo := repositories.NewMockUserRepository()
	postRepo := repositories.NewMockPostRepository()
	commentRepo := repositories.NewMockCommentRepository()
	txManager := &repositories.MockTxManager{}

	authSvc := services.NewAuthService(userRepo, cfg)
	postSvc := services.NewPostService(postRepo, commentRepo, txManager)
	commentSvc := services.NewCommentService(commentRepo, postRepo, txManager)

	authController := NewAuthController(authSvc)
	postController := NewPostController(postSvc)
	commentController := NewCommentController(commentSvc)

	r := gin.New()

	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authController.Register)
			authGroup.POST("/login", authController.Login)
			authGroup.GET("/me", middlewares.AuthMiddleware(cfg), authController.GetProfile)
		}

		postsGroup := v1.Group("/posts")
		{
			postsGroup.GET("", postController.List)
			postsGroup.GET("/:id", postController.GetByID)
			postsGroup.GET("/:id/comments", commentController.ListByPostID)
			postsGroup.GET("/:id/comments/:commentId", commentController.GetByID)

			postsProtected := postsGroup.Group("")
			postsProtected.Use(middlewares.AuthMiddleware(cfg))
			{
				postsProtected.POST("", postController.Create)
				postsProtected.PUT("/:id", postController.Update)
				postsProtected.DELETE("/:id", postController.Delete)

				postsProtected.POST("/:id/comments", commentController.Create)
				postsProtected.PUT("/:id/comments/:commentId", commentController.Update)
				postsProtected.DELETE("/:id/comments/:commentId", commentController.Delete)
			}
		}
	}

	// Pre-create user and token
	regResp, _ := authSvc.Register(context.Background(), &models.RegisterRequest{
		Name:     "Author User",
		Email:    "author@example.com",
		Password: "password123",
	})

	return r, regResp.Token, regResp.User.ID
}

func TestRegisterEndpoint(t *testing.T) {
	router, _, _ := setupTestRouter()

	body, _ := json.Marshal(models.RegisterRequest{
		Name:     "New User",
		Email:    "newuser@example.com",
		Password: "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	var resp models.SuccessResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil || !resp.Success {
		t.Errorf("expected success envelope, got %v", resp)
	}
}

func TestCreatePostWithoutAuth(t *testing.T) {
	router, _, _ := setupTestRouter()

	body, _ := json.Marshal(models.CreatePostRequest{
		Title:   "Unauthenticated Post",
		Content: "Should fail with 401 Unauthorized",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized, got %d", w.Code)
	}
}

func TestCreateAndGetPostWithAuth(t *testing.T) {
	router, token, _ := setupTestRouter()

	// 1. Create Post
	body, _ := json.Marshal(models.CreatePostRequest{
		Title:   "Integration Post Title",
		Content: "Detailed blog post content for test.",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d: %s", w.Code, w.Body.String())
	}

	// 2. Fetch Post List
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=1&limit=5", nil)
	listW := httptest.NewRecorder()

	router.ServeHTTP(listW, listReq)

	if listW.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d: %s", listW.Code, listW.Body.String())
	}
}

func TestValidationFailureEndpoint(t *testing.T) {
	router, _, _ := setupTestRouter()

	// Missing required fields
	body, _ := json.Marshal(models.RegisterRequest{
		Name:     "",
		Email:    "invalid-email",
		Password: "123",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", w.Code)
	}

	var errResp models.ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Success || len(errResp.Error.Details) == 0 {
		t.Errorf("expected validation error details, got: %v", errResp)
	}
}
