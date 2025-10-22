package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService services.UserService
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SimpleRegistrationRequest represents a request for simple user registration
type SimpleRegistrationRequest struct {
	TelegramID int64  `json:"telegram_id" validate:"required,min=1" example:"123456789"`
	FullName   string `json:"full_name,omitempty" validate:"omitempty,min=1,max=255" example:"Иван Иванов"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserHandler godoc
// @Summary Получение пользователя по ID
// @Description Возвращает информацию о пользователе по его идентификатору
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} models.User "Данные пользователя"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/{id} [get]
func (h *UserHandler) GetUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "User ID is required",
		})
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid user ID format",
		})
	}

	logger.Log.Info("get user accessed", zap.String("userId", userID))

	user, err := h.userService.GetUserProfile(c.Context(), id)
	if err != nil {
		logger.Log.Error("failed to get user", zap.Error(err), zap.Int("userId", id))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get user",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error: "User not found",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(user.User)
}

// RegisterUserSimpleHandler godoc
// @Summary Простая регистрация пользователя
// @Description Создает пользователя с Telegram ID и именем (для команды Start)
// @Tags users
// @Accept json
// @Produce json
// @Param request body SimpleRegistrationRequest true "Данные для простой регистрации"
// @Success 201 {object} models.User "Зарегистрированный пользователь"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 409 {object} ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/register/simple [post]
func (h *UserHandler) RegisterUserSimpleHandler(c *fiber.Ctx) error {
	var request SimpleRegistrationRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if request.TelegramID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Valid Telegram ID is required",
		})
	}

	// Use provided full name or default to "Telegram User"
	fullName := request.FullName
	if fullName == "" {
		fullName = "Telegram User"
	}

	user, err := h.userService.RegisterUserSimple(c.Context(), request.TelegramID, fullName)
	if err != nil {
		logger.Log.Error("failed to register user", zap.Error(err), zap.Int64("telegramId", request.TelegramID))

		if err.Error() == "user with this telegram ID already exists" {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "User with this Telegram ID already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to register user",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Status(fiber.StatusCreated).JSON(user)
}

// RegisterUserHandler godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param request body services.UserRegistration true "Данные для регистрации пользователя"
// @Success 201 {object} models.User "Зарегистрированный пользователь"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 409 {object} ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/register [post]
func (h *UserHandler) RegisterUserHandler(c *fiber.Ctx) error {
	var registrationData services.UserRegistration
	if err := c.BodyParser(&registrationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Get Telegram ID from context (should be set by middleware)
	telegramID, ok := c.Locals("telegram_id").(int64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Telegram ID is required",
		})
	}

	user, err := h.userService.RegisterUser(c.Context(), telegramID, registrationData)
	if err != nil {
		logger.Log.Error("failed to register user", zap.Error(err), zap.Int64("telegramId", telegramID))

		if err.Error() == "user with this telegram ID already exists" {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "User with this Telegram ID already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to register user",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUserHandler godoc
// @Summary Обновление данных пользователя
// @Description Обновляет информацию о пользователе
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param request body services.UserUpdate true "Данные для обновления"
// @Success 200 {object} SuccessResponse "Данные успешно обновлены"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/{id} [put]
func (h *UserHandler) UpdateUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "User ID is required",
		})
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid user ID format",
		})
	}

	var updateData services.UserUpdate
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	if err := h.userService.UpdateUserProfile(c.Context(), id, updateData); err != nil {
		logger.Log.Error("failed to update user", zap.Error(err), zap.Int("userId", id))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to update user",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(SuccessResponse{
		Message: "User updated successfully",
	})
}

// GetUserByTelegramHandler godoc
// @Summary Получение пользователя по Telegram ID
// @Description Возвращает информацию о пользователе по его Telegram ID
// @Tags users
// @Produce json
// @Param telegram_id query int64 true "Telegram ID пользователя"
// @Success 200 {object} models.User "Данные пользователя"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/telegram [get]
func (h *UserHandler) GetUserByTelegramHandler(c *fiber.Ctx) error {
	telegramIDStr := c.Query("telegram_id")
	if telegramIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Telegram ID is required",
		})
	}

	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid Telegram ID format",
		})
	}

	user, err := h.userService.GetUserByTelegramID(c.Context(), telegramID)
	if err != nil {
		logger.Log.Error("failed to get user by telegram ID", zap.Error(err), zap.Int64("telegramId", telegramID))

		if err.Error() == "user with telegram id not found" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get user",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(user)
}
