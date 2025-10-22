package validation

import (
	"github.com/go-playground/validator/v10"
)

// Validator is a global validator instance
var Validator = validator.New()

// ValidateStruct validates a struct using the global validator
func ValidateStruct(s interface{}) error {
	return Validator.Struct(s)
}

// ValidationError represents a validation error with field and message
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrorResponse represents the response format for validation errors
type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

// GetValidationErrors extracts validation errors into a structured format
func GetValidationErrors(err error) ValidationErrorResponse {
	var response ValidationErrorResponse

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			response.Errors = append(response.Errors, ValidationError{
				Field:   fieldError.Field(),
				Message: getValidationMessage(fieldError),
			})
		}
	}

	return response
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "Это поле обязательно для заполнения"
	case "min":
		return "Значение слишком короткое"
	case "max":
		return "Значение слишком длинное"
	case "email":
		return "Неверный формат email"
	case "e164":
		return "Телефон должен быть в международном формате (например, +79991234567)"
	case "oneof":
		return "Недопустимое значение"
	case "numeric":
		return "Должно быть числом"
	case "gte":
		return "Значение должно быть больше или равно " + fieldError.Param()
	case "lte":
		return "Значение должно быть меньше или равно " + fieldError.Param()
	default:
		return "Недопустимое значение"
	}
}
