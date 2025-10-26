package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// VetClinicHandler обрабатывает HTTP запросы для операций с ветеринарными клиниками
type VetClinicHandler struct {
	vetClinicService services.VetClinicService
}

// NewVetClinicHandler создает новый обработчик ветеринарных клиник
func NewVetClinicHandler(vetClinicService services.VetClinicService) *VetClinicHandler {
	return &VetClinicHandler{
		vetClinicService: vetClinicService,
	}
}

// RegisterClinicHandler godoc
// @Summary Регистрация новой ветеринарной клиники
// @Description Регистрирует новую ветеринарную клинику в системе
// @Tags vet-clinics
// @Accept json
// @Produce json
// @Param request body services.VetClinicRegistration true "Данные клиники"
// @Success 201 {object} models.VetClinic "Созданная клиника"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 409 {object} ErrorResponse "Клиника уже существует"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /vet-clinics/register [post]
func (h *VetClinicHandler) RegisterClinicHandler(c *fiber.Ctx) error {
	var clinicData services.VetClinicRegistration
	if err := ParseBody(c, &clinicData); err != nil {
		return err
	}

	logger.Log.Info("регистрация ветеринарной клиники", zap.String("clinicName", clinicData.Name))

	clinic, err := h.vetClinicService.RegisterClinic(c.Context(), clinicData)
	if err != nil {
		return err
	}

	return SendCreated(c, clinic)
}

// GetClinicProfileHandler godoc
// @Summary Получение профиля клиники по ID
// @Description Возвращает полный профиль ветеринарной клиники
// @Tags vet-clinics
// @Produce json
// @Param id path int true "ID клиники"
// @Success 200 {object} services.VetClinicProfile "Профиль клиники"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Клиника не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /vet-clinics/{id} [get]
func (h *VetClinicHandler) GetClinicProfileHandler(c *fiber.Ctx) error {
	clinicID, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение профиля клиники", zap.Int("clinicId", clinicID))

	profile, err := h.vetClinicService.GetClinicProfile(c.Context(), clinicID)
	if err != nil {
		return err
	}

	return SendJSON(c, profile)
}

// GetClinicsByLocationIDHandler godoc
// @Summary Получение всех клиник по ID локации
// @Description Возвращает список всех ветеринарных клиник в указанной локации
// @Tags vet-clinics
// @Produce json
// @Param location_id path int true "ID локации"
// @Success 200 {array} models.VetClinic "Список клиник"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /vet-clinics/location/{location_id} [get]
func (h *VetClinicHandler) GetClinicsByLocationIDHandler(c *fiber.Ctx) error {
	locationID, err := ParseIDParam(c, "location_id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение клиник по ID локации", zap.Int("locationId", locationID))

	clinics, err := h.vetClinicService.GetClinicsByLocationID(c.Context(), locationID)
	if err != nil {
		return err
	}

	return SendJSON(c, clinics)
}

// UpdateClinicProfileHandler godoc
// @Summary Обновление профиля клиники
// @Description Обновляет информацию о ветеринарной клинике
// @Tags vet-clinics
// @Accept json
// @Produce json
// @Param id path int true "ID клиники"
// @Param request body services.VetClinicUpdate true "Данные для обновления"
// @Success 200 {object} SuccessResponse "Данные успешно обновлены"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Клиника не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /vet-clinics/{id} [put]
func (h *VetClinicHandler) UpdateClinicProfileHandler(c *fiber.Ctx) error {
	clinicID, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	var updateData services.VetClinicUpdate
	if err := ParseBody(c, &updateData); err != nil {
		return err
	}

	logger.Log.Info("обновление профиля клиники", zap.Int("clinicId", clinicID))

	if err := h.vetClinicService.UpdateClinicProfile(c.Context(), clinicID, updateData); err != nil {
		return err
	}

	return SendSuccess(c, "Профиль клиники успешно обновлен")
}

// DeleteClinicHandler godoc
// @Summary Удаление клиники по ID
// @Description Удаляет клинику из системы (soft delete)
// @Tags vet-clinics
// @Produce json
// @Param id path int true "ID клиники"
// @Success 200 {object} SuccessResponse "Клиника успешно удалена"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Клиника не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /vet-clinics/{id} [delete]
func (h *VetClinicHandler) DeleteClinicHandler(c *fiber.Ctx) error {
	clinicID, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("удаление клиники", zap.Int("clinicId", clinicID))

	if err := h.vetClinicService.DeleteClinic(c.Context(), clinicID); err != nil {
		return err
	}

	return SendSuccess(c, "Клиника успешно удалена")
}
