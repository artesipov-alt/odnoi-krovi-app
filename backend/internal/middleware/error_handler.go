package middleware

import (
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorResponse представляет стандартный ответ с ошибкой
type ErrorResponse struct {
	Code    apperrors.ErrorCode    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorHandler middleware для централизованной обработки ошибок
// Использует этот обработчик в fiber.Config при создании приложения
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Пытаемся привести к AppError
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			// Логируем в зависимости от типа ошибки
			if appErr.Internal != nil {
				// Внутренняя ошибка - логируем как error с полным контекстом
				logger.Log.Error("internal error occurred",
					zap.String("code", string(appErr.Code)),
					zap.String("message", appErr.Message),
					zap.Error(appErr.Internal),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.String("ip", c.IP()),
					zap.Any("details", appErr.Details),
				)
			} else if appErr.HTTPStatus >= 500 {
				// Серверные ошибки без внутренней причины
				logger.Log.Error("server error",
					zap.String("code", string(appErr.Code)),
					zap.String("message", appErr.Message),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.Any("details", appErr.Details),
				)
			} else if appErr.HTTPStatus >= 400 {
				// Клиентские ошибки - логируем как warning
				logger.Log.Warn("client error",
					zap.String("code", string(appErr.Code)),
					zap.String("message", appErr.Message),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.Any("details", appErr.Details),
				)
			}

			// Отправляем JSON ответ
			c.Set("Content-Type", "application/json; charset=utf-8")
			return c.Status(appErr.HTTPStatus).JSON(ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			})
		}

		// Если это ошибка Fiber
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			// Для 404 используем INFO (без stack trace), для остальных - WARN
			if fiberErr.Code == 404 {
				logger.Log.Info("not found",
					zap.Int("status", fiberErr.Code),
					zap.String("message", fiberErr.Message),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)
			} else {
				logger.Log.Warn("fiber error",
					zap.Int("status", fiberErr.Code),
					zap.String("message", fiberErr.Message),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
				)
			}

			c.Set("Content-Type", "application/json; charset=utf-8")
			return c.Status(fiberErr.Code).JSON(ErrorResponse{
				Code:    apperrors.ErrCodeBadRequest,
				Message: fiberErr.Message,
			})
		}

		// Непредвиденная ошибка - логируем с максимумом информации
		logger.Log.Error("unexpected error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		// Не показываем детали непредвиденных ошибок клиенту
		c.Set("Content-Type", "application/json; charset=utf-8")
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Code:    apperrors.ErrCodeInternal,
			Message: "Внутренняя ошибка сервера",
		})
	}
}

// RecoveryMiddleware ловит панику и конвертирует в ошибку
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				logger.Log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("path", c.Path()),
					zap.String("method", c.Method()),
					zap.Stack("stack"),
				)

				// Конвертируем панику в AppError
				err := apperrors.Internal(
					errors.New("panic recovered"),
					"Произошла критическая ошибка",
				)

				// Обрабатываем через ErrorHandler
				_ = ErrorHandler()(c, err)
			}
		}()

		return c.Next()
	}
}
