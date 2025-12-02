package middleware

import (
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/pkg/logger"
	"go.uber.org/zap"
)

// ZapLogger - middleware для логирования HTTP запросов с использованием zap логгера
func ZapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем response writer для отслеживания статуса
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		// Обрабатываем запрос
		next.ServeHTTP(ww, r)

		// Вычисляем время выполнения
		duration := time.Since(start)

		// Получаем Request ID из контекста
		requestID := GetRequestIDFromContext(r.Context())
		if requestID == "" {
			// Если нет в контексте, пробуем из заголовка
			requestID = r.Header.Get("X-Request-ID")
		}

		// Определяем уровень логирования на основе статуса кода
		logEntry := logger.Log.With(
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("ip", r.RemoteAddr),
			zap.String("user_agent", r.UserAgent()),
			zap.String("proto", r.Proto),
			zap.String("host", r.Host),
			zap.Int("status", ww.status),
			zap.Duration("duration", duration),
			zap.Int64("bytes", ww.bytesWritten),
			zap.String("request_id", requestID),
		)

		// Логируем с разными уровнями в зависимости от статуса кода
		switch {
		case ww.status >= 500:
			logEntry.Error("HTTP request error")
		case ww.status >= 400:
			logEntry.Warn("HTTP request client error")
		case ww.status >= 300:
			logEntry.Info("HTTP request redirect")
		default:
			logEntry.Info("HTTP request success")
		}
	})
}

// responseWriter - обертка для http.ResponseWriter для отслеживания статуса и размера ответа
type responseWriter struct {
	http.ResponseWriter
	status       int
	bytesWritten int64
	wroteHeader  bool
}

// WriteHeader перехватывает статус код
func (rw *responseWriter) WriteHeader(statusCode int) {
	if !rw.wroteHeader {
		rw.status = statusCode
		rw.wroteHeader = true
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write перехватывает запись данных
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += int64(n)
	return n, err
}

// Status возвращает статус код ответа
func (rw *responseWriter) Status() int {
	return rw.status
}

// BytesWritten возвращает количество записанных байт
func (rw *responseWriter) BytesWritten() int64 {
	return rw.bytesWritten
}
