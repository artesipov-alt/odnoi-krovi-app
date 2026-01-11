package handlers

// import (
// 	"strconv"

// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
// 	"github.com/labstack/echo/v4"
// 	"go.uber.org/zap"
// )

// // UserHandler обрабатывает HTTP запросы для операций с пользователями
// type UserHandler struct {
// 	userService services.UserService
// }

// // NewUserHandler создает новый обработчик пользователей
// func NewUserHandler(userService services.UserService) *UserHandler {
// 	return &UserHandler{
// 		userService: userService,
// 	}
// }

// // GetUserHandler godoc
// // @Summary Получение пользователя по ID
// // @Description Возвращает информацию о пользователе по его идентификатору
// // @Tags users-v1
// // @Produce json
// // @Param id path string true "ID пользователя"
// // @Success 200 {object} ent.User "Данные пользователя"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /v1/user/{id} [get]
// func (h *UserHandler) GetUserHandler(c echo.Context) error {
// 	id := c.Param("id")
// 	if id == "" {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid ID"})
// 	}

// 	logger.Log.Info("получение пользователя", zap.String("userId", id))

// 	user, err := h.userService.GetUserByID(c.Request().Context(), id)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(200, user)
// }

// // RegisterUserSimpleHandler godoc
// // @Summary Простая регистрация пользователя
// // @Description Создает пользователя с Telegram ID и именем (для команды Start)
// // @Tags users-v1
// // @Accept json
// // @Produce json
// // @Param request body dto.SimpleRegistrationRequest true "Данные для простой регистрации"
// // @Success 201 {object} ent.User "Зарегистрированный пользователь"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 409 {object} utils.ErrorResponse "Пользователь уже существует"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /v1/user/register/simple [post]
// func (h *UserHandler) RegisterUserSimpleHandler(c echo.Context) error {
// 	var request dto.SimpleRegistrationRequest
// 	if err := c.Bind(&request); err != nil {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid request"})
// 	}

// 	// Использовать предоставленное полное имя или установить значение по умолчанию "Пользователь Telegram"
// 	fullName := request.FullName
// 	if fullName == "" {
// 		fullName = "Пользователь Telegram"
// 	}

// 	logger.Log.Info("регистрация пользователя", zap.Int64("telegramId", request.TelegramID))

// 	user, err := h.userService.RegisterUserSimple(c.Request().Context(), request.TelegramID, fullName)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(201, user)
// }

// // RegisterUserHandler godoc
// // @Summary Регистрация нового пользователя
// // @Description Регистрирует нового пользователя в системе
// // @Tags users-v1
// // @Accept json
// // @Produce json
// // @Param request body dto.UserRegistration true "Данные для регистрации пользователя"
// // @Success 201 {object} ent.User "Зарегистрированный пользователь"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 409 {object} utils.ErrorResponse "Пользователь уже существует"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Deprecated
// // @Router /v1/user/register [post]
// func (h *UserHandler) RegisterUserHandler(c echo.Context) error {
// 	var registrationData dto.UserRegistration
// 	if err := c.Bind(&registrationData); err != nil {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid request"})
// 	}

// 	// Получить Telegram ID из контекста (должен быть установлен промежуточным ПО)
// 	telegramID, ok := c.Get("telegram_id").(int64)
// 	if !ok {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Telegram ID обязателен"})
// 	}

// 	logger.Log.Info("регистрация пользователя", zap.Int64("telegramId", telegramID))

// 	user, err := h.userService.RegisterUser(c.Request().Context(), telegramID, registrationData)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(201, user)
// }

// // UpdateUserHandler godoc
// // @Summary Обновление данных пользователя
// // @Description Обновляет информацию о пользователе
// // @Tags users-v1
// // @Accept json
// // @Produce json
// // @Param id path string true "ID пользователя"
// // @Param request body ent.User true "Данные для обновления"
// // @Success 200 {object} utils.SuccessResponse "Данные успешно обновлены"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /v1/user/{id} [put]
// func (h *UserHandler) UpdateUserHandler(c echo.Context) error {
// 	id := c.Param("id")
// 	if id == "" {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid ID"})
// 	}

// 	var updateData dto.UserUpdate
// 	if err := c.Bind(&updateData); err != nil {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid request"})
// 	}

// 	logger.Log.Info("обновление пользователя", zap.String("userId", id))

// 	if err := h.userService.UpdateUserProfile(c.Request().Context(), id, updateData); err != nil {
// 		return err
// 	}

// 	return c.JSON(200, map[string]string{"message": "Пользователь успешно обновлен"})
// }

// // GetUserByTelegramHandler godoc
// // @Summary Получение пользователя по Telegram ID
// // @Description Возвращает информацию о пользователе по его Telegram ID
// // @Tags users-v1
// // @Produce json
// // @Param telegram_id query int64 true "Telegram ID пользователя"
// // @Success 200 {object} ent.User "Данные пользователя"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /v1/user/telegram [get]
// func (h *UserHandler) GetUserByTelegramHandler(c echo.Context) error {
// 	telegramIDStr := c.QueryParam("telegram_id")
// 	if telegramIDStr == "" {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid telegram_id"})
// 	}
// 	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
// 	if err != nil {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid telegram_id"})
// 	}

// 	logger.Log.Info("получение пользователя по Telegram ID", zap.Int64("telegramId", telegramID))

// 	user, err := h.userService.GetUserByTelegramID(c.Request().Context(), telegramID)
// 	if err != nil {
// 		return err
// 	}

// 	return c.JSON(200, user)
// }

// // DeleteUserHandler godoc
// // @Summary Удаление пользователя по ID
// // @Description Удаляет пользователя из системы (soft delete)
// // @Tags users-v1
// // @Produce json
// // @Param id path string true "ID пользователя"
// // @Success 200 {object} utils.SuccessResponse "Пользователь успешно удален"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /v1/user/{id} [delete]
// func (h *UserHandler) DeleteUserHandler(c echo.Context) error {
// 	id := c.Param("id")
// 	if id == "" {
// 		return c.JSON(400, utils.ErrorResponse{Message: "Invalid ID"})
// 	}

// 	logger.Log.Info("удаление пользователя", zap.String("userId", id))

// 	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
// 		return err
// 	}

// 	return c.JSON(200, map[string]string{"message": "Пользователь успешно удален"})
// }
