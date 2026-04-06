package middleware

//+++====================================================+++
// 					DEPRECATED Middleware
// +++====================================================+++

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/danielgtaylor/huma/v2"
)

// ErrorResponse представляет стандартный ответ с ошибкой, совместимый с форматом Huma.
// Это гарантирует, что ошибки из middleware и ошибки из Huma-обработчиков выглядят одинаково.
type ErrorResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

// ErrorHandler — централизованная функция для отправки ошибок в формате JSON.
// Она умеет распознавать ошибки Huma (StatusError) и правильно их форматировать.
func ErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	var status int
	var body ErrorResponse

	// Пытаемся привести ошибку к интерфейсу Huma StatusError
	var humaErr huma.StatusError
	if errors.As(err, &humaErr) {
		status = humaErr.GetStatus()
		body = ErrorResponse{
			Status: status,
			Title:  http.StatusText(status),
			Detail: humaErr.Error(),
		}
	} else {
		// Для всех остальных (непредвиденных) ошибок возвращаем 500
		status = http.StatusInternalServerError
		body = ErrorResponse{
			Status: status,
			Title:  "Internal Server Error",
			Detail: "Внутренняя ошибка сервера",
		}
	}

	// Логируем ошибку в зависимости от её критичности
	if status >= 500 {
		slog.ErrorContext(r.Context(), "internal server error",
			"error", err.Error(),
			"path", r.URL.Path,
			"stack", string(debug.Stack()),
		)
	} else {
		slog.WarnContext(r.Context(), "application warning",
			"status", status,
			"error", err.Error(),
			"path", r.URL.Path,
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// RecoveryMiddleware перехватывает паники и возвращает клиенту 500 ошибку вместо падения процесса.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "panic recovered",
					"panic", rec,
					"stack", string(debug.Stack()),
				)

				// Создаем ошибку через Huma, чтобы ErrorHandler её правильно обработал
				err := huma.Error500InternalServerError("Критическая ошибка сервера")
				ErrorHandler(err, w, r)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
