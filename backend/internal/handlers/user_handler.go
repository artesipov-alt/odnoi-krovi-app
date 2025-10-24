package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler обрабатывает HTTP запросы для операций с пользователями
type UserHandler struct {
	userService services.UserService
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// SimpleRegistrationRequest представляет запрос на простую регистрацию пользователя
type SimpleRegistrationRequest struct {
	TelegramID int64  `json:"telegram_id" validate:"required,min=1" example:"123456789"`
	FullName   string `json:"full_name,omitempty" validate:"omitempty,min=1,max=255" example:"Иван Иванов"`
}

// SuccessResponse представляет успешный ответ
type SuccessResponse struct {
	Message string `json:"message"`
}

// NewUserHandler создает новый обработчик пользователей
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
			Error: "ID пользователя обязателен",
		})
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID пользователя",
		})
	}

	logger.Log.Info("доступ к получению пользователя", zap.String("userId", userID))

	user, err := h.userService.GetUserProfile(c.Context(), id)
	if err != nil {
		logger.Log.Error("не удалось получить пользователя", zap.Error(err), zap.Int("userId", id))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error: "Пользователь не найден",
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
			Error: "Неверное тело запроса",
		})
	}

	// Validate request
	if err := validation.ValidateStruct(request); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	if request.TelegramID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Корректный Telegram ID обязателен",
		})
	}

	// Использовать предоставленное полное имя или установить значение по умолчанию "Пользователь Telegram"
	fullName := request.FullName
	if fullName == "" {
		fullName = "Пользователь Telegram"
	}

	user, err := h.userService.RegisterUserSimple(c.Context(), request.TelegramID, fullName)
	if err != nil {
		logger.Log.Error("не удалось зарегистрировать пользователя", zap.Error(err), zap.Int64("telegramId", request.TelegramID))

		if err.Error() == "пользователь с этим telegram ID уже существует" {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
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
// @Deprecated
// @Router /user/register [post]
func (h *UserHandler) RegisterUserHandler(c *fiber.Ctx) error {
	var registrationData services.UserRegistration
	if err := c.BodyParser(&registrationData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверное тело запроса",
		})
	}

	// Validate registration data
	if err := validation.ValidateStruct(registrationData); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	// Получить Telegram ID из контекста (должен быть установлен промежуточным ПО)
	telegramID, ok := c.Locals("telegram_id").(int64)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Telegram ID обязателен",
		})
	}

	user, err := h.userService.RegisterUser(c.Context(), telegramID, registrationData)
	if err != nil {
		logger.Log.Error("не удалось зарегистрировать пользователя", zap.Error(err), zap.Int64("telegramId", telegramID))

		if err.Error() == "пользователь с этим telegram ID уже существует" {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
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
			Error: "ID пользователя обязателен",
		})
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID пользователя",
		})
	}

	var updateData services.UserUpdate
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверное тело запроса",
		})
	}

	// Validate update data
	if err := validation.ValidateStruct(updateData); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	if err := h.userService.UpdateUserProfile(c.Context(), id, updateData); err != nil {
		logger.Log.Error("не удалось обновить пользователя", zap.Error(err), zap.Int("userId", id))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(SuccessResponse{
		Message: "Пользователь успешно обновлен",
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
			Error: "Telegram ID обязателен",
		})
	}

	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат Telegram ID",
		})
	}

	user, err := h.userService.GetUserByTelegramID(c.Context(), telegramID)
	if err != nil {
		logger.Log.Error("не удалось получить пользователя по Telegram ID", zap.Error(err), zap.Int64("telegramId", telegramID))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(user)
}

// DeleteUserHandler godoc
// @Summary Удаление пользователя по ID
// @Description Удаляет пользователя из системы (soft delete)
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} SuccessResponse "Пользователь успешно удален"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID пользователя обязателен",
		})
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID пользователя",
		})
	}

	logger.Log.Info("доступ к удалению пользователя", zap.Int("userId", id))

	if err := h.userService.DeleteUser(c.Context(), id); err != nil {
		logger.Log.Error("не удалось удалить пользователя", zap.Error(err), zap.Int("userId", id))

		if err.Error() == "пользователь не найден" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(SuccessResponse{
		Message: "Пользователь успешно удален",
	})
}
