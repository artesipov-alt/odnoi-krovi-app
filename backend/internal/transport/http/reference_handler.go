package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/reference/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/enums"
	"github.com/danielgtaylor/huma/v2"
)

// ReferenceHandler обрабатывает HTTP запросы для справочных данных
type ReferenceHandler struct {
	getAllBreedsHandler          *query.GetAllBreedsHandler
	getBreedsByTypeHandler       *query.GetBreedsByPetTypeHandler
	getAllLocationsHandler       *query.GetAllLocationsHandler
	getAllBloodComponentsHandler *query.GetAllBloodComponentsHandler
	getBloodGroupsByTypeHandler  *query.GetBloodGroupsByPetTypeHandler
}

// NewReferenceHandler создает новый обработчик справочных данных
func NewReferenceHandler(
	getAllBreedsHandler *query.GetAllBreedsHandler,
	getBreedsByTypeHandler *query.GetBreedsByPetTypeHandler,
	getAllLocationsHandler *query.GetAllLocationsHandler,
	getAllBloodComponentsHandler *query.GetAllBloodComponentsHandler,
	getBloodGroupsByTypeHandler *query.GetBloodGroupsByPetTypeHandler,
) *ReferenceHandler {
	return &ReferenceHandler{
		getAllBreedsHandler:          getAllBreedsHandler,
		getBreedsByTypeHandler:       getBreedsByTypeHandler,
		getAllLocationsHandler:       getAllLocationsHandler,
		getAllBloodComponentsHandler: getAllBloodComponentsHandler,
		getBloodGroupsByTypeHandler:  getBloodGroupsByTypeHandler,
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

	huma.Register(api, huma.Operation{
		OperationID: "get-donor-restrictions",
		Method:      http.MethodGet,
		Path:        "/v1/reference/donor-restrictions",
		Summary:     "Получение всех ограничений для доноров",
		Description: "Возвращает все доступные ограничения и предупреждения для доноров для выбора на фронтенде",
		Tags:        []string{"reference-v1"},
	}, h.GetDonorRestrictions)

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
		Tags:        []string{"reference-v1"},
	}, h.GetUserRoles)

	// // Получение всех ролей питомцев
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

func (h *ReferenceHandler) GetPetTypes(ctx context.Context, input *dto.GetPetTypesInput) (*dto.GetPetTypesOutput, error) {
	petTypes := enums.GetAllEntPetTypes()
	items := make([]dto.ReferenceItem, len(petTypes))

	for i, petType := range petTypes {
		ruValue := enums.LocalizeEntPetType(petType)

		items[i] = dto.ReferenceItem{
			Value: string(petType),
			Label: ruValue,
		}
	}

	return &dto.GetPetTypesOutput{Body: dto.PetTypesList{Data: items}}, nil
}

func (h *ReferenceHandler) GetDonorRestrictions(ctx context.Context, input *dto.GetDonorRestrictionsInput) (*dto.GetDonorRestrictionsOutput, error) {
	allFactors := petmodel.GetAllFactors()
	var stopFactors []dto.FactorDescription
	var warnFactors []dto.FactorDescription

	for code, desc := range allFactors {
		factor := dto.FactorDescription{
			Code:           string(code),
			Description:    desc.Description,
			SubDescription: desc.SubDescription,
		}
		if strings.HasPrefix(string(code), "STOP_") {
			stopFactors = append(stopFactors, factor)
		} else if strings.HasPrefix(string(code), "WARN_") {
			warnFactors = append(warnFactors, factor)
		}
	}

	response := dto.DonorRestrictionsDetail{
		StopFactors: stopFactors,
		WarnFactors: warnFactors,
	}

	return &dto.GetDonorRestrictionsOutput{Body: response}, nil
}

func (h *ReferenceHandler) GetGenders(ctx context.Context, input *dto.GetGendersInput) (*dto.GetGendersOutput, error) {
	genders := enums.GetAllEntGenders()
	items := make([]dto.ReferenceItem, len(genders))

	for i, gender := range genders {
		ruValue := enums.LocalizeEntGender(gender)

		items[i] = dto.ReferenceItem{
			Value: string(gender),
			Label: ruValue,
		}
	}

	return &dto.GetGendersOutput{Body: dto.GendersList{Data: items}}, nil
}

func (h *ReferenceHandler) GetLivingConditions(ctx context.Context, input *dto.GetLivingConditionsInput) (*dto.GetLivingConditionsOutput, error) {
	conditions := enums.GetAllEntLivingConditions()
	items := make([]dto.ReferenceItem, len(conditions))

	for i, condition := range conditions {
		ruValue := enums.LocalizeEntLivingCondition(condition)

		items[i] = dto.ReferenceItem{
			Value: string(condition),
			Label: ruValue,
		}
	}

	return &dto.GetLivingConditionsOutput{Body: dto.LivingConditionsList{Data: items}}, nil
}

