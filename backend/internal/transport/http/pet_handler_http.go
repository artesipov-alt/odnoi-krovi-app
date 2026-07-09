package http

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	petcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/cmd"
	petquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	mapper "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dtomapper"
	"github.com/danielgtaylor/huma/v2"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	createHandler    *petcmd.CreateHandler
	updateHandler    *petcmd.UpdateHandler
	deleteHandler    *petcmd.DeleteHandler
	getByIDHandler   *petquery.GetByIDHandler
	getByUserHandler *petquery.GetByUserHandler
	petMapper        *mapper.PetMapper
	storage          filestorage.Repository
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(
	createHandler *petcmd.CreateHandler,
	updateHandler *petcmd.UpdateHandler,
	deleteHandler *petcmd.DeleteHandler,
	getByIDHandler *petquery.GetByIDHandler,
	getByUserHandler *petquery.GetByUserHandler,
	storage filestorage.Repository,
) *PetHandler {
	return &PetHandler{
		createHandler:    createHandler,
		updateHandler:    updateHandler,
		deleteHandler:    deleteHandler,
		getByIDHandler:   getByIDHandler,
		getByUserHandler: getByUserHandler,
		petMapper:        mapper.NewPetMapper(storage),
		storage:          storage,
	}
}

// Register регистрирует маршруты питомцев в Huma API
func (h *PetHandler) Register(api huma.API) {
	// Создание нового питомца
	huma.Register(api, huma.Operation{
		OperationID:   "create-pet",
		Method:        http.MethodPost,
		Path:          "/v1/pet/user/{user_id}",
		Summary:       "Создание нового питомца",
		Description:   "Создает нового питомца для пользователя",
		Tags:          []string{"pets-v1"},
		DefaultStatus: http.StatusCreated,
	}, h.CreatePet)

	// Получение питомца по ID
	huma.Register(api, huma.Operation{
		OperationID: "get-pet-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/pet/{pet_id}",
		Summary:     "Получение питомца по ID",
		Description: "Возвращает информацию о питомце по его идентификатору",
		Tags:        []string{"pets-v1"},
	}, h.GetPet)

	// Получение питомцев пользователя
	huma.Register(api, huma.Operation{
		OperationID: "get-user-pets",
		Method:      http.MethodGet,
		Path:        "/v1/pet/user/{user_id}",
		Summary:     "Получение питомцев пользователя",
		Description: "Возвращает всех питомцев конкретного пользователя",
		Tags:        []string{"pets-v1"},
	}, h.GetUserPets)

	// Обновление данных питомца
	huma.Register(api, huma.Operation{
		OperationID: "update-pet",
		Method:      http.MethodPut,
		Path:        "/v1/pet/{pet_id}",
		Summary:     "Обновление данных питомца",
		Description: "Обновляет информацию о питомце",
		Tags:        []string{"pets-v1"},
	}, h.UpdatePet)

	// Удаление питомца по ID
	huma.Register(api, huma.Operation{
		OperationID: "delete-pet",
		Method:      http.MethodDelete,
		Path:        "/v1/pet/{pet_id}",
		Summary:     "Удаление питомца по ID",
		Description: "Удаляет питомца из системы",
		Tags:        []string{"pets-v1"},
	}, h.DeletePet)
}

// CreatePet создает нового питомца
func (h *PetHandler) CreatePet(ctx context.Context, input *dto.CreatePetInput) (*dto.CreatePetOutput, error) {
	petDomain, err := h.petMapper.FromCreate(input.Body)
	if err != nil {
		return nil, apperrors.Validation("некорректные данные питомца", map[string]any{"error": err.Error()})
	}

	userID := input.UserIDPath.ID

	createdPet, err := h.createHandler.Handle(ctx, userID, petDomain)
	if err != nil {
		return nil, err
	}

	return &dto.CreatePetOutput{
		Body: dto.CreatePetResult{
			ID:        createdPet.ID,
			CreatedAt: createdPet.CreatedAt,
		},
	}, nil
}

// UpdatePet обновляет данные питомца
func (h *PetHandler) UpdatePet(ctx context.Context, input *dto.UpdatePetInput) (*dto.UpdatePetOutput, error) {
	petDomain := h.petMapper.ToUpdateModel(input.Body)

	updatedPet, err := h.updateHandler.Handle(ctx, input.ID, petDomain)
	if err != nil {
		return nil, err
	}

	return &dto.UpdatePetOutput{
		Body: dto.UpdatePetResult{
			ID:        updatedPet.ID,
			UpdatedAt: updatedPet.UpdatedAt,
		},
	}, nil
}

// GetPet возвращает питомца по ID
func (h *PetHandler) GetPet(ctx context.Context, input *dto.GetPetByIDInput) (*dto.GetPetByIDOutput, error) {
	opts := pet.PetPreloadOptions{
		WithAll: input.WithAll,
	}

	petResult, err := h.getByIDHandler.Handle(ctx, input.ID, opts)
	if err != nil {
		return nil, err
	}

	return &dto.GetPetByIDOutput{
		Body: h.petMapper.ToResponse(*petResult),
	}, nil
}

// GetUserPets возвращает всех питомцев пользователя
func (h *PetHandler) GetUserPets(ctx context.Context, input *dto.GetPetsByUserInput) (*dto.GetPetsByUserOutput, error) {
	opts := pet.PetPreloadOptions{
		WithHealth:     input.WithHealth,
		WithTreatments: input.WithTreatments,
		WithAnalyses:   input.WithAnalysis,
		WithBonuses:    input.WithBonuses,
		WithAll:        input.WithAll,
	}

	result, err := h.getByUserHandler.Handle(ctx, input.UserIDPath.ID, opts)
	if err != nil {
		return nil, err
	}

	return &dto.GetPetsByUserOutput{
		Body: dto.GetPetsByUserResult{
			Pets:                    h.petMapper.ToResponseSlice(result.Pets),
			TotalPets:               result.TotalPets,
			TotalPlannedDonations:   result.TotalPlannedDonations,
			TotalCompletedDonations: result.TotalCompletedDonations,
			TotalPrioritySearch:     result.TotalPrioritySearch,
			TotalBonuses:            result.TotalBonuses,
		},
	}, nil
}

// DeletePet удаляет питомца
func (h *PetHandler) DeletePet(ctx context.Context, input *dto.DeletePetInput) (*dto.DeletePetOutput, error) {
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.DeletePetOutput{
		Body: dto.DeletePetResult{
			Message: "Питомец успешно удален",
		},
	}, nil
}
