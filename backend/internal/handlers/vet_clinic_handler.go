package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
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
	if err := c.BodyParser(&clinicData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверное тело запроса",
		})
	}

	// Validate clinic data
	if err := validation.ValidateStruct(clinicData); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	logger.Log.Info("регистрация ветеринарной клиники", zap.String("clinicName", clinicData.Name))

	clinic, err := h.vetClinicService.RegisterClinic(c.Context(), clinicData)
	if err != nil {
		logger.Log.Error("не удалось зарегистрировать клинику", zap.Error(err))

		if err.Error() == "клиника уже существует" {
			return c.Status(fiber.StatusConflict).JSON(ErrorResponse{
				Error: "Клиника уже существует",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось зарегистрировать клинику",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Status(fiber.StatusCreated).JSON(clinic)
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
	clinicIDStr := c.Params("id")
	if clinicIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID клиники обязателен",
		})
	}

	clinicID, err := strconv.Atoi(clinicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID клиники",
		})
	}

	logger.Log.Info("получение профиля клиники", zap.Int("clinicId", clinicID))

	profile, err := h.vetClinicService.GetClinicProfile(c.Context(), clinicID)
	if err != nil {
		logger.Log.Error("не удалось получить профиль клиники", zap.Error(err), zap.Int("clinicId", clinicID))

		if err.Error() == "клиника не найдена" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Клиника не найдена",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось получить профиль клиники",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(profile)
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
	locationIDStr := c.Params("location_id")
	if locationIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID локации обязателен",
		})
	}

	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID локации",
		})
	}

	logger.Log.Info("получение клиник по ID локации", zap.Int("locationId", locationID))

	clinics, err := h.vetClinicService.GetClinicsByLocationID(c.Context(), locationID)
	if err != nil {
		logger.Log.Error("не удалось получить клиники", zap.Error(err), zap.Int("locationId", locationID))

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось получить клиники",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(clinics)
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
	clinicIDStr := c.Params("id")
	if clinicIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID клиники обязателен",
		})
	}

	clinicID, err := strconv.Atoi(clinicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID клиники",
		})
	}

	var updateData services.VetClinicUpdate
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

	logger.Log.Info("обновление профиля клиники", zap.Int("clinicId", clinicID))

	if err := h.vetClinicService.UpdateClinicProfile(c.Context(), clinicID, updateData); err != nil {
		logger.Log.Error("не удалось обновить профиль клиники", zap.Error(err), zap.Int("clinicId", clinicID))

		if err.Error() == "клиника не найдена" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Клиника не найдена",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось обновить профиль клиники",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(SuccessResponse{
		Message: "Профиль клиники успешно обновлен",
	})
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
	clinicIDStr := c.Params("id")
	if clinicIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID клиники обязателен",
		})
	}

	clinicID, err := strconv.Atoi(clinicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID клиники",
		})
	}

	logger.Log.Info("удаление клиники", zap.Int("clinicId", clinicID))

	if err := h.vetClinicService.DeleteClinic(c.Context(), clinicID); err != nil {
		logger.Log.Error("не удалось удалить клинику", zap.Error(err), zap.Int("clinicId", clinicID))

		if err.Error() == "клиника не найдена" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Клиника не найдена",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось удалить клинику",
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(SuccessResponse{
		Message: "Клиника успешно удалена",
	})
}
