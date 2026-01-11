package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// statusWriter — обертка над http.ResponseWriter для перехвата HTTP-статуса
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// LoggingMiddleware фиксирует каждый входящий запрос, его длительность и итоговый статус.
// Он не логирует детали ошибок, так как это задача ErrorHandler.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Обертка для получения статуса ответа
		ww := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		// Передача управления следующему обработчику
		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		// Логируем факт завершения запроса.
		// Используем Info уровень для всех запросов, так как детали ошибок
		// будут залогированы отдельно в ErrorHandler с уровнем Error/Warn.
		slog.InfoContext(r.Context(), "http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.status,
			"duration", duration,
			"ip", r.RemoteAddr,
		)
	})
}
