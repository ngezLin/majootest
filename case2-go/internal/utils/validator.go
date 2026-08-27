package utils

import (
	"fmt"
	"majootest/case2-go/internal/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct based on its struct tags and returns detailed error information.
func ValidateStruct(s interface{}) ([]models.ErrorDetail, error) {
	err := validate.Struct(s)
	if err == nil {
		return nil, nil
	}

	var details []models.ErrorDetail
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, errField := range validationErrors {
			var message string
			switch errField.Tag() {
			case "required":
				message = fmt.Sprintf("Field '%s' is required", errField.Field())
			case "email":
				message = fmt.Sprintf("Field '%s' must be a valid email address", errField.Field())
			case "min":
				message = fmt.Sprintf("Field '%s' must be at least %s characters long", errField.Field(), errField.Param())
			case "max":
				message = fmt.Sprintf("Field '%s' cannot exceed %s characters", errField.Field(), errField.Param())
			default:
				message = fmt.Sprintf("Field '%s' failed validation for tag '%s'", errField.Field(), errField.Tag())
			}

			details = append(details, models.ErrorDetail{
				Field:   errField.Field(),
				Message: message,
			})
		}
	}

	return details, err
}
