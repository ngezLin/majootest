package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"majootest/case2-go/internal/config"
	"majootest/case2-go/internal/utils"
)

const (
	UserIDKey = "userID"
	EmailKey  = "email"
)

// AuthMiddleware creates a JWT authentication middleware for Gin.
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing Authorization header", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid Authorization header format. Expected 'Bearer <token>'", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			utils.SendError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", nil)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(EmailKey, claims.Email)
		c.Next()
	}
}

// GetUserIDFromGinContext retrieves the authenticated user's ID from the Gin context.
func GetUserIDFromGinContext(c *gin.Context) (int64, bool) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := val.(int64)
	return id, ok
}
