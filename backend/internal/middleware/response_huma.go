package middleware

// import (
// 	"context"
// 	"log/slog"
// 	"net/http"

// 	"github.com/danielgtaylor/huma/v2"
// )

// // func ResponseHandler(next huma.Context) huma.Handler {
// // 	return func(ctx context.Context, req *http.Request, res *http.Response) error {
// // 		id := ctx.Value("request_id").(string)
// // 		slog.InfoContext(ctx, "req", req.Method, req.URL.Path, "id", id)
// // 		return next(ctx, req, res)
// // 	}
// // }
