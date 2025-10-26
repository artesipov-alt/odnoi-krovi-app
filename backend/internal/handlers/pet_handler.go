package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
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
	userID, err := ParseIDParam(c, "user_id")
	if err != nil {
		return err
	}

	var petData services.PetCreate
	if err := ParseBody(c, &petData); err != nil {
		return err
	}

	logger.Log.Info("создание питомца", zap.Int("userId", userID), zap.String("petName", petData.Name))

	pet, err := h.petService.CreatePet(c.Context(), userID, petData)
	if err != nil {
		return err
	}

	return SendCreated(c, pet)
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
	petID, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение питомца", zap.Int("petId", petID))

	pet, err := h.petService.GetPetByID(c.Context(), petID)
	if err != nil {
		return err
	}

	return SendJSON(c, pet)
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
	userID, err := ParseIDParam(c, "user_id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение питомцев пользователя", zap.Int("userId", userID))

	pets, err := h.petService.GetUserPets(c.Context(), userID)
	if err != nil {
		return err
	}

	return SendJSON(c, pets)
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
	petID, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	var updateData services.PetUpdate
	if err := ParseBody(c, &updateData); err != nil {
		return err
	}

	logger.Log.Info("обновление питомца", zap.Int("petId", petID))

	if err := h.petService.UpdatePet(c.Context(), petID, updateData); err != nil {
		return err
	}

	return SendSuccess(c, "Питомец успежно обновлен")
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
	petID, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("удаление питомца", zap.Int("petId", petID))

	if err := h.petService.DeletePet(c.Context(), petID); err != nil {
		return err
	}

	return SendSuccess(c, "Питомец успешно удален")
}
