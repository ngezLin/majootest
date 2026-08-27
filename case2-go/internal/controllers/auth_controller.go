package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"majootest/case2-go/internal/middlewares"
	"majootest/case2-go/internal/models"
	"majootest/case2-go/internal/services"
	"majootest/case2-go/internal/utils"
)

// AuthController handles authentication HTTP endpoints using Gin.
type AuthController struct {
	authService services.AuthService
}

// NewAuthController creates a new AuthController instance.
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register handles user registration.
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body", nil)
		return
	}

	if details, err := utils.ValidateStruct(req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
		return
	}

	resp, err := ctrl.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			utils.SendError(c, http.StatusConflict, "USER_ALREADY_EXISTS", err.Error(), nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "REGISTRATION_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "User registered successfully", resp, nil)
}

// Login handles user authentication and JWT issuance.
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON request body", nil)
		return
	}

	if details, err := utils.ValidateStruct(req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
		return
	}

	resp, err := ctrl.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			utils.SendError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "LOGIN_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Login successful", resp, nil)
}

// GetProfile returns the authenticated user's profile.
func (ctrl *AuthController) GetProfile(c *gin.Context) {
	userID, ok := middlewares.GetUserIDFromGinContext(c)
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated", nil)
		return
	}

	user, err := ctrl.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			utils.SendError(c, http.StatusNotFound, "USER_NOT_FOUND", "User profile not found", nil)
			return
		}
		utils.SendError(c, http.StatusInternalServerError, "PROFILE_FETCH_FAILED", err.Error(), nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Profile retrieved successfully", user, nil)
}
