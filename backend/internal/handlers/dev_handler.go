package handlers

import (
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/labstack/echo/v4"
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
// @Param id path string true "ID пользователя"
// @Success 200 {object} DevResponse "Успешный сброс пользователя"
// @Router /dev/reset-user/{id} [post]
func (h *DevHandler) ResetUserHandler(c echo.Context) error {
	logger.Log.Info("Сброс пользователя к заводским настройкам")

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	logger.Log.Info("Сброс пользователя", zap.String("userId", id))

	if err := h.userRepo.ResetUser(c.Request().Context(), id); err != nil {
		logger.Log.Error("Ошибка при сбросе пользователя", zap.Error(err), zap.String("userId", id))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	logger.Log.Info("Пользователь успешно сброшен", zap.String("userId", id))

	return c.JSON(http.StatusOK, DevResponse{
		Status:  true,
		Message: "Пользователь успешно сброшен к заводским настройкам",
	})
}

// RestoreUserHandler godoc
// @Summary Восстановление удаленного пользователя
// @Description Восстанавливает мягко удаленного пользователя, устанавливая deleted_at в NULL
// @Tags dev
// @Produce json
// @Param id path string true "ID пользователя"
// @Success 200 {object} DevResponse "Успешное восстановление пользователя"
// @Router /dev/restore-user/{id} [post]
func (h *DevHandler) RestoreUserHandler(c echo.Context) error {
	logger.Log.Info("Восстановление удаленного пользователя")

	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	logger.Log.Info("Восстановление пользователя", zap.String("userId", id))

	if err := h.userRepo.RestoreUser(c.Request().Context(), id); err != nil {
		logger.Log.Error("Ошибка при восстановлении пользователя", zap.Error(err), zap.String("userId", id))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	logger.Log.Info("Пользователь успешно восстановлен", zap.String("userId", id))

	return c.JSON(http.StatusOK, DevResponse{
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
func (h *DevHandler) GetDeletedUsersHandler(c echo.Context) error {
	logger.Log.Info("Получение списка удаленных пользователей")

	users, err := h.userRepo.GetDeletedUsers(c.Request().Context())
	if err != nil {
		logger.Log.Error("Ошибка при получении удаленных пользователей", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	logger.Log.Info("Удаленные пользователи успешно получены", zap.Int("count", len(users)))

	return c.JSON(http.StatusOK, GetDeletedUsersResponse{
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
