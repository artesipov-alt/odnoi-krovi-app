package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ReferenceHandler обрабатывает HTTP запросы для справочных данных
type DevHandler struct {
	userRepo repositories.UserRepository
}

// NewReferenceHandler создает новый обработчик справочных данных
func NewDevHandler(userRepo repositories.UserRepository) *DevHandler {
	return &DevHandler{
		userRepo: userRepo,
	}
}

// ReferenceResponse представляет ответ со справочными данными
type DevResponse struct {
	Status  bool
	Message string
}

// ResetUserHandler godoc
// @Summary Сброс пользователя к начальным настройкам
// @Description Сбрасывает пользователя к заводским настройкам на этапе команды старт от бота
// @Tags dev
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} DevResponse "Успешный сброс пользователя"
// @Router /dev/reset-user/{id} [post]
func (h *DevHandler) ResetUserHandler(c *fiber.Ctx) error {
	logger.Log.Info("Сброс пользователя к заводским настройкам")

	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("Сброс пользователя", zap.Int("userId", id))

	if err := h.userRepo.ResetUser(c.Context(), id); err != nil {
		logger.Log.Error("Ошибка при сбросе пользователя", zap.Error(err), zap.Int("userId", id))
		return err
	}

	logger.Log.Info("Пользователь успешно сброшен", zap.Int("userId", id))

	return SendJSON(c, DevResponse{
		Status:  true,
		Message: "Пользователь успешно сброшен к заводским настройкам",
	})
}

// RestoreUserHandler godoc
// @Summary Восстановление удаленного пользователя
// @Description Восстанавливает мягко удаленного пользователя, устанавливая deleted_at в NULL
// @Tags dev
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} DevResponse "Успешное восстановление пользователя"
// @Router /dev/restore-user/{id} [post]
func (h *DevHandler) RestoreUserHandler(c *fiber.Ctx) error {
	logger.Log.Info("Восстановление удаленного пользователя")

	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("Восстановление пользователя", zap.Int("userId", id))

	if err := h.userRepo.RestoreUser(c.Context(), id); err != nil {
		logger.Log.Error("Ошибка при восстановлении пользователя", zap.Error(err), zap.Int("userId", id))
		return err
	}

	logger.Log.Info("Пользователь успешно восстановлен", zap.Int("userId", id))

	return SendJSON(c, DevResponse{
		Status:  true,
		Message: "Пользователь успешно восстановлен",
	})
}

// GetDeletedUsersHandler godoc
// @Summary Получение всех удаленных пользователей
// @Description Возвращает список всех мягко удаленных пользователей
// @Tags dev
// @Produce json
// @Success 200 {object} GetDeletedUsersResponse "Список удаленных пользователей"
// @Router /dev/deleted-users [get]
func (h *DevHandler) GetDeletedUsersHandler(c *fiber.Ctx) error {
	logger.Log.Info("Получение списка удаленных пользователей")

	users, err := h.userRepo.GetDeletedUsers(c.Context())
	if err != nil {
		logger.Log.Error("Ошибка при получении удаленных пользователей", zap.Error(err))
		return err
	}

	logger.Log.Info("Удаленные пользователи успешно получены", zap.Int("count", len(users)))

	return SendJSON(c, GetDeletedUsersResponse{
		Status:  true,
		Message: "Удаленные пользователи успешно получены",
		Users:   users,
	})
}

// GetDeletedUsersResponse представляет ответ со списком удаленных пользователей
type GetDeletedUsersResponse struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Users   []*models.User `json:"users"`
}
