package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/enums"
	"github.com/danielgtaylor/huma/v2"
)

// ReferenceHandler обрабатывает HTTP запросы для справочных данных
type ReferenceHandler struct {
	breedRepo    services.BreedRepository
	bloodRepo    services.BloodInfoRepository
	locationRepo services.LocationRepository
}

// NewReferenceHandler создает новый обработчик справочных данных
func NewReferenceHandler(breedRepo services.BreedRepository, bloodTypeRepo services.BloodInfoRepository, locationRepo services.LocationRepository) *ReferenceHandler {
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

func (h *ReferenceHandler) GetPetTypes(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	petTypes := enums.GetAllEntPetTypes()
	items := make([]dto.ReferenceItem, len(petTypes))

	for i, petType := range petTypes {
		ruValue := enums.LocalizeEntPetType(petType)

		items[i] = dto.ReferenceItem{
			Value: string(petType),
			Label: ruValue,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetGenders(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	genders := enums.GetAllEntGenders()
	items := make([]dto.ReferenceItem, len(genders))

	for i, gender := range genders {
		ruValue := enums.LocalizeEntGender(gender)

		items[i] = dto.ReferenceItem{
			Value: string(gender),
			Label: ruValue,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetLivingConditions(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	conditions := enums.GetAllEntLivingConditions()
	items := make([]dto.ReferenceItem, len(conditions))

	for i, condition := range conditions {
		ruValue := enums.LocalizeEntLivingCondition(condition)

		items[i] = dto.ReferenceItem{
			Value: string(condition),
			Label: ruValue,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetUserRoles(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	roles := enums.GetAllEntUserRoles()
	items := make([]dto.ReferenceItem, len(roles))

	for i, role := range roles {
		ruValue := enums.LocalizeEntUserRole(role)

		items[i] = dto.ReferenceItem{
			Value: string(role),
			Label: ruValue,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetPetRoles(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	roles := enums.GetAllEntPetStatuses()
	items := make([]dto.ReferenceItem, len(roles))

	for i, role := range roles {
		ruValue := enums.LocalizeEntPetStatus(role)

		items[i] = dto.ReferenceItem{
			Value: string(role),
			Label: ruValue,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBreeds(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	slog.DebugContext(ctx, "getting all breeds")
	breeds, err := h.breedRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReferenceItem, len(breeds))
	for i, breed := range breeds {
		items[i] = dto.ReferenceItem{
			Value: breed.ID,
			Label: breed.Name,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetLocations(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	slog.DebugContext(ctx, "getting all locations")
	locations, err := h.locationRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReferenceItem, len(locations))
	for i, location := range locations {
		items[i] = dto.ReferenceItem{
			Value: location.ID,
			Label: location.Name,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBreedsByType(ctx context.Context, input *dto.PetTypeQuery) (*dto.ReferenceResponse, error) {
	slog.DebugContext(ctx, "getting breeds by type", "pet_type", input.PetType)
	petTypeStr := input.PetType
	if petTypeStr == "" {
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
		return nil, apperrors.BadRequest("Неверный тип животного")
	}

	breeds, err := h.breedRepo.GetByPetType(ctx, breed.Type(petTypeStr))
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReferenceItem, len(breeds))
	for i, breed := range breeds {
		items[i] = dto.ReferenceItem{
			Value: breed.ID,
			Label: breed.Name,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBloodComponents(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	slog.DebugContext(ctx, "getting blood components")
	bloodComponents, err := h.bloodRepo.AllComponents(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReferenceItem, len(bloodComponents))
	for i, bloodType := range bloodComponents {
		items[i] = dto.ReferenceItem{
			Value: bloodType.ID,
			Label: bloodType.Name,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetBloodGroups(ctx context.Context, input *dto.PetTypePath) (*dto.ReferenceResponse, error) {
	slog.DebugContext(ctx, "getting blood groups", "pet_type", input.PetType)
	petType := input.PetType
	if petType == "" {
		return nil, apperrors.BadRequest("Необходимо указать тип животного")
	}

	bloodGroups, err := h.bloodRepo.BloodGroupsByPetType(ctx, bloodgroup.PetType(petType))
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReferenceItem, len(bloodGroups))
	for i, bloodGroup := range bloodGroups {
		items[i] = dto.ReferenceItem{
			Value: bloodGroup.ID,
			Label: bloodGroup.BloodGroup,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetHealthStatuses(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	statuses := enums.GetAllEntHealthStatuses()
	items := make([]dto.ReferenceItem, len(statuses))

	for i, status := range statuses {
		ruValue := enums.LocalizeEntHealthStatus(status)

		items[i] = dto.ReferenceItem{
			Value: string(status),
			Label: ruValue,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}

func (h *ReferenceHandler) GetReproductiveStatuses(ctx context.Context, input *struct{}) (*dto.ReferenceResponse, error) {
	statuses := enums.GetAllEntReproductiveStatuses()
	items := make([]dto.ReferenceItem, len(statuses))

	for i, status := range statuses {
		ruValue := enums.LocalizeEntReproductiveStatus(status)

		items[i] = dto.ReferenceItem{
			Value: string(status),
			Label: ruValue,
		}
	}

	return &dto.ReferenceResponse{Body: dto.ReferenceData{Data: items}}, nil
}
