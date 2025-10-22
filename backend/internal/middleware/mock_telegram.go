package middleware

import (
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// MockTelegramConfig содержит конфигурацию для моковой Telegram аутентификации
type MockTelegramConfig struct {
	// MockTelegramID фиксированный Telegram ID для разработки
	MockTelegramID int64
	// MockUsername имя пользователя для разработки
	MockUsername string
	// MockFirstName имя для разработки
	MockFirstName string
	// MockLastName фамилия для разработки
	MockLastName string
	// EnableQueryParam включает получение Telegram ID из query параметра
	EnableQueryParam bool
	// QueryParamName имя query параметра для Telegram ID
	QueryParamName string
}

// DefaultMockTelegramConfig возвращает конфигурацию по умолчанию
func DefaultMockTelegramConfig() MockTelegramConfig {
	return MockTelegramConfig{
		MockTelegramID:   123456789,
		MockUsername:     "test_user",
		MockFirstName:    "Тестовый",
		MockLastName:     "Пользователь",
		EnableQueryParam: true,
		QueryParamName:   "mock_telegram_id",
	}
}

// MockTelegramAuthMiddleware создает middleware для моковой аутентификации Telegram пользователей
func MockTelegramAuthMiddleware(config MockTelegramConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var telegramID int64 = config.MockTelegramID

		// Если включено получение из query параметра, пробуем получить оттуда
		if config.EnableQueryParam {
			if idStr := c.Query(config.QueryParamName); idStr != "" {
				if id, err := parseTelegramID(idStr); err == nil {
					telegramID = id
				}
			}
		}

		// Создаем моковые данные пользователя
		userData := &TelegramUserData{
			ID:        telegramID,
			FirstName: config.MockFirstName,
			LastName:  config.MockLastName,
			Username:  config.MockUsername,
			AuthDate:  1700000000, // Фиксированная дата для разработки
		}

		// Сохраняем данные пользователя в контекст
		c.Locals("telegram_user", userData)
		c.Locals("telegram_id", userData.ID)

		logger.Log.Info("Mock Telegram user authenticated",
			zap.Int64("telegram_id", userData.ID),
			zap.String("username", userData.Username),
			zap.Bool("mock", true),
		)

		return c.Next()
	}
}

// parseTelegramID парсит Telegram ID из строки
func parseTelegramID(idStr string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("failed to parse telegram ID: %w", err)
	}
	return id, nil
}
