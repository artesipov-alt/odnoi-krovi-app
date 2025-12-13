package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	bloodrequestv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodrequest/v1"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	petService         services.PetService
	bloodRequestClient services.BloodRequestClient
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(petService services.PetService, bloodRequestClient services.BloodRequestClient) *PetHandler {
	return &PetHandler{
		petService:         petService,
		bloodRequestClient: bloodRequestClient,
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

	return SendSuccess(c, "Питомец успешно обновлен")
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

// AddPetToBloodRequestPoolHandler godoc
// @Summary Добавить питомца в пул поиска крови
// @Description Добавляет питомца-реципиента в пул поиска крови
// @Tags pets, blood-Request
// @Accept json
// @Produce json
// @Param request body BloodRequestPetRequest true "Данные питомца для пула поиска крови"
// @Success 201 {object} BloodRequestPetResponse "Статус добавления питомца"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/blood-request/pool [post]
func (h *PetHandler) AddPetToBloodRequestPool(c *fiber.Ctx) error {
	var petReq models.BloodSearchPetRequest
	if err := ParseBody(c, &petReq); err != nil {
		return err
	}

	// Конвертируем DTO в protobuf структуру
	pet := &bloodrequestv1.BloodRequest{
		PetId:                  petReq.PetID,
		PetType:                petReq.PetType,
		BloodGroup:             petReq.BloodGroup,
		BloodComponents:        petReq.BloodComponents,
		BloodVolumeNeeded:      petReq.BloodVolumeNeeded,
		BloodVolumeReserved:    petReq.BloodVolumeReserved,
		Regions:                petReq.Regions,
		SmallPetsNotifyAllowed: petReq.SmallPetsNotifyAllowed,
		Status:                 petReq.Status,
		Description:            petReq.Description,
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
	statusResp := models.BloodSearchPetResponse{
		PetID:  status.PetId,
		Status: status.Status,
	}

	return SendCreated(c, statusResp)
}

// GetPetsFromBloodRequestPoolHandler godoc
// @Summary Получить питомцев из пула поиска крови
// @Description Возвращает список питомцев-реципиентов по фильтрам
// @Tags pets, blood-Request
// @Accept json
// @Produce json
// @Param request body BloodSearchFilterRequest true "Фильтры поиска: тип, группа крови, регионы"
// @Success 200 {object} BloodSearchPetsResponse "Список питомцев"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/blood-Request/pool/Request [post]
func (h *PetHandler) GetPetsFromBloodRequestPool(c *fiber.Ctx) error {
	var filterReq models.BloodSearchFilterRequest
	if err := ParseBody(c, &filterReq); err != nil {
		return err
	}

	// Конвертируем DTO в protobuf структуру
	filter := &bloodrequestv1.GetBloodRequests{
		PetId:      filterReq.PetID,
		PetType:    filterReq.PetType,
		BloodGroup: filterReq.BloodGroup,
		Regions:    filterReq.Regions,
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
		petsResp.Pets = append(petsResp.Pets, models.BloodSearchPetRequest{
			PetID:                  pet.PetId,
			PetType:                pet.PetType,
			BloodGroup:             pet.BloodGroup,
			BloodComponents:        pet.BloodComponents,
			BloodVolumeNeeded:      pet.BloodVolumeNeeded,
			BloodVolumeReserved:    pet.BloodVolumeReserved,
			Regions:                pet.Regions,
			SmallPetsNotifyAllowed: pet.SmallPetsNotifyAllowed,
			Status:                 pet.Status,
			Description:            pet.Description,
		})
	}

	return SendJSON(c, petsResp)
}

// GetAvatarUploadURL godoc
// @Summary Получить ссылку для загрузки фотографии питомца
// @Description Возвращает временную ссылку для загрузки фотографии питомца по ID
// @Tags pets
// @Produce json
// @Param id path string true "ID питомца"
// @Success 200 {object} map[string]string "Ссылка для загрузки фотографии и путь к файлу"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/upload/avatar/{id} [get]
func (h *PetHandler) GetAvatarUploadURL(c *fiber.Ctx) error {
	petID, err := ParseStringParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение ссылки для загрузки фотографии питомца", zap.String("petId", petID))

	url, path, err := h.petService.GetAvatarUploadURL(c.Context(), petID)
	if err != nil {
		return err
	}

	return SendJSON(c, map[string]string{"url": url, "path": path})
}

// ConfirmPetAvatarUpload godoc
// @Summary Подтверждение загрузки аватарки питомца
// @Description Подтверждает загрузку аватарки питомца, делает её публичной и возвращает публичную ссылку
// @Tags pets
// @Produce json
// @Param path path string true "Путь к аватарке питомца (например: pets/PET-25-000001/avatar.jpg)"
// @Success 200 {object} map[string]string "Публичная ссылка на аватарку"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Питомец не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /pets/upload/avatar/confirm/{path} [post]
func (h *PetHandler) ConfirmPetAvatarUpload(c *fiber.Ctx) error {
	avatarPath, err := ParseStringParam(c, "path")
	if err != nil {
		return err
	}

	logger.Log.Info("подтверждение загрузки аватарки питомца", zap.String("avatarPath", avatarPath))

	publicURL, err := h.petService.UpdatePetAvatar(c.Context(), avatarPath)
	if err != nil {
		return err
	}

	return SendJSON(c, map[string]string{"publicUrl": publicURL})
}
