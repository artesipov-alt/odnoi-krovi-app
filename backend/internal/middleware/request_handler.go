package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/jaevor/go-nanoid"
)

func RequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gen, _ := nanoid.Standard(12)
		id := gen()
		ctx := context.WithValue(r.Context(), "request_id", id)
		slog.InfoContext(r.Context(), "req", r.Method, r.URL.Path, "id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
