package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService        services.PetService
	bloodSearchClient services.BloodSearchClient
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService services.PetService, bloodSearchClient services.BloodSearchClient) *PetHandler {
	return &PetHandler{
		petService:        petService,
		bloodSearchClient: bloodSearchClient,
	}
}

// CreatePetHandler godoc
// @Summary Создание нового питомца
// @Description Создает нового питомца для пользователя
// @Tags pets
// @Accept json
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param request body services.PetCreate true "Данные питомца"
// @Success 201 {object} models.Pet "Созданный питомец"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/user/{user_id} [post]
func (h *PetHandler) CreatePetHandler(c *fiber.Ctx) error {
	userID, err := ParseStringParam(c, "user_id")
	if err != nil {
		return err
	}

	var petData services.PetCreate
	if err := ParseBody(c, &petData); err != nil {
		return err
	}

	logger.Log.Info("создание питомца", zap.String("userId", userID), zap.String("petName", petData.Name))

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
// @Param id path string true "ID питомца"
// @Success 200 {object} models.Pet "Данные питомца"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [get]
func (h *PetHandler) GetPetHandler(c *fiber.Ctx) error {
	petID, err := ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение питомца", zap.String("petId", petID))

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
// @Param user_id path string true "ID пользователя"
// @Success 200 {array} models.Pet "Список питомцев"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Пользователь не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/user/{user_id} [get]
func (h *PetHandler) GetUserPetsHandler(c *fiber.Ctx) error {
	userID, err := ParseStringParam(c, "user_id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение питомцев пользователя", zap.String("userId", userID))

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
// @Param id path string true "ID питомца"
// @Param request body services.PetUpdate true "Данные для обновления"
// @Success 200 {object} SuccessResponse "Данные успешно обновлены"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [put]
func (h *PetHandler) UpdatePetHandler(c *fiber.Ctx) error {
	petID, err := ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	var updateData services.PetUpdate
	if err := ParseBody(c, &updateData); err != nil {
		return err
	}

	logger.Log.Info("обновление питомца", zap.String("petId", petID))

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
// @Param id path string true "ID питомца"
// @Success 200 {object} SuccessResponse "Питомец успешно удален"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/{id} [delete]
func (h *PetHandler) DeletePetHandler(c *fiber.Ctx) error {
	petID, err := ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("удаление питомца", zap.String("petId", petID))

	if err := h.petService.DeletePet(c.Context(), petID); err != nil {
		return err
	}

	return SendSuccess(c, "Питомец успешно удален")
}

// AddPetToBloodSearchPoolHandler godoc
// @Summary Добавить питомца в пул поиска крови
// @Description Добавляет питомца-реципиента в пул поиска крови
// @Tags pets, blood-search
// @Accept json
// @Produce json
// @Param request body bloodsearchv1.PetRow true "Данные питомца для пула поиска крови"
// @Success 201 {object} bloodsearchv1.PetRowStatus "Статус добавления питомца"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/blood-search/pool [post]
func (h *PetHandler) AddPetToBloodSearchPoolHandler(c *fiber.Ctx) error {
	var pet bloodsearchv1.PetRow
	if err := ParseBody(c, &pet); err != nil {
		return err
	}

	logger.Log.Info(
		"добавление питомца в пул поиска крови",
		zap.String("petId", pet.PetId),
		zap.String("petType", pet.PetType),
		zap.String("bloodGroup", pet.BloodGroup),
	)

	status, err := h.bloodSearchClient.AddPet(c.Context(), &pet)
	if err != nil {
		logger.Log.Error("failed to add pet to blood search pool", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось добавить питомца в пул поиска крови")
	}

	return SendCreated(c, status)
}

// GetPetsFromBloodSearchPoolHandler godoc
// @Summary Получить питомцев из пула поиска крови
// @Description Возвращает список питомцев-реципиентов по фильтрам
// @Tags pets, blood-search
// @Accept json
// @Produce json
// @Param request body bloodsearchv1.GetPetRows true "Фильтры поиска: тип, группа крови, регионы"
// @Success 200 {object} bloodsearchv1.PetRows "Список питомцев"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/blood-search/pool/search [post]
func (h *PetHandler) GetPetsFromBloodSearchPoolHandler(c *fiber.Ctx) error {
	var filter bloodsearchv1.GetPetRows
	if err := ParseBody(c, &filter); err != nil {
		return err
	}

	logger.Log.Info(
		"получение питомцев из пула поиска крови",
		zap.String("petType", filter.PetType),
		zap.String("bloodGroup", filter.BloodGroup),
		zap.Int("regionsCount", len(filter.Regions)),
	)

	pets, err := h.bloodSearchClient.GetPets(c.Context(), &filter)
	if err != nil {
		logger.Log.Error("failed to get pets from blood search pool", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось получить питомцев из пула поиска крови")
	}

	return SendJSON(c, pets)
}
