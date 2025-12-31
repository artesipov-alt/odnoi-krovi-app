package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	bloodrequestv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodrequest/v1"
)

// BloodRequestHandler обрабатывает HTTP запросы для операций с пулом запросов крови
type BloodRequestHandler struct {
	bloodRequestClient services.BloodRequestClient
}

// NewBloodRequestHandler создает новый обработчик для пула запросов крови
func NewBloodRequestHandler(bloodRequestClient services.BloodRequestClient) *BloodRequestHandler {
	return &BloodRequestHandler{
		bloodRequestClient: bloodRequestClient,
	}
}

// AddPetToBloodRequestPool godoc
// @Summary Добавить питомца в пул поиска крови
// @Description Добавляет питомца-реципиента в пул поиска крови
// @Tags pets, blood-request
// @Accept json
// @Produce json
// @Param request body models.BloodSearchPetRequest true "Данные питомца для пула поиска крови"
// @Success 201 {object} models.BloodSearchPetResponse "Статус добавления питомца"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-request/pool [post]
func (h *BloodRequestHandler) AddPetToBloodRequestPool(c *fiber.Ctx) error {
	var petReq models.BloodSearchPetRequest
	if err := utils.ParseBody(c, &petReq); err != nil {
		return err
	}

	// Конвертируем DTO в protobuf структуру
	pet := &bloodrequestv1.BloodRequest{}
	if err := mapstructure.Decode(petReq, pet); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка конвертации данных")
	}

	logger.Log.Info(
		"добавление питомца в пул поиска крови",
		zap.String("petId", pet.PetId),
		zap.String("petType", pet.PetType),
		zap.Strings("bloodGroup", pet.BloodGroup),
	)

	status, err := h.bloodRequestClient.AddPet(c.Context(), pet)
	if err != nil {
		logger.Log.Error("failed to add pet to blood Request pool", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось добавить питомца в пул поиска крови")
	}

	// Конвертируем protobuf ответ в DTO
	var statusResp models.BloodSearchPetResponse
	if err := mapstructure.Decode(status, &statusResp); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка конвертации ответа")
	}

	return utils.SendCreated(c, statusResp)
}

// GetPetsFromBloodRequestPool godoc
// @Summary Получить питомцев из пула поиска крови
// @Description Возвращает список питомцев-реципиентов по фильтрам
// @Tags pets, blood-request
// @Accept json
// @Produce json
// @Param request body models.BloodSearchFilterRequest true "Фильтры поиска: тип, группа крови, регионы"
// @Success 200 {object} models.BloodSearchPetsResponse "Список питомцев"
// @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-request/pool/search [post]
func (h *BloodRequestHandler) GetPetsFromBloodRequestPool(c *fiber.Ctx) error {
	var filterReq models.BloodSearchFilterRequest
	if err := utils.ParseBody(c, &filterReq); err != nil {
		return err
	}

	// Конвертируем DTO в protobuf структуру
	filter := &bloodrequestv1.GetBloodRequests{}
	if err := mapstructure.Decode(filterReq, filter); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Ошибка конвертации фильтров")
	}

	logger.Log.Info(
		"получение питомцев из пула поиска крови",
		zap.String("petId", filterReq.PetID),
		zap.String("petType", filterReq.PetType),
		zap.String("bloodGroup", filterReq.BloodGroup),
		zap.Int("regionsCount", len(filterReq.Regions)),
	)

	pets, err := h.bloodRequestClient.GetPets(c.Context(), filter)
	if err != nil {
		logger.Log.Error("failed to get pets from blood Request pool", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Не удалось получить питомцев из пула поиска крови")
	}

	// Конвертируем protobuf ответ в DTO
	var petsResp models.BloodSearchPetsResponse
	for _, pet := range pets.Pets {
		var dtoPet models.BloodSearchPetRequest
		if err := mapstructure.Decode(pet, &dtoPet); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Ошибка конвертации данных питомца")
		}
		petsResp.Pets = append(petsResp.Pets, dtoPet)
	}

	return utils.SendJSON(c, petsResp)
}
