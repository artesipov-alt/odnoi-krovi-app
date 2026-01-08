package middleware

import (
	"errors"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ErrorResponse представляет стандартный ответ с ошибкой
type ErrorResponse struct {
	Code    apperrors.ErrorCode `json:"code"`
	Message string              `json:"message"`
	Details map[string]any      `json:"details,omitempty"`
}

// ErrorHandler middleware для централизованной обработки ошибок
// Использует этот обработчик в Echo при создании приложения
func ErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		// Если ответ уже отправлен, ничего не делаем
		if c.Response().Committed {
			return
		}

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
					zap.String("method", c.Request().Method),
					zap.String("ip", c.RealIP()),
					zap.Any("details", appErr.Details),
				)
			} else if appErr.HTTPStatus >= 500 {
				// Серверные ошибки без внутренней причины
				logger.Log.Error("server error",
					zap.String("code", string(appErr.Code)),
					zap.String("message", appErr.Message),
					zap.String("path", c.Path()),
					zap.String("method", c.Request().Method),
					zap.Any("details", appErr.Details),
				)
			} else if appErr.HTTPStatus >= 400 {
				// Клиентские ошибки - логируем как warning
				logger.Log.Warn("client error",
					zap.String("code", string(appErr.Code)),
					zap.String("message", appErr.Message),
					zap.String("path", c.Path()),
					zap.String("method", c.Request().Method),
					zap.Any("details", appErr.Details),
				)
			}

			// Отправляем JSON ответ и выходим
			c.JSON(appErr.HTTPStatus, ErrorResponse{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			})
			return
		}

		// Если это ошибка Echo
		var echoErr *echo.HTTPError
		if errors.As(err, &echoErr) {
			// Получаем сообщение с type assertion
			message, ok := echoErr.Message.(string)
			if !ok {
				message = "Unknown error"
			}

			// Для 404 используем INFO (без stack trace), для остальных - WARN
			if echoErr.Code == 404 {
				logger.Log.Info("not found",
					zap.Int("status", echoErr.Code),
					zap.String("message", message),
					zap.String("path", c.Path()),
					zap.String("method", c.Request().Method),
				)
			} else {
				logger.Log.Warn("echo error",
					zap.Int("status", echoErr.Code),
					zap.String("message", message),
					zap.String("path", c.Path()),
					zap.String("method", c.Request().Method),
				)
			}

			c.JSON(echoErr.Code, ErrorResponse{
				Code:    apperrors.ErrCodeBadRequest,
				Message: message,
			})
			return
		}

		// Непредвиденная ошибка - логируем с максимумом информации
		logger.Log.Error("unexpected error",
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Request().Method),
			zap.String("ip", c.RealIP()),
			zap.String("user_agent", c.Request().Header.Get("User-Agent")),
		)

		// Не показываем детали непредвиденных ошибок клиенту
		c.JSON(500, ErrorResponse{
			Code:    apperrors.ErrCodeInternal,
			Message: "Внутренняя ошибка сервера",
		})
	}
}

// RecoveryMiddleware ловит панику и конвертирует в ошибку
func RecoveryMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					logger.Log.Error("panic recovered",
						zap.Any("panic", r),
						zap.String("path", c.Path()),
						zap.String("method", c.Request().Method),
						zap.Stack("stack"),
					)

					// Конвертируем панику в AppError
					err := apperrors.Internal(
						errors.New("panic recovered"),
						"Произошла критическая ошибка",
					)

					// Обрабатываем через ErrorHandler
					c.Error(err)
				}
			}()

			return next(c)
		}
	}
}
