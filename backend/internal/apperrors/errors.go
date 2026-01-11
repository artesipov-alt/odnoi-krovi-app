package apperrors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ErrorCode представляет код ошибки для API
type ErrorCode string

const (
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrCodeBadRequest    ErrorCode = "BAD_REQUEST"
	ErrCodeConflict      ErrorCode = "CONFLICT"
)

// AppError представляет ошибку приложения с метаданными.
// Реализует huma.StatusError для интеграции с Huma.
type AppError struct {
	Code       ErrorCode      `json:"code"`
	Message    string         `json:"message"`
	Internal   error          `json:"-"` // Не сериализуем внутреннюю ошибку
	Details    map[string]any `json:"details,omitempty"`
	HTTPStatus int            `json:"-"` // Используется методом GetStatus
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// GetStatus возвращает HTTP статус код для Huma
func (e *AppError) GetStatus() int {
	if e.HTTPStatus == 0 {
		return http.StatusInternalServerError
	}
	return e.HTTPStatus
}

// Unwrap позволяет использовать errors.Is и errors.As
func (e *AppError) Unwrap() error {
	return e.Internal
}

// InitHuma связывает AppError с механизмом создания ошибок в Huma.
// Это позволяет Huma автоматически использовать ваш формат ошибок.
func InitHuma() {
	huma.NewError = func(status int, message string, errs ...error) huma.StatusError {
		code := ErrCodeInternal
		switch status {
		case http.StatusNotFound:
			code = ErrCodeNotFound
		case http.StatusBadRequest, http.StatusUnprocessableEntity:
			code = ErrCodeValidation
		case http.StatusUnauthorized:
			code = ErrCodeUnauthorized
		case http.StatusForbidden:
			code = ErrCodeForbidden
		case http.StatusConflict:
			code = ErrCodeConflict
		}

		details := make(map[string]any)
		if len(errs) > 0 {
			for i, err := range errs {
				details[fmt.Sprintf("error_%d", i)] = err.Error()
			}
		}

		return &AppError{
			Code:       code,
			Message:    message,
			HTTPStatus: status,
			Details:    details,
		}
	}
}

// Конструкторы для частых типов ошибок

// NotFound создает ошибку "не найдено"
func NotFound(message string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
	}
}

// AlreadyExists создает ошибку "уже существует"
func AlreadyExists(message string) *AppError {
	return &AppError{
		Code:       ErrCodeAlreadyExists,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

// Validation создает ошибку валидации с деталями
func Validation(message string, details map[string]any) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		Details:    details,
		HTTPStatus: http.StatusBadRequest,
	}
}

// Internal создает внутреннюю ошибку сервера
func Internal(err error, message string) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		Internal:   err,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// BadRequest создает ошибку неверного запроса
func BadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

// Unauthorized создает ошибку неавторизованного доступа
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

// Forbidden создает ошибку запрещенного доступа
func Forbidden(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

// Conflict создает ошибку конфликта
func Conflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

// Wrap оборачивает любую ошибку в AppError
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Internal(err, message)
}

// WithDetails добавляет детали к ошибке
func (e *AppError) WithDetails(details map[string]any) *AppError {
	e.Details = details
	return e
}

// WithInternal добавляет внутреннюю ошибку
func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	return e
}
