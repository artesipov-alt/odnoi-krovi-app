package middleware

import (
	"net/http"

	sloghttp "github.com/samber/slog-http"
)

func TraceIDResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// trace ID уже есть в контексте если sloghttp идёт раньше
		if traceID := sloghttp.GetRequestIDFromContext(r.Context()); traceID != "" {
			w.Header().Set("X-Trace-Id", traceID)
		}
		next.ServeHTTP(w, r)
	})
}
