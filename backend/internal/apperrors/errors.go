package apperrors

import (
	"errors"
	"fmt"
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

// AppError представляет ошибку приложения с метаданными
type AppError struct {
	Code       ErrorCode      // Код ошибки для API
	Message    string         // Сообщение для пользователя
	Internal   error          // Внутренняя ошибка (для логов)
	Details    map[string]any // Дополнительные детали
	HTTPStatus int            // HTTP статус код
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// Unwrap позволяет использовать errors.Is и errors.As
func (e *AppError) Unwrap() error {
	return e.Internal
}

// Конструкторы для частых типов ошибок

// NotFound создает ошибку "не найдено"
func NotFound(message string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    message,
		HTTPStatus: 404,
	}
}

// AlreadyExists создает ошибку "уже существует"
func AlreadyExists(message string) *AppError {
	return &AppError{
		Code:       ErrCodeAlreadyExists,
		Message:    message,
		HTTPStatus: 409,
	}
}

// Validation создает ошибку валидации с деталями
func Validation(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		Details:    details,
		HTTPStatus: 400,
	}
}

// Internal создает внутреннюю ошибку сервера
func Internal(err error, message string) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		Internal:   err,
		HTTPStatus: 500,
	}
}

// BadRequest создает ошибку неверного запроса
func BadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		HTTPStatus: 400,
	}
}

// Unauthorized создает ошибку неавторизованного доступа
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: 401,
	}
}

// Forbidden создает ошибку запрещенного доступа
func Forbidden(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: 403,
	}
}

// Conflict создает ошибку конфликта
func Conflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: 409,
	}
}

// Wrap оборачивает любую ошибку в AppError
func Wrap(err error, message string) *AppError {
	if err == nil {
		return nil
	}

	// Если уже AppError - возвращаем как есть
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	// Иначе оборачиваем как Internal
	return Internal(err, message)
}

// WithDetails добавляет детали к ошибке
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithInternal добавляет внутреннюю ошибку
func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	return e
}
