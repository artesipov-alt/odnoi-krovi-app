package middleware

//+++====================================================+++
// 					DEPRECATED Middleware
// +++====================================================+++

import (
	"context"
	"net/http"

	"github.com/jaevor/go-nanoid"
)

// contextKey — приватный тип для ключей контекста, чтобы избежать коллизий с другими пакетами.
type contextKey string

// RequestIDKey — ключ для хранения ID запроса в контексте.
const RequestIDKey contextKey = "request_id"

// RequestIDMiddleware генерирует уникальный идентификатор для каждого входящего запроса.
// Этот ID используется для трассировки логов и возвращается клиенту в заголовке X-Request-ID.
func RequestIDMiddleware(next http.Handler) http.Handler {
	// Инициализируем генератор один раз для middleware
	idGenerator, _ := nanoid.Standard(12)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Генерируем новый ID
		id := idGenerator()

		// Кладем ID в контекст запроса, используя типизированный ключ
		ctx := context.WithValue(r.Context(), RequestIDKey, id)

		// Прокидываем в заголовок ответа, чтобы клиент мог сообщить этот ID при ошибке
		w.Header().Set("X-Request-ID", id)

		// Передаем управление дальше с обновленным контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
