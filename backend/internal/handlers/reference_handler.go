package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ReferenceHandler обрабатывает HTTP запросы для справочных данных
type ReferenceHandler struct {
	breedRepo     repositories.BreedRepository
	bloodTypeRepo repositories.BloodTypeRepository
}

// NewReferenceHandler создает новый обработчик справочных данных
func NewReferenceHandler(breedRepo repositories.BreedRepository, bloodTypeRepo repositories.BloodTypeRepository) *ReferenceHandler {
	return &ReferenceHandler{
		breedRepo:     breedRepo,
		bloodTypeRepo: bloodTypeRepo,
	}
}

// ReferenceResponse представляет ответ со справочными данными
type ReferenceResponse struct {
	Data []ReferenceItem `json:"data"`
}

// ReferenceItem представляет элемент справочника
type ReferenceItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// GetPetTypesHandler godoc
// @Summary Получение всех типов животных
// @Description Возвращает все доступные типы животных для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список типов животных"
// @Router /reference/pet-types [get]
func (h *ReferenceHandler) GetPetTypesHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника типов животных")

	petTypes := models.GetAllPetTypes()
	items := make([]ReferenceItem, len(petTypes))

	for i, petType := range petTypes {
		items[i] = ReferenceItem{
			Value: string(petType),
			Label: petType.String(),
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetGendersHandler godoc
// @Summary Получение всех значений пола
// @Description Возвращает все доступные значения пола для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список значений пола"
// @Router /reference/genders [get]
func (h *ReferenceHandler) GetGendersHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника полов")

	genders := models.GetAllGenders()
	items := make([]ReferenceItem, len(genders))

	for i, gender := range genders {
		items[i] = ReferenceItem{
			Value: string(gender),
			Label: gender.String(),
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetLivingConditionsHandler godoc
// @Summary Получение всех условий проживания
// @Description Возвращает все доступные условия проживания для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список условий проживания"
// @Router /reference/living-conditions [get]
func (h *ReferenceHandler) GetLivingConditionsHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника условий проживания")

	conditions := models.GetAllLivingConditions()
	items := make([]ReferenceItem, len(conditions))

	for i, condition := range conditions {
		items[i] = ReferenceItem{
			Value: string(condition),
			Label: condition.String(),
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetUserRolesHandler godoc
// @Summary Получение всех ролей пользователей
// @Description Возвращает все доступные роли пользователей для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список ролей пользователей"
// @Router /reference/user-roles [get]
func (h *ReferenceHandler) GetUserRolesHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника ролей пользователей")

	roles := models.GetAllUserRoles()
	items := make([]ReferenceItem, len(roles))

	for i, role := range roles {
		items[i] = ReferenceItem{
			Value: string(role),
			Label: role.String(),
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetBloodSearchStatusesHandler godoc
// @Summary Получение всех статусов поиска крови
// @Description Возвращает все доступные статусы поиска крови для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список статусов поиска крови"
// @Router /reference/blood-search-statuses [get]
func (h *ReferenceHandler) GetBloodSearchStatusesHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника статусов поиска крови")

	statuses := models.GetAllBloodSearchStatuses()
	items := make([]ReferenceItem, len(statuses))

	for i, status := range statuses {
		items[i] = ReferenceItem{
			Value: string(status),
			Label: status.String(),
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetBloodStockStatusesHandler godoc
// @Summary Получение всех статусов запаса крови
// @Description Возвращает все доступные статусы запаса крови для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список статусов запаса крови"
// @Router /reference/blood-stock-statuses [get]
func (h *ReferenceHandler) GetBloodStockStatusesHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника статусов запаса крови")

	statuses := models.GetAllBloodStockStatuses()
	items := make([]ReferenceItem, len(statuses))

	for i, status := range statuses {
		items[i] = ReferenceItem{
			Value: string(status),
			Label: status.String(),
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetDonationStatusesHandler godoc
// @Summary Получение всех статусов донорства
// @Description Возвращает все доступные статусы донорства для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список статусов донорства"
// @Router /reference/donation-statuses [get]
func (h *ReferenceHandler) GetDonationStatusesHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника статусов донорства")

	statuses := models.GetAllDonationStatuses()
	items := make([]ReferenceItem, len(statuses))

	for i, status := range statuses {
		items[i] = ReferenceItem{
			Value: string(status),
			Label: status.String(),
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetBreedsHandler godoc
// @Summary Получение всех пород животных
// @Description Возвращает список всех пород животных в базе для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список пород животных"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/breeds [get]
func (h *ReferenceHandler) GetBreedsHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника пород животных")

	breeds, err := h.breedRepo.GetAll(c.Context())
	if err != nil {
		logger.Log.Error("не удалось получить породы из БД", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось получить список пород",
		})
	}

	items := make([]ReferenceItem, len(breeds))
	for i, breed := range breeds {
		items[i] = ReferenceItem{
			Value: breed.Name,
			Label: breed.Name,
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetBreedsByTypeHandler godoc
// @Summary Получение пород животных по типу
// @Description Возвращает список пород животных для указанного типа животного для выбора на фронтенде
// @Tags reference
// @Produce json
// @Param petType query string true "Тип животного (dog, cat, etc.)"
// @Success 200 {object} ReferenceResponse "Список пород животных"
// @Failure 400 {object} ErrorResponse "Неверный тип животного"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/breeds-by-type [get]
func (h *ReferenceHandler) GetBreedsByTypeHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника пород животных по типу животного")

	petTypeStr := c.Query("petType")
	if petTypeStr == "" {
		logger.Log.Error("не указан тип животного")
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Необходимо указать тип животного",
		})
	}

	petType, err := models.ValidatePetType(petTypeStr)
	if err != nil {
		logger.Log.Error("неверный тип животного", zap.String("petType", petTypeStr), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный тип животного",
		})
	}

	breeds, err := h.breedRepo.GetByPetType(c.Context(), petType)
	if err != nil {
		logger.Log.Error("не удалось получить породы из БД", zap.Error(err), zap.String("petType", petTypeStr))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось получить список пород",
		})
	}

	items := make([]ReferenceItem, len(breeds))
	for i, breed := range breeds {
		items[i] = ReferenceItem{
			Value: breed.Name,
			Label: breed.Name,
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}

// GetBloodGroupsHandler godoc
// @Summary Получение групп крови животных
// @Description Возвращает список групп крови животных для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список групп крови"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/blood-groups [get]
func (h *ReferenceHandler) GetBloodGroupsHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение справочника групп крови")

	bloodTypes, err := h.bloodTypeRepo.GetAll(c.Context())
	if err != nil {
		logger.Log.Error("не удалось получить группы крови из БД", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Не удалось получить список групп крови",
		})
	}

	items := make([]ReferenceItem, len(bloodTypes))
	for i, bloodType := range bloodTypes {
		items[i] = ReferenceItem{
			Value: bloodType.Name,
			Label: bloodType.Name,
		}
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(ReferenceResponse{Data: items})
}
