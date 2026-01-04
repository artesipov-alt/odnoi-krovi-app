package middleware

import (
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware логирует HTTP запросы с использованием кастомного логгера
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			req := c.Request()
			res := c.Response()

			logger.Log.Info("HTTP request",
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", res.Status),
				zap.Duration("latency", stop.Sub(start)),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", req.Header.Get("User-Agent")),
			)

			return err
		}
	}
}
