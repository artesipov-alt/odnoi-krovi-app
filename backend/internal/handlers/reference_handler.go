package handlers

import (
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/enums"
	validation "github.com/artesipov-alt/odnoi-krovi-app/internal/utils/enums"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ReferenceHandler обрабатывает HTTP запросы для справочных данных
type ReferenceHandler struct {
	breedRepo    repositories.BreedRepository
	bloodRepo    repositories.BloodInfoRepository
	locationRepo repositories.LocationRepository
}

// NewReferenceHandler создает новый обработчик справочных данных
func NewReferenceHandler(breedRepo repositories.BreedRepository, bloodTypeRepo repositories.BloodInfoRepository, locationRepo repositories.LocationRepository) *ReferenceHandler {
	return &ReferenceHandler{
		breedRepo:    breedRepo,
		bloodRepo:    bloodTypeRepo,
		locationRepo: locationRepo,
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

// ReferenceResponseDB представляет ответ со справочными данными из базы данных
type ReferenceResponseDB struct {
	Data []ReferenceItemDB `json:"data"`
}

// ReferenceItemDB представляет элемент справочника из базы данных с ID
type ReferenceItemDB struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// GetPetTypesHandler godoc
// @Summary Получение всех типов животных
// @Description Возвращает все доступные типы животных для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список типов животных"
// @Router /reference/pet-types [get]
func (h *ReferenceHandler) GetPetTypesHandler(c echo.Context) error {
	logger.Log.Info("получение справочника типов животных")

	petTypes := enums.GetAllPetTypes()
	items := make([]ReferenceItem, len(petTypes))

	for i, petType := range petTypes {
		ruValue, err := validation.LocalizePetType(string(petType))
		if err != nil {
			return c.JSON(500, utils.ErrorResponse{
				Message: apperrors.ErrInvalidPetType.Error(),
			})
		}

		items[i] = ReferenceItem{
			Value: string(petType),
			Label: string(ruValue),
		}
	}

	return c.JSON(200, ReferenceResponse{Data: items})
}

// GetGendersHandler godoc
// @Summary Получение всех значений пола
// @Description Возвращает все доступные значения пола для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список значений пола"
// @Router /reference/genders [get]
func (h *ReferenceHandler) GetGendersHandler(c echo.Context) error {
	logger.Log.Info("получение справочника полов")

	genders := enums.GetAllGenders()
	items := make([]ReferenceItem, len(genders))

	for i, gender := range genders {
		ruValue, err := validation.LocalizeGender(string(gender))
		if err != nil {
			return c.JSON(500, utils.ErrorResponse{
				Message: apperrors.ErrInvalidGender.Error(),
			})
		}

		items[i] = ReferenceItem{
			Value: string(gender),
			Label: string(ruValue),
		}
	}

	return c.JSON(200, ReferenceResponse{Data: items})
}

// GetLivingConditionsHandler godoc
// @Summary Получение всех условий проживания
// @Description Возвращает все доступные условия проживания для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список условий проживания"
// @Router /reference/living-conditions [get]
func (h *ReferenceHandler) GetLivingConditionsHandler(c echo.Context) error {
	logger.Log.Info("получение справочника условий проживания")

	conditions := enums.GetAllLivingConditions()
	items := make([]ReferenceItem, len(conditions))

	for i, condition := range conditions {
		ruValue, err := validation.LocalizeLivingCondition(string(condition))
		if err != nil {
			return c.JSON(500, utils.ErrorResponse{
				Message: apperrors.ErrInvalidLivingCondition.Error(),
			})
		}

		items[i] = ReferenceItem{
			Value: string(condition),
			Label: string(ruValue),
		}
	}

	return c.JSON(200, ReferenceResponse{Data: items})
}

// GetUserRolesHandler godoc
// @Summary Получение всех ролей пользователей
// @Description Возвращает все доступные роли пользователей для выбора на фронтенде
// @Tags reference, users
// @Produce json
// @Success 200 {object} ReferenceResponse "Список ролей пользователей"
// @Router /reference/user-roles [get]
func (h *ReferenceHandler) GetUserRolesHandler(c echo.Context) error {
	logger.Log.Info("получение справочника ролей пользователей")

	roles := enums.GetAllUserRoles()
	items := make([]ReferenceItem, len(roles))

	for i, role := range roles {
		ruValue, err := validation.LocalizeUserRole(string(role))
		if err != nil {
			return c.JSON(500, utils.ErrorResponse{
				Message: apperrors.ErrUserInvalidRole.Error(),
			})
		}

		items[i] = ReferenceItem{
			Value: string(role),
			Label: string(ruValue),
		}
	}

	return c.JSON(200, ReferenceResponse{Data: items})
}

// GetPetRolesHandler godoc
// @Summary Получение всех ролей питомцев
// @Description Возвращает все доступные роли питомцев для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список ролей питомцев"
// @Router /reference/pet-roles [get]
func (h *ReferenceHandler) GetPetRolesHandler(c echo.Context) error {
	logger.Log.Info("получение справочника ролей питомцев")

	roles := enums.GetAllEntPetStatuses()
	items := make([]ReferenceItem, len(roles))

	for i, role := range roles {
		ruValue := enums.LocalizeEntPetStatus(role)

		items[i] = ReferenceItem{
			Value: string(role),
			Label: ruValue,
		}
	}

	return c.JSON(200, ReferenceResponse{Data: items})
}

// GetBreedsHandler godoc
// @Summary Получение всех пород животных
// @Description Возвращает список всех пород животных в базе для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список пород животных"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/breeds [get]
func (h *ReferenceHandler) GetBreedsHandler(c echo.Context) error {
	logger.Log.Info("получение справочника пород животных")

	breeds, err := h.breedRepo.GetAll(c.Request().Context())
	if err != nil {
		logger.Log.Error("не удалось получить породы из БД", zap.Error(err))
		return c.JSON(500, utils.ErrorResponse{
			Message: "Не удалось получить список пород",
		})
	}

	items := make([]ReferenceItemDB, len(breeds))
	for i, breed := range breeds {
		items[i] = ReferenceItemDB{
			Value: breed.ID,
			Label: breed.Name,
		}
	}

	return c.JSON(200, ReferenceResponseDB{Data: items})
}

// GetLocationsHandler godoc
// @Summary Получение всех локаций
// @Description Возвращает список всех локаций в системе для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponseDB "Список локаций"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/locations [get]
func (h *ReferenceHandler) GetLocationsHandler(c echo.Context) error {
	logger.Log.Info("получение справочника локаций")

	locations, err := h.locationRepo.GetAll(c.Request().Context())
	if err != nil {
		logger.Log.Error("не удалось получить локации из БД", zap.Error(err))
		return c.JSON(500, utils.ErrorResponse{
			Message: "Не удалось получить список локаций",
		})
	}

	items := make([]ReferenceItemDB, len(locations))
	for i, location := range locations {
		items[i] = ReferenceItemDB{
			Value: location.ID,
			Label: location.Name,
		}
	}

	return c.JSON(200, ReferenceResponseDB{Data: items})
}

// GetBreedsByTypeHandler godoc
// @Summary Получение пород животных по типу
// @Description Возвращает список пород животных для указанного типа животного для выбора на фронтенде
// @Tags reference
// @Produce json
// @Param petType query string true "Тип животного (dog, cat, etc.)"
// @Success 200 {object} ReferenceResponse "Список пород животных"
// @Failure 400 {object} utils.ErrorResponse "Неверный тип животного"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/breeds-by-type [get]
func (h *ReferenceHandler) GetBreedsByTypeHandler(c echo.Context) error {
	logger.Log.Info("получение справочника пород животных по типу животного")

	petTypeStr := c.QueryParam("petType")
	if petTypeStr == "" {
		logger.Log.Error("не указан тип животного")
		return c.JSON(400, utils.ErrorResponse{
			Message: "Необходимо указать тип животного",
		})
	}

	if _, err := validation.LocalizePetType(petTypeStr); err != nil {
		logger.Log.Error("неверный тип животного", zap.String("petType", petTypeStr), zap.Error(err))
		return c.JSON(400, utils.ErrorResponse{
			Message: "Неверный тип животного",
		})
	}

	breeds, err := h.breedRepo.GetByPetType(c.Request().Context(), breed.Type(petTypeStr))
	if err != nil {
		logger.Log.Error("не удалось получить породы из БД", zap.Error(err), zap.String("petType", petTypeStr))
		return c.JSON(500, utils.ErrorResponse{
			Message: "Не удалось получить список пород",
		})
	}

	items := make([]ReferenceItemDB, len(breeds))
	for i, breed := range breeds {
		items[i] = ReferenceItemDB{
			Value: breed.ID,
			Label: breed.Name,
		}
	}

	return c.JSON(200, ReferenceResponseDB{Data: items})
}

// GetBloodComponentsHandler godoc
// @Summary Получение компонентов крови животных
// @Description Возвращает список компонентов крови животных для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список компонентов крови"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/blood-components [get]
func (h *ReferenceHandler) GetBloodComponentsHandler(c echo.Context) error {
	logger.Log.Info("получение справочника компоненотов крови")

	bloodComponents, err := h.bloodRepo.AllComponents(c.Request().Context())
	if err != nil {
		logger.Log.Error("не удалось получить компоненотов крови из БД", zap.Error(err))
		return c.JSON(500, utils.ErrorResponse{
			Message: "Не удалось получить список компоненотов крови",
		})
	}

	items := make([]ReferenceItemDB, len(bloodComponents))
	for i, bloodType := range bloodComponents {
		items[i] = ReferenceItemDB{
			Value: bloodType.ID,
			Label: bloodType.Name,
		}
	}

	return c.JSON(200, ReferenceResponseDB{Data: items})
}

// GetBloodGroupsHandler godoc
// @Summary Получение групп крови животных по типу животного
// @Description Возвращает список групп крови животных для выбора на фронтенде
// @Tags reference
// @Produce json
// @Param pet_type path string true "Тип животного"
// @Success 200 {object} ReferenceResponseDB "Список групп крови"
// @Failure 400 {object} utils.ErrorResponse "Неверный тип животного"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /reference/blood-groups/{pet_type} [get]
func (h *ReferenceHandler) GetBloodGroupsHandler(c echo.Context) error {
	petType := c.Param("pet_type")
	if petType == "" {
		logger.Log.Error("не указан тип животного")
		return c.JSON(400, utils.ErrorResponse{
			Message: "Необходимо указать тип животного",
		})
	}

	logger.Log.Info("получение групп крови по типу животного", zap.String("petType", petType))

	bloodGroups, err := h.bloodRepo.BloodGroupsByPetType(c.Request().Context(), bloodgroup.PetType(petType))
	if err != nil {
		logger.Log.Error("не удалось получить группы крови из БД", zap.Error(err), zap.String("petType", petType))
		return c.JSON(500, utils.ErrorResponse{
			Message: "Не удалось получить список групп крови",
		})
	}

	items := make([]ReferenceItemDB, len(bloodGroups))
	for i, bloodGroup := range bloodGroups {
		items[i] = ReferenceItemDB{
			Value: bloodGroup.ID,
			Label: bloodGroup.BloodGroup,
		}
	}

	return c.JSON(200, ReferenceResponseDB{Data: items})
}

// GetHealthStatusesHandler godoc
// @Summary Получение всех статусов здоровья
// @Description Возвращает все доступные статусы здоровья для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список статусов здоровья"
// @Router /reference/health-statuses [get]
func (h *ReferenceHandler) GetHealthStatusesHandler(c echo.Context) error {
	logger.Log.Info("получение справочника статусов здоровья")

	statuses := enums.GetAllHealthStatuses()
	items := make([]ReferenceItem, len(statuses))

	for i, status := range statuses {
		ruValue, err := validation.LocalizeHealthStatus(string(status))
		if err != nil {
			return c.JSON(500, utils.ErrorResponse{
				Message: "Не удалось получить список статусов здоровья",
			})
		}

		items[i] = ReferenceItem{
			Value: string(status),
			Label: string(ruValue),
		}
	}

	return c.JSON(200, ReferenceResponse{Data: items})
}

// GetReproductiveStatusesHandler godoc
// @Summary Получение всех репродуктивных состояний
// @Description Возвращает все доступные репродуктивные состояния для выбора на фронтенде
// @Tags reference
// @Produce json
// @Success 200 {object} ReferenceResponse "Список репродуктивных состояний"
// @Router /reference/reproductive-statuses [get]
func (h *ReferenceHandler) GetReproductiveStatusesHandler(c echo.Context) error {
	logger.Log.Info("получение справочника репродуктивных состояний")

	statuses := enums.GetAllReproductiveStatuses()
	items := make([]ReferenceItem, len(statuses))

	for i, status := range statuses {
		ruValue, err := validation.LocalizeReproductiveStatus(string(status))
		if err != nil {
			return c.JSON(500, utils.ErrorResponse{
				Message: "Не удалось получить список репродуктивных состояний",
			})
		}

		items[i] = ReferenceItem{
			Value: string(status),
			Label: string(ruValue),
		}
	}

	return c.JSON(200, ReferenceResponse{Data: items})
}
