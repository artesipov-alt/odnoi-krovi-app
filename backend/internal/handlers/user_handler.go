package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler обрабатывает HTTP запросы для операций с пользователями
type UserHandler struct {
	userService services.UserService
}

// SimpleRegistrationRequest представляет запрос на простую регистрацию пользователя
type SimpleRegistrationRequest struct {
	TelegramID int64  `json:"telegramId" validate:"required,min=1" example:"123456789"`
	FullName   string `json:"fullName,omitempty" validate:"omitempty,min=1,max=255" example:"Иван Иванов"`
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
	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение пользователя", zap.Int("userId", id))

	profile, err := h.userService.GetUserProfile(c.Context(), id)
	if err != nil {
		return err
	}

	return SendJSON(c, profile.User)
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
	if err := ParseBody(c, &request); err != nil {
		return err
	}

	// Использовать предоставленное полное имя или установить значение по умолчанию "Пользователь Telegram"
	fullName := request.FullName
	if fullName == "" {
		fullName = "Пользователь Telegram"
	}

	logger.Log.Info("регистрация пользователя", zap.Int64("telegramId", request.TelegramID))

	user, err := h.userService.RegisterUserSimple(c.Context(), request.TelegramID, fullName)
	if err != nil {
		return err
	}

	return SendCreated(c, user)
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
	if err := ParseBody(c, &registrationData); err != nil {
		return err
	}

	// Получить Telegram ID из контекста (должен быть установлен промежуточным ПО)
	telegramID, ok := c.Locals("telegram_id").(int64)
	if !ok {
		return apperrors.BadRequest("Telegram ID обязателен")
	}

	logger.Log.Info("регистрация пользователя", zap.Int64("telegramId", telegramID))

	user, err := h.userService.RegisterUser(c.Context(), telegramID, registrationData)
	if err != nil {
		return err
	}

	return SendCreated(c, user)
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
	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	var updateData services.UserUpdate
	if err := ParseBody(c, &updateData); err != nil {
		return err
	}

	logger.Log.Info("обновление пользователя", zap.Int("userId", id))

	if err := h.userService.UpdateUserProfile(c.Context(), id, updateData); err != nil {
		return err
	}

	return SendSuccess(c, "Пользователь успешно обновлен")
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
	telegramID, err := ParseInt64Query(c, "telegram_id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение пользователя по Telegram ID", zap.Int64("telegramId", telegramID))

	user, err := h.userService.GetUserByTelegramID(c.Context(), telegramID)
	if err != nil {
		return err
	}

	return SendJSON(c, user)
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
	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("удаление пользователя", zap.Int("userId", id))

	if err := h.userService.DeleteUser(c.Context(), id); err != nil {
		return err
	}

	return SendSuccess(c, "Пользователь успешно удален")
}
