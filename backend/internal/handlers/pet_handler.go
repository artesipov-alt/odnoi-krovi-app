package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService services.PetService
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService services.PetService) *PetHandler {
	return &PetHandler{
		petService: petService,
	}
}

// getPreloads извлекает список связей для предзагрузки из query-параметров
func (h *PetHandler) getPreloads(c echo.Context) []string {
	var preloads []string
	if c.QueryParam("with_health") == "true" {
		preloads = append(preloads, "Health")
	}
	if c.QueryParam("with_treatments") == "true" {
		preloads = append(preloads, "Treatments")
	}
	if c.QueryParam("with_analysis") == "true" {
		preloads = append(preloads, "Analysis")
	}
	if c.QueryParam("with_bonuses") == "true" {
		preloads = append(preloads, "Bonuses")
	}
	if c.QueryParam("with_all") == "true" {
		return []string{"Health", "Treatments", "Analysis", "Bonuses"}
	}
	return preloads
}

// CreatePetHandler godoc
// @Summary Создание нового питомца
// @Description Создает нового питомца для пользователя
// @Tags pets
// @Accept json
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param request body dto.PetCreate true "Данные питомца"
// @Success 201 {object} ent.Pet "Созданный питомец"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/user/{user_id} [post]
func (h *PetHandler) CreatePetHandler(c echo.Context) error {
	userID, err := utils.ParseStringParam(c, "user_id")
	if err != nil {
		return err
	}

	var petData dto.PetCreate
	if err := utils.ParseBody(c, &petData); err != nil {
		return err
	}

	logger.Log.Info("создание питомца", zap.String("userId", userID), zap.String("petName", petData.Name))

	pet, err := h.petService.CreatePet(c.Request().Context(), userID, petData)
	if err != nil {
		return err
	}

	return utils.SendCreated(c, pet)
}

// GetPetHandler godoc
// @Summary Получение питомца по ID
// @Description Возвращает информацию о питомце по его идентификатору
// @Tags pets
// @Produce json
// @Param id path string true "ID питомца"
// @Param with_health query bool false "Включить данные о здоровье"
// @Param with_treatments query bool false "Включить данные о ветеринарных обработках"
// @Param with_analysis query bool false "Включить данные об анализах"
// @Param with_bonuses query bool false "Включить данные о бонусах"
// @Param with_all query bool false "Включить все связанные данные"
// @Success 200 {object} ent.Pet "Данные питомца"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 404 {object} utils.ErrorResponse "Питомец не найден"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [get]
func (h *PetHandler) GetPetHandler(c echo.Context) error {
	petID, err := utils.ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	preloads := h.getPreloads(c)

	logger.Log.Info("получение питомца", zap.String("petId", petID), zap.Strings("preloads", preloads))

	pet, err := h.petService.GetPetByID(c.Request().Context(), petID, preloads...)
	if err != nil {
		return err
	}

	return utils.SendJSON(c, pet)
}

// GetUserPetsHandler godoc
// @Summary Получение питомцев пользователя
// @Description Возвращает всех питомцев конкретного пользователя
// @Tags pets
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param with_health query bool false "Включить данные о здоровье"
// @Param with_treatments query bool false "Включить данные о ветеринарных обработках"
// @Param with_analysis query bool false "Включить данные об анализах"
// @Param with_bonuses query bool false "Включить данные о бонусах"
// @Param with_all query bool false "Включить все связанные данные"
// @Success 200 {array} ent.Pet "Список питомцев"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 404 {object} utils.ErrorResponse "Пользователь не найден"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/user/{user_id} [get]
func (h *PetHandler) GetUserPetsHandler(c echo.Context) error {
	userID, err := utils.ParseStringParam(c, "user_id")
	if err != nil {
		return err
	}

	preloads := h.getPreloads(c)

	logger.Log.Info("получение питомцев пользователя", zap.String("userId", userID), zap.Strings("preloads", preloads))

	pets, err := h.petService.GetUserPets(c.Request().Context(), userID, preloads...)
	if err != nil {
		return err
	}

	return utils.SendJSON(c, pets)
}

// UpdatePetHandler godoc
// @Summary Обновление данных питомца
// @Description Обновляет информацию о питомце
// @Tags pets
// @Accept json
// @Produce json
// @Param id path string true "ID питомца"
// @Param request body dto.PetUpdate true "Данные для обновления"
// @Success 200 {object} utils.SuccessResponse "Данные успешно обновлены"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 404 {object} utils.ErrorResponse "Питомец не найден"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [put]
func (h *PetHandler) UpdatePetHandler(c echo.Context) error {
	petID, err := utils.ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	var updateData dto.PetUpdate
	if err := utils.ParseBody(c, &updateData); err != nil {
		return err
	}

	logger.Log.Info("обновление питомца", zap.String("petId", petID))

	if err := h.petService.UpdatePet(c.Request().Context(), petID, updateData); err != nil {
		return err
	}

	return utils.SendSuccess(c, "Питомец успешно обновлен")
}

// DeletePetHandler godoc
// @Summary Удаление питомца по ID
// @Description Удаляет питомца из системы
// @Tags pets
// @Produce json
// @Param id path string true "ID питомца"
// @Success 200 {object} utils.SuccessResponse "Питомец успешно удален"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 404 {object} utils.ErrorResponse "Питомец не найден"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [delete]
func (h *PetHandler) DeletePetHandler(c echo.Context) error {
	petID, err := utils.ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("удаление питомца", zap.String("petId", petID))

	if err := h.petService.DeletePet(c.Request().Context(), petID); err != nil {
		return err
	}

	return utils.SendSuccess(c, "Питомец успешно удален")
}

// GetAvatarUploadURL godoc
// @Summary Получить ссылку для загрузки фотографии питомца
// @Description Возвращает временную ссылку для загрузки фотографии питомца по ID
// @Tags pets
// @Produce json
// @Param id path string true "ID питомца"
// @Success 200 {object} map[string]string "Ссылка для загрузки фотографии и путь к файлу"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 404 {object} utils.ErrorResponse "Питомец не найден"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/upload/avatar/{id} [get]
func (h *PetHandler) GetAvatarUploadURL(c echo.Context) error {
	petID, err := utils.ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение ссылки для загрузки фотографии питомца", zap.String("petId", petID))

	url, path, err := h.petService.GetAvatarUploadURL(c.Request().Context(), petID)
	if err != nil {
		return err
	}

	return utils.SendJSON(c, map[string]string{"url": url, "path": path})
}

// ConfirmPetAvatarUpload godoc
// @Summary Подтверждение загрузки аватарки питомца
// @Description Подтверждает загрузку аватарки питомца, делает её публичной и возвращает публичную ссылку
// @Tags pets
// @Produce json
// @Param path path string true "Путь к аватарке питомца (например: pets/PET-25-000001/avatar.jpg)"
// @Success 200 {object} map[string]string "Публичная ссылка на аватарку"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 404 {object} utils.ErrorResponse "Питомец не найден"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/upload/avatar/confirm/{path} [post]
func (h *PetHandler) ConfirmPetAvatarUpload(c echo.Context) error {
	avatarPath, err := utils.ParseStringParam(c, "path")
	if err != nil {
		return err
	}

	logger.Log.Info("подтверждение загрузки аватарки питомца", zap.String("avatarPath", avatarPath))

	publicURL, err := h.petService.UpdatePetAvatar(c.Request().Context(), avatarPath)
	if err != nil {
		return err
	}

	return utils.SendJSON(c, map[string]string{"publicUrl": publicURL})
}
