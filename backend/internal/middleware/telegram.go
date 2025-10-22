package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// TelegramAuthConfig содержит конфигурацию для Telegram аутентификации
type TelegramAuthConfig struct {
	// BotToken токен бота для проверки подписи
	BotToken string
	// RequireAuth требует ли обязательная аутентификация
	RequireAuth bool
}

// TelegramUserData представляет данные пользователя из Telegram
type TelegramUserData struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}

// DefaultTelegramAuthConfig возвращает конфигурацию по умолчанию
func DefaultTelegramAuthConfig() TelegramAuthConfig {
	return TelegramAuthConfig{
		BotToken:    "",
		RequireAuth: true,
	}
}

// TelegramAuthMiddleware создает middleware для аутентификации Telegram пользователей
func TelegramAuthMiddleware(config TelegramAuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем данные авторизации из заголовка
		initData := c.Get("Authorization")
		if initData == "" {
			if config.RequireAuth {
				logger.Log.Warn("Telegram auth data missing")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Telegram authorization data required",
				})
			}
			// Пропускаем если аутентификация не обязательна
			return c.Next()
		}

		// Убираем префикс "tma " если есть
		initData = strings.TrimPrefix(initData, "tma ")

		// Парсим данные авторизации
		userData, err := parseTelegramInitData(initData)
		if err != nil {
			logger.Log.Warn("Failed to parse telegram init data", zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid telegram authorization data",
			})
		}

		// Проверяем подпись если указан токен бота
		if config.BotToken != "" {
			if !verifyTelegramSignature(config.BotToken, initData, userData.Hash) {
				logger.Log.Warn("Invalid telegram signature")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid telegram signature",
				})
			}
		}

		// Сохраняем данные пользователя в контекст
		c.Locals("telegram_user", userData)
		c.Locals("telegram_id", userData.ID)

		logger.Log.Info("Telegram user authenticated",
			zap.Int64("telegram_id", userData.ID),
			zap.String("username", userData.Username),
		)

		return c.Next()
	}
}

// parseTelegramInitData парсит данные инициализации из Telegram
func parseTelegramInitData(initData string) (*TelegramUserData, error) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse init data: %w", err)
	}

	userData := &TelegramUserData{}

	// Парсим основные поля
	if userStr := values.Get("user"); userStr != "" {
		var user map[string]interface{}
		if err := json.Unmarshal([]byte(userStr), &user); err != nil {
			return nil, fmt.Errorf("failed to parse user data: %w", err)
		}

		if id, ok := user["id"].(float64); ok {
			userData.ID = int64(id)
		}
		if firstName, ok := user["first_name"].(string); ok {
			userData.FirstName = firstName
		}
		if lastName, ok := user["last_name"].(string); ok {
			userData.LastName = lastName
		}
		if username, ok := user["username"].(string); ok {
			userData.Username = username
		}
		if photoURL, ok := user["photo_url"].(string); ok {
			userData.PhotoURL = photoURL
		}
	}

	// Парсим остальные поля
	if authDate := values.Get("auth_date"); authDate != "" {
		if date, err := parseAuthDate(authDate); err == nil {
			userData.AuthDate = date
		}
	}

	userData.Hash = values.Get("hash")

	if userData.ID == 0 {
		return nil, fmt.Errorf("invalid telegram user data: missing user ID")
	}

	return userData, nil
}

// verifyTelegramSignature проверяет подпись Telegram
func verifyTelegramSignature(botToken, initData, hash string) bool {
	// Разбираем данные
	values, err := url.ParseQuery(initData)
	if err != nil {
		return false
	}

	// Убираем поле hash
	values.Del("hash")

	// Сортируем ключи
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Собираем строку для проверки
	var dataCheckStrings []string
	for _, key := range keys {
		value := values.Get(key)
		dataCheckStrings = append(dataCheckStrings, fmt.Sprintf("%s=%s", key, value))
	}
	dataCheckString := strings.Join(dataCheckStrings, "\n")

	// Создаем HMAC-SHA256 подпись
	secretKey := hmac.New(sha256.New, []byte("WebAppData"))
	secretKey.Write([]byte(botToken))
	secret := secretKey.Sum(nil)

	hmacHash := hmac.New(sha256.New, secret)
	hmacHash.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(hmacHash.Sum(nil))

	return hash == expectedHash
}

// parseAuthDate парсит дату авторизации
func parseAuthDate(authDate string) (int64, error) {
	var date int64
	_, err := fmt.Sscanf(authDate, "%d", &date)
	if err != nil {
		return 0, fmt.Errorf("failed to parse auth date: %w", err)
	}
	return date, nil
}

// GetTelegramUserFromContext извлекает данные Telegram пользователя из контекста
func GetTelegramUserFromContext(c *fiber.Ctx) *TelegramUserData {
	if user, ok := c.Locals("telegram_user").(*TelegramUserData); ok {
		return user
	}
	return nil
}

// GetTelegramIDFromContext извлекает Telegram ID из контекста
func GetTelegramIDFromContext(c *fiber.Ctx) int64 {
	if id, ok := c.Locals("telegram_id").(int64); ok {
		return id
	}
	return 0
}
