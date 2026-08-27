package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"majootest/case2-go/internal/middlewares"
	"majootest/case2-go/internal/models"
	"majootest/case2-go/internal/services"
	"majootest/case2-go/internal/utils"
)

// PostController handles post-related HTTP endpoints using Gin.
type PostController struct {
	postService services.PostService
}

// NewPostController creates a new PostController instance.
func NewPostController(postService services.PostService) *PostController {
	return &PostController{postService: postService}
}

// Create handles the creation of a new blog post.
func (ctrl *PostController) Create(c *gin.Context) {
	userID, ok := middlewares.GetUserIDFromGinContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body", nil)
		return
	}

	if details, err := utils.ValidateStruct(req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
		return
	}

	post, err := ctrl.postService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "POST_CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Post created successfully", post, nil)
}

// GetByID handles retrieving a single blog post by its ID.
func (ctrl *PostController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Post ID must be a valid integer", nil)
		return
	}

	post, err := ctrl.postService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "POST_NOT_FOUND", "Post not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "POST_FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Post retrieved successfully", post, nil)
}

// List handles retrieving a paginated list of posts with optional search filter.
func (ctrl *PostController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	posts, meta, err := ctrl.postService.List(c.Request.Context(), page, limit, search)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "POSTS_FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Posts retrieved successfully", posts, meta)
}

// Update handles updating an existing blog post.
func (ctrl *PostController) Update(c *gin.Context) {
	userID, ok := middlewares.GetUserIDFromGinContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Post ID must be a valid integer", nil)
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body", nil)
		return
	}

	if details, err := utils.ValidateStruct(req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
		return
	}

	post, err := ctrl.postService.Update(c.Request.Context(), id, userID, &req)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "POST_NOT_FOUND", "Post not found", nil)
			return
		}
		if errors.Is(err, models.ErrForbidden) {
			utils.SendError(c, http.StatusForbidden, "FORBIDDEN", "You are not authorized to edit this post", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "POST_UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Post updated successfully", post, nil)
}

// Delete handles deleting a blog post and its comments atomically.
func (ctrl *PostController) Delete(c *gin.Context) {
	userID, ok := middlewares.GetUserIDFromGinContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Post ID must be a valid integer", nil)
		return
	}

	if err := ctrl.postService.Delete(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "POST_NOT_FOUND", "Post not found", nil)
			return
		}
		if errors.Is(err, models.ErrForbidden) {
			utils.SendError(c, http.StatusForbidden, "FORBIDDEN", "You are not authorized to delete this post", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "POST_DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Post deleted successfully", nil, nil)
}
