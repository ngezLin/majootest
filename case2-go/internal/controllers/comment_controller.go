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

// CommentController handles comment-related HTTP endpoints using Gin.
type CommentController struct {
	commentService services.CommentService
}

// NewCommentController creates a new CommentController instance.
func NewCommentController(commentService services.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

// Create handles creating a new comment on a post.
func (ctrl *CommentController) Create(c *gin.Context) {
	userID, ok := middlewares.GetUserIDFromGinContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	postIDStr := c.Param("id")
	if postIDStr == "" {
		postIDStr = c.Param("postId")
	}
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Post ID must be a valid integer", nil)
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body", nil)
		return
	}

	if details, err := utils.ValidateStruct(req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
		return
	}

	comment, err := ctrl.commentService.Create(c.Request.Context(), postID, userID, &req)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "POST_NOT_FOUND", "Cannot comment on non-existent post", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "COMMENT_CREATE_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Comment created successfully", comment, nil)
}

// GetByID handles retrieving a single comment by ID.
func (ctrl *CommentController) GetByID(c *gin.Context) {
	idStr := c.Param("commentId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Comment ID must be a valid integer", nil)
		return
	}

	comment, err := ctrl.commentService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "COMMENT_NOT_FOUND", "Comment not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "COMMENT_FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Comment retrieved successfully", comment, nil)
}

// ListByPostID handles retrieving a paginated list of comments for a given post.
func (ctrl *CommentController) ListByPostID(c *gin.Context) {
	postIDStr := c.Param("id")
	if postIDStr == "" {
		postIDStr = c.Param("postId")
	}
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Post ID must be a valid integer", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	comments, meta, err := ctrl.commentService.ListByPostID(c.Request.Context(), postID, page, limit)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "POST_NOT_FOUND", "Post not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "COMMENTS_FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Comments retrieved successfully", comments, meta)
}

// Update handles updating an existing comment.
func (ctrl *CommentController) Update(c *gin.Context) {
	userID, ok := middlewares.GetUserIDFromGinContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	idStr := c.Param("commentId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Comment ID must be a valid integer", nil)
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body", nil)
		return
	}

	if details, err := utils.ValidateStruct(req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
		return
	}

	comment, err := ctrl.commentService.Update(c.Request.Context(), id, userID, &req)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "COMMENT_NOT_FOUND", "Comment not found", nil)
			return
		}
		if errors.Is(err, models.ErrForbidden) {
			utils.SendError(c, http.StatusForbidden, "FORBIDDEN", "You are not authorized to edit this comment", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "COMMENT_UPDATE_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Comment updated successfully", comment, nil)
}

// Delete handles deleting a comment and decrementing the post comment count.
func (ctrl *CommentController) Delete(c *gin.Context) {
	userID, ok := middlewares.GetUserIDFromGinContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	idStr := c.Param("commentId")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_PARAM", "Comment ID must be a valid integer", nil)
		return
	}

	if err := ctrl.commentService.Delete(c.Request.Context(), id, userID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "COMMENT_NOT_FOUND", "Comment not found", nil)
			return
		}
		if errors.Is(err, models.ErrForbidden) {
			utils.SendError(c, http.StatusForbidden, "FORBIDDEN", "You are not authorized to delete this comment", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "COMMENT_DELETE_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Comment deleted successfully", nil, nil)
}
