package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
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

// CreatePetHandler godoc
// @Summary Создание нового питомца
// @Description Создает нового питомца для пользователя
// @Tags pets
// @Accept json
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Param request body services.PetCreate true "Данные питомца"
// @Success 201 {object} models.Pet "Созданный питомец"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/user/{user_id} [post]
func (h *PetHandler) CreatePetHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID пользователя обязателен",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID пользователя",
		})
	}

	var petData services.PetCreate
	if err := c.BodyParser(&petData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверное тело запроса",
		})
	}

	// Validate pet data
	if err := validation.ValidateStruct(petData); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	logger.Log.Info("создание питомца", zap.Int("userId", userID), zap.String("petName", petData.Name))

	pet, err := h.petService.CreatePet(c.Context(), userID, petData)
	if err != nil {
		logger.Log.Error("не удалось создать питомца", zap.Error(err), zap.Int("userId", userID))

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
	return c.Status(fiber.StatusCreated).JSON(pet)
}

// GetPetHandler godoc
// @Summary Получение питомца по ID
// @Description Возвращает информацию о питомце по его идентификатору
// @Tags pets
// @Produce json
// @Param id path int true "ID питомца"
// @Success 200 {object} models.Pet "Данные питомца"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [get]
func (h *PetHandler) GetPetHandler(c *fiber.Ctx) error {
	petIDStr := c.Params("id")
	if petIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID питомца обязателен",
		})
	}

	petID, err := strconv.Atoi(petIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID питомца",
		})
	}

	logger.Log.Info("получение питомца", zap.Int("petId", petID))

	pet, err := h.petService.GetPetByID(c.Context(), petID)
	if err != nil {
		logger.Log.Error("не удалось получить питомца", zap.Error(err), zap.Int("petId", petID))

		if err.Error() == "питомец не найден" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(pet)
}

// GetUserPetsHandler godoc
// @Summary Получение питомцев пользователя
// @Description Возвращает всех питомцев конкретного пользователя
// @Tags pets
// @Produce json
// @Param user_id path int true "ID пользователя"
// @Success 200 {array} models.Pet "Список питомцев"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/user/{user_id} [get]
func (h *PetHandler) GetUserPetsHandler(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID пользователя обязателен",
		})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID пользователя",
		})
	}

	logger.Log.Info("получение питомцев пользователя", zap.Int("userId", userID))

	pets, err := h.petService.GetUserPets(c.Context(), userID)
	if err != nil {
		logger.Log.Error("не удалось получить питомцев пользователя", zap.Error(err), zap.Int("userId", userID))

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
	return c.JSON(pets)
}

// UpdatePetHandler godoc
// @Summary Обновление данных питомца
// @Description Обновляет информацию о питомце
// @Tags pets
// @Accept json
// @Produce json
// @Param id path int true "ID питомца"
// @Param request body services.PetUpdate true "Данные для обновления"
// @Success 200 {object} SuccessResponse "Данные успешно обновлены"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [put]
func (h *PetHandler) UpdatePetHandler(c *fiber.Ctx) error {
	petIDStr := c.Params("id")
	if petIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID питомца обязателен",
		})
	}

	petID, err := strconv.Atoi(petIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID питомца",
		})
	}

	var updateData services.PetUpdate
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

	logger.Log.Info("обновление питомца", zap.Int("petId", petID))

	if err := h.petService.UpdatePet(c.Context(), petID, updateData); err != nil {
		logger.Log.Error("не удалось обновить питомца", zap.Error(err), zap.Int("petId", petID))

		if err.Error() == "питомец не найден" {
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
		Message: "Питомец успешно обновлен",
	})
}

// DeletePetHandler godoc
// @Summary Удаление питомца по ID
// @Description Удаляет питомца из системы
// @Tags pets
// @Produce json
// @Param id path int true "ID питомца"
// @Success 200 {object} SuccessResponse "Питомец успешно удален"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [delete]
func (h *PetHandler) DeletePetHandler(c *fiber.Ctx) error {
	petIDStr := c.Params("id")
	if petIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID питомца обязателен",
		})
	}

	petID, err := strconv.Atoi(petIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID питомца",
		})
	}

	logger.Log.Info("удаление питомца", zap.Int("petId", petID))

	if err := h.petService.DeletePet(c.Context(), petID); err != nil {
		logger.Log.Error("не удалось удалить питомца", zap.Error(err), zap.Int("petId", petID))

		if err.Error() == "питомец не найден" {
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
		Message: "Питомец успешно удален",
	})
}
