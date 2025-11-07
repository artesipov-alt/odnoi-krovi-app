package handlers

import (
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
