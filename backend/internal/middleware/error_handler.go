package middleware

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
)

// ErrorResponse представляет стандартный ответ с ошибкой
type ErrorResponse struct {
	Code    apperrors.ErrorCode `json:"code"`
	Message string              `json:"message"`
	Details map[string]any      `json:"details,omitempty"`
}

// ErrorHandler функция для централизованной обработки ошибок.
// Она используется в тех местах, где нужно вручную отправить ошибку в формате AppError.
func ErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	var status int
	var body ErrorResponse

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		status = appErr.GetStatus()
		body = ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		}

		// Логируем в зависимости от статуса
		if status >= 500 {
			slog.ErrorContext(r.Context(), "internal server error",
				"code", appErr.Code,
				"message", appErr.Message,
				"internal_err", appErr.Internal,
				"details", appErr.Details,
				"path", r.URL.Path,
			)
		} else {
			slog.WarnContext(r.Context(), "application warning",
				"code", appErr.Code,
				"message", appErr.Message,
				"status", status,
				"path", r.URL.Path,
			)
		}
	} else {
		// Непредвиденная ошибка
		status = http.StatusInternalServerError
		body = ErrorResponse{
			Code:    apperrors.ErrCodeInternal,
			Message: "Внутренняя ошибка сервера",
		}

		slog.ErrorContext(r.Context(), "unexpected error occurred",
			"error", err.Error(),
			"path", r.URL.Path,
			"stack", string(debug.Stack()),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// RecoveryMiddleware ловит панику и предотвращает падение приложения.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "panic recovered",
					"panic", rec,
					"stack", string(debug.Stack()),
				)

				// Используем наш ErrorHandler для отправки ответа
				err := apperrors.Internal(errors.New("panic"), "Критическая ошибка сервера")
				ErrorHandler(err, w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