func (h *ReferenceHandler) GetUserRoles(ctx context.Context, input *dto.GetUserRolesInput) (*dto.GetUserRolesOutput, error) {
	roles := enums.GetAllEntUserRoles()
	items := make([]dto.ReferenceItem, len(roles))

	for i, role := range roles {
		ruValue := enums.LocalizeEntUserRole(role)

		items[i] = dto.ReferenceItem{
			Value: string(role),
			Label: ruValue,
		}
	}

	return &dto.GetUserRolesOutput{Body: dto.UserRolesList{Data: items}}, nil
}

func (h *ReferenceHandler) GetPetRoles(ctx context.Context, input *dto.GetPetRolesInput) (*dto.GetPetRolesOutput, error) {
	roles := enums.GetAllEntPetStatuses()
	items := make([]dto.ReferenceItem, len(roles))

	for i, role := range roles {
		ruValue := enums.LocalizeEntPetStatus(role)

		items[i] = dto.ReferenceItem{
			Value: string(role),
			Label: ruValue,
		}
	}

	return &dto.GetPetRolesOutput{Body: dto.PetRolesList{Data: items}}, nil
}

func (h *ReferenceHandler) GetBreeds(ctx context.Context, input *dto.GetBreedsInput) (*dto.GetBreedsOutput, error) {
	breeds, err := h.getAllBreedsHandler.Handle(ctx)
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

	return &dto.GetBreedsOutput{Body: dto.BreedsList{Data: items}}, nil
}

func (h *ReferenceHandler) GetLocations(ctx context.Context, input *dto.GetLocationsInput) (*dto.GetLocationsOutput, error) {
	locations, err := h.getAllLocationsHandler.Handle(ctx)
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

	return &dto.GetLocationsOutput{Body: dto.LocationsList{Data: items}}, nil
}

func (h *ReferenceHandler) GetBreedsByType(ctx context.Context, input *dto.GetBreedsByTypeInput) (*dto.GetBreedsByTypeOutput, error) {
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

	breeds, err := h.getBreedsByTypeHandler.Handle(ctx, breed.Type(petTypeStr))
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

	return &dto.GetBreedsByTypeOutput{Body: dto.BreedsList{Data: items}}, nil
}

func (h *ReferenceHandler) GetBloodComponents(ctx context.Context, input *dto.GetBloodComponentsInput) (*dto.GetBloodComponentsOutput, error) {
	bloodComponents, err := h.getAllBloodComponentsHandler.Handle(ctx)
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

	return &dto.GetBloodComponentsOutput{Body: dto.BloodComponentsList{Data: items}}, nil
}

func (h *ReferenceHandler) GetBloodGroups(ctx context.Context, input *dto.GetBloodGroupsInput) (*dto.GetBloodGroupsOutput, error) {
	petTypeStr := input.PetType
	if petTypeStr == "" {
		return nil, apperrors.BadRequest("Необходимо указать тип животного")
	}

	petType := common.PetType(petTypeStr)
	bloodGroups, err := h.getBloodGroupsByTypeHandler.Handle(ctx, petType)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReferenceItem, len(bloodGroups))
	for i, bloodGroup := range bloodGroups {
		items[i] = dto.ReferenceItem{
			Value: "",
			Label: bloodGroup,
		}
	}

	return &dto.GetBloodGroupsOutput{Body: dto.BloodGroupsList{Data: items}}, nil
}

func (h *ReferenceHandler) GetHealthStatuses(ctx context.Context, input *dto.GetHealthStatusesInput) (*dto.GetHealthStatusesOutput, error) {
	statuses := enums.GetAllEntHealthStatuses()
	items := make([]dto.ReferenceItem, len(statuses))

	for i, status := range statuses {
		ruValue := enums.LocalizeEntHealthStatus(status)

		items[i] = dto.ReferenceItem{
			Value: string(status),
			Label: ruValue,
		}
	}

	return &dto.GetHealthStatusesOutput{Body: dto.HealthStatusesList{Data: items}}, nil
}

func (h *ReferenceHandler) GetReproductiveStatuses(ctx context.Context, input *dto.GetReproductiveStatusesInput) (*dto.GetReproductiveStatusesOutput, error) {
	statuses := enums.GetAllEntReproductiveStatuses()
	items := make([]dto.ReferenceItem, len(statuses))

	for i, status := range statuses {
		ruValue := enums.LocalizeEntReproductiveStatus(status)

		items[i] = dto.ReferenceItem{
			Value: string(status),
			Label: ruValue,
		}
	}

	return &dto.GetReproductiveStatusesOutput{Body: dto.ReproductiveStatusesList{Data: items}}, nil
}
