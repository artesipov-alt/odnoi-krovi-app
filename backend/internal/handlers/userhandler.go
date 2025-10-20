package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var users = make(map[int]*models.User)

// GetUserHandler godoc
// @Summary Получение пользователя по ID
// @Description Возвращает информацию о пользователе по его идентификатору
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User "Данные пользователя"
// @Router /user/{id} [get]
func GetUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	logger.Log.Info("get user accessed", zap.String("userId", userID))

	// Добавляем пользователя в глобальную переменную

	// Возвращаем mock-данные пользователя
	user := users[id]

	// Устанавливаем правильный Content-Type с кодировкой UTF-8
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(user)

}

// AddUserHandler godoc
// @Summary Добавление пользователя
// @Description Возвращает статус результата добавления пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "Пользователь"
// @Success 200 {object} map[string]string "Статус операции"
// @Router /user [post]
func AddUserHandler(c *fiber.Ctx) error {
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if newUser.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required and must be non-zero",
		})
	}

	// Добавляем пользователя в глобальную переменную
	users[newUser.ID] = &newUser

	logger.Log.Info("user added", zap.Int("userId", newUser.ID))

	// Устанавливаем правильный Content-Type с кодировкой UTF-8
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
