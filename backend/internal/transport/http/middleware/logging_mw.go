package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// statusWriter — обертка над http.ResponseWriter для перехвата HTTP-статуса ответа.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// LoggingMiddleware фиксирует каждый входящий запрос, его длительность и итоговый статус.
// Он безопасно извлекает RequestID из контекста и логирует информацию о запросе.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Обертка для получения статуса ответа (по умолчанию 200 OK)
		ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		// Передача управления следующему обработчику
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Безопасно получаем request_id из контекста.
		// Используем типизированный ключ RequestIDKey из req_id_mw.go.
		// Если ID не найден или тип не совпадает, используем "unknown".
		id, ok := r.Context().Value(RequestIDKey).(string)
		if !ok {
			id = "unknown"
		}

		// Логируем запрос с уровнем в зависимости от статуса
		var message string
		var level slog.Level
		if ww.status >= 500 {
			level = slog.LevelError
			message = "Internal server error"
		} else if ww.status >= 400 {
			level = slog.LevelWarn
			message = "Client error"
		} else {
			level = slog.LevelInfo
			message = "http request"
		}

		slog.Log(r.Context(), level, message,
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.status,
			"duration", duration,
			"ip", r.RemoteAddr,
			"request_id", id,
		)
	})
}
