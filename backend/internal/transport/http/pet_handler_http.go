package http

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	petcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/cmd"
	petquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/reference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/mapper"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/danielgtaylor/huma/v2"
)

// PetHandler обрабатывает HTTP запросы для операций с питомцами
type PetHandler struct {
	createHandler     *petcmd.CreateHandler
	updateHandler     *petcmd.UpdateHandler
	deleteHandler     *petcmd.DeleteHandler
	revalidateHandler *petcmd.RevalidateDonorHandler
	getByIDHandler    *petquery.GetByIDHandler
	getByUserHandler  *petquery.GetByUserHandler
	bloodInfoRepo     reference.BloodInfoRepository
	petMapper         *mapper.PetMapper
}

// NewPetHandler создает новый обработчик питомцев
func NewPetHandler(
	createHandler *petcmd.CreateHandler,
	updateHandler *petcmd.UpdateHandler,
	deleteHandler *petcmd.DeleteHandler,
	revalidateHandler *petcmd.RevalidateDonorHandler,
	getByIDHandler *petquery.GetByIDHandler,
	getByUserHandler *petquery.GetByUserHandler,
	bloodInfoRepo reference.BloodInfoRepository,
) *PetHandler {
	return &PetHandler{
		createHandler:     createHandler,
		updateHandler:     updateHandler,
		deleteHandler:     deleteHandler,
		revalidateHandler: revalidateHandler,
		getByIDHandler:    getByIDHandler,
		getByUserHandler:  getByUserHandler,
		bloodInfoRepo:     bloodInfoRepo,
		petMapper:         mapper.NewPetMapper(),
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
		Path:        "/v1/pet/{id}",
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
		Path:        "/v1/pet/{id}",
		Summary:     "Обновление данных питомца",
		Description: "Обновляет информацию о питомце",
		Tags:        []string{"pets-v1"},
	}, h.UpdatePet)

	// Удаление питомца по ID
	huma.Register(api, huma.Operation{
		OperationID: "delete-pet",
		Method:      http.MethodDelete,
		Path:        "/v1/pet/{id}",
		Summary:     "Удаление питомца по ID",
		Description: "Удаляет питомца из системы",
		Tags:        []string{"pets-v1"},
	}, h.DeletePet)

	// Валидация донора по ID (изменено на POST)
	huma.Register(api, huma.Operation{
		OperationID:   "validate-donor",
		Method:        http.MethodPost,
		Path:          "/v1/pet/validate-donor/{id}",
		Summary:       "Валидация донора по ID",
		Description:   "Пересчитывает и сохраняет факторы валидации донора для питомца",
		Tags:          []string{"pets-v1"},
		DefaultStatus: http.StatusOK,
	}, h.ValidateDonor)

}

//==========Handlers==============================

func (h *PetHandler) CreatePet(ctx context.Context, input *struct {
	dto.PetUserIDPath
	Body dto.PetCreate
}) (*dto.BodyPetCreateResponse, error) {
	body := &input.Body

	petDomain, err := h.petMapper.FromCreate(*body)
	if err != nil {
		return nil, apperrors.Validation("invalid pet data", map[string]interface{}{"error": err.Error()})
	}

	createdPet, err := h.createHandler.Handle(ctx, input.ID, petDomain)
	if err != nil {
		return nil, err
	}

	return &dto.BodyPetCreateResponse{Body: dto.PetCreateResponse{ID: createdPet.ID, CreatedAt: createdPet.CreatedAt}}, nil
}

func (h *PetHandler) UpdatePet(ctx context.Context, input *struct {
	dto.IDPathStr
	Body dto.PetUpdate
}) (*dto.PetResponse, error) {
	body := &input.Body

	// Create empty domain model and apply updates to it
	petDomain := &model.Pet{}
	h.petMapper.ApplyUpdate(*body, petDomain)

	updatedPet, err := h.updateHandler.Handle(ctx, input.ID, petDomain)
	if err != nil {
		return nil, err
	}

	return &dto.PetResponse{Body: h.petMapper.ToResponse(*updatedPet)}, nil
}

func (h *PetHandler) GetPet(ctx context.Context,
	input *struct {
		dto.IDPathStr
		dto.PetPreloadQuery
	}) (*dto.PetResponse, error) {
	// Добавляем опции к запросу.
	opts := pet.PetPreloadOptions{
		WithHealth:     input.WithHealth,
		WithTreatments: input.WithTreatments,
		WithAnalyses:   input.WithAnalysis,
		WithBonuses:    input.WithBonuses,
		WithAll:        input.WithAll,
	}

	pet, err := h.getByIDHandler.Handle(ctx, input.ID, opts)
	if err != nil {
		return nil, err
	}

	return &dto.PetResponse{Body: h.petMapper.ToResponse(*pet)}, nil
}

func (h *PetHandler) GetUserPets(ctx context.Context,
	input *struct {
		dto.PetUserIDPath
		dto.PetPreloadQuery
	}) (*dto.PetsResponse, error) {

	opts := pet.PetPreloadOptions{
		WithHealth:     input.WithHealth,
		WithTreatments: input.WithTreatments,
		WithAnalyses:   input.WithAnalysis,
		WithBonuses:    input.WithBonuses,
		WithAll:        input.WithAll,
	}

	pets, err := h.getByUserHandler.Handle(ctx, input.ID, opts)
	if err != nil {
		return nil, err
	}

	return &dto.PetsResponse{Body: h.petMapper.ToResponseSlice(pets)}, nil
}

func (h *PetHandler) DeletePet(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "Питомец удален"
	return resp, nil
}

func (h *PetHandler) ValidateDonor(ctx context.Context, input *dto.IDPathStr) (*dto.PetResponse, error) {
	pet, err := h.revalidateHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.PetResponse{Body: h.petMapper.ToResponse(*pet)}, nil
}
