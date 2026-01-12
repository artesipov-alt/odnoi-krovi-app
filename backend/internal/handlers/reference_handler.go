package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/enums"
	"github.com/danielgtaylor/huma/v2"
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

// Register регистрирует маршруты справочников в Huma API
func (h *ReferenceHandler) Register(api huma.API) {
	// Получение всех типов животных
	huma.Register(api, huma.Operation{
		OperationID: "get-pet-types",
		Method:      http.MethodGet,
		Path:        "/v1/reference/pet-types",
		Summary:     "Получение всех типов животных",
		Description: "Возвращает все доступные типы животных для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetPetTypes)

	// Получение всех значений пола
	huma.Register(api, huma.Operation{
		OperationID: "get-genders",
		Method:      http.MethodGet,
		Path:        "/v1/reference/genders",
		Summary:     "Получение всех значений пола",
		Description: "Возвращает все доступные значения пола для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetGenders)

	// Получение всех условий проживания
	huma.Register(api, huma.Operation{
		OperationID: "get-living-conditions",
		Method:      http.MethodGet,
		Path:        "/v1/reference/living-conditions",
		Summary:     "Получение всех условий проживания",
		Description: "Возвращает все доступные условия проживания для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetLivingConditions)

	// Получение всех ролей пользователей
	huma.Register(api, huma.Operation{
		OperationID: "get-user-roles",
		Method:      http.MethodGet,
		Path:        "/v1/reference/user-roles",
		Summary:     "Получение всех ролей пользователей",
		Description: "Возвращает все доступные роли пользователей для выбора на фронтенде",
		Tags:        []string{"reference-v1", "users-v1"},
	}, h.GetUserRoles)

	// Получение всех ролей питомцев
	huma.Register(api, huma.Operation{
		OperationID: "get-pet-roles",
		Method:      http.MethodGet,
		Path:        "/v1/reference/pet-roles",
		Summary:     "Получение всех ролей питомцев",
		Description: "Возвращает все доступные роли питомцев для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetPetRoles)

	// Получение всех пород животных
	huma.Register(api, huma.Operation{
		OperationID: "get-breeds",
		Method:      http.MethodGet,
		Path:        "/v1/reference/breeds",
		Summary:     "Получение всех пород животных",
		Description: "Возвращает список всех пород животных в базе для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetBreeds)

	// Получение всех локаций
	huma.Register(api, huma.Operation{
		OperationID: "get-locations",
		Method:      http.MethodGet,
		Path:        "/v1/reference/locations",
		Summary:     "Получение всех локаций",
		Description: "Возвращает список всех локаций в системе для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetLocations)

	// Получение пород животных по типу
	huma.Register(api, huma.Operation{
		OperationID: "get-breeds-by-type",
		Method:      http.MethodGet,
		Path:        "/v1/reference/breeds-by-type",
		Summary:     "Получение пород животных по типу",
		Description: "Возвращает список пород животных для указанного типа животного для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetBreedsByType)

	// Получение компонентов крови животных
	huma.Register(api, huma.Operation{
		OperationID: "get-blood-components",
		Method:      http.MethodGet,
		Path:        "/v1/reference/blood-components",
		Summary:     "Получение компонентов крови животных",
		Description: "Возвращает список компонентов крови животных для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetBloodComponents)

	// Получение групп крови животных по типу животного
	huma.Register(api, huma.Operation{
		OperationID: "get-blood-groups",
		Method:      http.MethodGet,
		Path:        "/v1/reference/blood-groups/{pet_type}",
		Summary:     "Получение групп крови животных по типу животного",
		Description: "Возвращает список групп крови животных для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetBloodGroups)

	// Получение всех статусов здоровья
	huma.Register(api, huma.Operation{
		OperationID: "get-health-statuses",
		Method:      http.MethodGet,
		Path:        "/v1/reference/health-statuses",
		Summary:     "Получение всех статусов здоровья",
		Description: "Возвращает все доступные статусы здоровья для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetHealthStatuses)

	// Получение всех репродуктивных состояний
	huma.Register(api, huma.Operation{
		OperationID: "get-reproductive-statuses",
		Method:      http.MethodGet,
		Path:        "/v1/reference/reproductive-statuses",
		Summary:     "Получение всех репродуктивных состояний",
		Description: "Возвращает все доступные репродуктивные состояния для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetReproductiveStatuses)
}

// Handlers

func (h *ReferenceHandler) GetPetTypes(ctx context.Context, input *struct{}) (*dto.ReferenceCodeResponse, error) {
	petTypes := enums.GetAllEntPetTypes()
	items := make([]dto.ReferenceCodeItem, len(petTypes))

	for i, petType := range petTypes {
		ruValue := enums.LocalizeEntPetType(petType)

		items[i] = dto.ReferenceCodeItem{
			Value: string(petType),
			Label: ruValue,
		}
	}

	return &dto.ReferenceCodeResponse{Body: dto.ReferenceCodeData{Data: items}}, nil
}

func (h *ReferenceHandler) GetGenders(ctx context.Context, input *struct{}) (*dto.ReferenceCodeResponse, error) {
	genders := enums.GetAllEntGenders()
	items := make([]dto.ReferenceCodeItem, len(genders))

	for i, gender := range genders {
		ruValue := enums.LocalizeEntGender(gender)

		items[i] = dto.ReferenceCodeItem{
			Value: string(gender),
			Label: ruValue,
		}
	}

	return &dto.ReferenceCodeResponse{Body: dto.ReferenceCodeData{Data: items}}, nil
}

func (h *ReferenceHandler) GetLivingConditions(ctx context.Context, input *struct{}) (*dto.ReferenceCodeResponse, error) {
	conditions := enums.GetAllEntLivingConditions()
	items := make([]dto.ReferenceCodeItem, len(conditions))

	for i, condition := range conditions {
		ruValue := enums.LocalizeEntLivingCondition(condition)

		items[i] = dto.ReferenceCodeItem{
			Value: string(condition),
			Label: ruValue,
		}
	}

	return &dto.ReferenceCodeResponse{Body: dto.ReferenceCodeData{Data: items}}, nil
}

func (h *ReferenceHandler) GetUserRoles(ctx context.Context, input *struct{}) (*dto.ReferenceCodeResponse, error) {
	roles := enums.GetAllEntUserRoles()
	items := make([]dto.ReferenceCodeItem, len(roles))

	for i, role := range roles {
		ruValue := enums.LocalizeEntUserRole(role)

		items[i] = dto.ReferenceCodeItem{
			Value: string(role),
			Label: ruValue,
		}
	}

	return &dto.ReferenceCodeResponse{Body: dto.ReferenceCodeData{Data: items}}, nil
}

func (h *ReferenceHandler) GetPetRoles(ctx context.Context, input *struct{}) (*dto.ReferenceCodeResponse, error) {
	roles := enums.GetAllEntPetStatuses()
	items := make([]dto.ReferenceCodeItem, len(roles))

	for i, role := range roles {
		ruValue := enums.LocalizeEntPetStatus(role)

		items[i] = dto.ReferenceCodeItem{
			Value: string(role),
			Label: ruValue,
		}
	}

	return &dto.ReferenceCodeResponse{Body: dto.ReferenceCodeData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBreeds(ctx context.Context, input *struct{}) (*dto.ReferenceDBResponse, error) {
	breeds, err := h.breedRepo.GetAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get all breeds", "error", err.Error())
		return nil, apperrors.Internal(err, "Ошибка сервера при получении пород")
	}

	items := make([]dto.ReferenceDBItem, len(breeds))
	for i, breed := range breeds {
		items[i] = dto.ReferenceDBItem{
			Value: breed.ID,
			Label: breed.Name,
		}
	}

	return &dto.ReferenceDBResponse{Body: dto.ReferenceDBData{Data: items}}, nil
}

func (h *ReferenceHandler) GetLocations(ctx context.Context, input *struct{}) (*dto.ReferenceDBResponse, error) {
	locations, err := h.locationRepo.GetAll(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get all locations", "error", err.Error())
		return nil, apperrors.Internal(err, "Ошибка сервера при получении локаций")
	}

	items := make([]dto.ReferenceDBItem, len(locations))
	for i, location := range locations {
		items[i] = dto.ReferenceDBItem{
			Value: location.ID,
			Label: location.Name,
		}
	}

	return &dto.ReferenceDBResponse{Body: dto.ReferenceDBData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBreedsByType(ctx context.Context, input *dto.PetTypeQuery) (*dto.ReferenceDBResponse, error) {
	petTypeStr := input.PetType
	if petTypeStr == "" {
		slog.DebugContext(ctx, "pet type not specified for GetBreedsByType")
		return nil, apperrors.BadRequest("Необходимо указать тип животного")
	}

	isValid := false
	for _, pt := range enums.GetAllEntPetTypes() {
		if string(pt) == petTypeStr {
			isValid = true
			break
		}
	}
	if !isValid {
		slog.DebugContext(ctx, "invalid pet type for GetBreedsByType", "pet_type", petTypeStr)
		return nil, apperrors.BadRequest("Неверный тип животного")
	}

	breeds, err := h.breedRepo.GetByPetType(ctx, breed.Type(petTypeStr))
	if err != nil {
		slog.ErrorContext(ctx, "failed to get breeds by pet type", "pet_type", petTypeStr, "error", err.Error())
		return nil, apperrors.Internal(err, "Ошибка сервера при получении пород по типу животного")
	}

	items := make([]dto.ReferenceDBItem, len(breeds))
	for i, breed := range breeds {
		items[i] = dto.ReferenceDBItem{
			Value: breed.ID,
			Label: breed.Name,
		}
	}

	return &dto.ReferenceDBResponse{Body: dto.ReferenceDBData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBloodComponents(ctx context.Context, input *struct{}) (*dto.ReferenceDBResponse, error) {
	bloodComponents, err := h.bloodRepo.AllComponents(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get all blood components", "error", err.Error())
		return nil, apperrors.Internal(err, "Ошибка сервера при получении компонентов крови")
	}

	items := make([]dto.ReferenceDBItem, len(bloodComponents))
	for i, bloodType := range bloodComponents {
		items[i] = dto.ReferenceDBItem{
			Value: bloodType.ID,
			Label: bloodType.Name,
		}
	}

	return &dto.ReferenceDBResponse{Body: dto.ReferenceDBData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBloodGroups(ctx context.Context, input *dto.PetTypePath) (*dto.ReferenceDBResponse, error) {
	petType := input.PetType
	if petType == "" {
		slog.DebugContext(ctx, "pet type not specified for GetBloodGroups")
		return nil, apperrors.BadRequest("Необходимо указать тип животного")
	}

	bloodGroups, err := h.bloodRepo.BloodGroupsByPetType(ctx, bloodgroup.PetType(petType))
	if err != nil {
		slog.ErrorContext(ctx, "failed to get blood groups by pet type", "pet_type", petType, "error", err.Error())
		return nil, apperrors.Internal(err, "Ошибка сервера при получении групп крови по типу животного")
	}

	items := make([]dto.ReferenceDBItem, len(bloodGroups))
	for i, bloodGroup := range bloodGroups {
		items[i] = dto.ReferenceDBItem{
			Value: bloodGroup.ID,
			Label: bloodGroup.BloodGroup,
		}
	}

	return &dto.ReferenceDBResponse{Body: dto.ReferenceDBData{Data: items}}, nil
}

func (h *ReferenceHandler) GetHealthStatuses(ctx context.Context, input *struct{}) (*dto.ReferenceCodeResponse, error) {
	statuses := enums.GetAllEntHealthStatuses()
	items := make([]dto.ReferenceCodeItem, len(statuses))

	for i, status := range statuses {
		ruValue := enums.LocalizeEntHealthStatus(status)

		items[i] = dto.ReferenceCodeItem{
			Value: string(status),
			Label: ruValue,
		}
	}

	return &dto.ReferenceCodeResponse{Body: dto.ReferenceCodeData{Data: items}}, nil
}

func (h *ReferenceHandler) GetReproductiveStatuses(ctx context.Context, input *struct{}) (*dto.ReferenceCodeResponse, error) {
	statuses := enums.GetAllEntReproductiveStatuses()
	items := make([]dto.ReferenceCodeItem, len(statuses))

	for i, status := range statuses {
		ruValue := enums.LocalizeEntReproductiveStatus(status)

		items[i] = dto.ReferenceCodeItem{
			Value: string(status),
			Label: ruValue,
		}
	}

	return &dto.ReferenceCodeResponse{Body: dto.ReferenceCodeData{Data: items}}, nil
}
