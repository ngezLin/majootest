package utils

import (
	"github.com/gin-gonic/gin"
	"majootest/case2-go/internal/models"
)

// SendSuccess formats and sends a standardized JSON success response.
func SendSuccess(c *gin.Context, status int, message string, data interface{}, meta *models.PaginationMeta) {
	c.JSON(status, models.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// SendError formats and sends a standardized JSON error response.
func SendError(c *gin.Context, status int, code, message string, details []models.ErrorDetail) {
	c.JSON(status, models.ErrorResponse{
		Success: false,
		Error: models.ErrorPayload{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
