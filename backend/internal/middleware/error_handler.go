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

// ErrorHandler функция для централизованной обработки ошибок
func ErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	// Пытаемся привести к AppError
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		// Логируем в зависимости от типа ошибки
		if appErr.Internal != nil {
			// Внутренняя ошибка - логируем как error с полным контекстом
			slog.ErrorContext(r.Context(), "internal error occurred",
				"code", string(appErr.Code),
				"message", appErr.Message,
				"error", appErr.Internal,
				"path", r.URL.Path,
				"method", r.Method,
				"ip", r.RemoteAddr,
				"details", appErr.Details,
			)
		} else if appErr.HTTPStatus >= 500 {
			// Серверные ошибки без внутренней причины
			slog.ErrorContext(r.Context(), "server error",
				"code", string(appErr.Code),
				"message", appErr.Message,
				"path", r.URL.Path,
				"method", r.Method,
				"details", appErr.Details,
			)
		} else if appErr.HTTPStatus >= 400 {
			// Клиентские ошибки - логируем как warning
			slog.WarnContext(r.Context(), "client error",
				"code", string(appErr.Code),
				"message", appErr.Message,
				"path", r.URL.Path,
				"method", r.Method,
				"details", appErr.Details,
			)
		}

		// Отправляем JSON ответ и выходим
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.HTTPStatus)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		})
		return
	}

	// Непредвиденная ошибка - логируем с максимумом информации
	slog.ErrorContext(r.Context(), "unexpected error",
		"error", err,
		"path", r.URL.Path,
		"method", r.Method,
		"ip", r.RemoteAddr,
		"user_agent", r.Header.Get("User-Agent"),
	)

	// Не показываем детали непредвиденных ошибок клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    apperrors.ErrCodeInternal,
		Message: "Внутренняя ошибка сервера",
	})
}

// RecoveryMiddleware ловит панику и конвертирует в ошибку.
// Это стандартный middleware-обертка.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "panic recovered",
					"panic", rec,
					"path", r.URL.Path,
					"method", r.Method,
					"stack", string(debug.Stack()),
				)

				// Конвертируем панику в AppError
				err := apperrors.Internal(
					errors.New("panic recovered"),
					"Произошла критическая ошибка",
				)

				// Обрабатываем через ErrorHandler
				ErrorHandler(err, w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
