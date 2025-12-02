package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	// RequestIDKey ключ для хранения Request ID в контексте
	RequestIDKey contextKey = "request_id"
	// RequestIDHeader заголовок для передачи Request ID
	RequestIDHeader = "X-Request-ID"
)

// RequestID middleware добавляет уникальный идентификатор к каждому запросу
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем или генерируем Request ID
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Устанавливаем заголовок в запрос для логирования
		r.Header.Set(RequestIDHeader, requestID)

		// Устанавливаем заголовок в ответ
		w.Header().Set(RequestIDHeader, requestID)

		// Добавляем Request ID в контекст запроса
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Продолжаем обработку запроса
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestIDFromContext извлекает Request ID из контекста
func GetRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// generateRequestID генерирует уникальный идентификатор запроса
func generateRequestID() string {
	return uuid.New().String()
}
