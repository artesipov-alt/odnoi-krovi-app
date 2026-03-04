package http

import (
	"context"
	"log/slog"
	"net/http"

	bloodcmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/cmd"
	bloodquery "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/mapper"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"

	"github.com/danielgtaylor/huma/v2"
)

// calculateDonationAmount вычисляет максимальный объем донации крови для питомца (до 20% циркулирующей крови, но не более лимита)
// Для собак: не более 17.6 мл/кг
// Для кошек: не более 13.2 мл/кг
func calculateDonationAmount(petType string, weightKg float64) int32 {
	var limitPerKg float64
	switch petType {
	case "dog":
		limitPerKg = 17.6
	case "cat":
		limitPerKg = 13.2
	default:
		return 0
	}
	amount := limitPerKg * weightKg
	return int32(amount)
}

// BloodRequestHandler обрабатывает HTTP запросы для операций с заявками на поиск крови
type BloodRequestHandler struct {
	createHandler       *bloodcmd.CreateRequestHandler
	updateHandler       *bloodcmd.UpdateRequestHandler
	updateStatusHandler *bloodcmd.UpdateStatusHandler
	deleteHandler       *bloodcmd.DeleteRequestHandler
	applyHandler        *bloodcmd.ApplyForRequestHandler
	getByIDHandler      *bloodquery.GetByIDHandler
	getByPetIDHandler   *bloodquery.GetByPetIDHandler
	listHandler         *bloodquery.ListRequestsHandler
	bloodRequestMapper  *mapper.BloodRequestMapper
	storage             filestorage.Repository
}

// NewBloodRequestHandler создает новый обработчик для заявок на поиск крови
func NewBloodRequestHandler(
	createHandler *bloodcmd.CreateRequestHandler,
	updateHandler *bloodcmd.UpdateRequestHandler,
	updateStatusHandler *bloodcmd.UpdateStatusHandler,
	deleteHandler *bloodcmd.DeleteRequestHandler,
	applyHandler *bloodcmd.ApplyForRequestHandler,
	getByIDHandler *bloodquery.GetByIDHandler,
	getByPetIDHandler *bloodquery.GetByPetIDHandler,
	listHandler *bloodquery.ListRequestsHandler,
	storage filestorage.Repository,
) *BloodRequestHandler {
	return &BloodRequestHandler{
		createHandler:       createHandler,
		updateHandler:       updateHandler,
		updateStatusHandler: updateStatusHandler,
		deleteHandler:       deleteHandler,
		applyHandler:        applyHandler,
		getByIDHandler:      getByIDHandler,
		getByPetIDHandler:   getByPetIDHandler,
		listHandler:         listHandler,
		bloodRequestMapper:  mapper.NewBloodRequestMapper(storage),
		storage:             storage,
	}
}

// Register регистрирует маршруты заявок на поиск крови в Huma API
func (h *BloodRequestHandler) Register(api huma.API) {
	// Добавить питомца в пул поиска крови
	huma.Register(api, huma.Operation{
		OperationID:   "add-pet-to-blood-request-pool",
		Method:        http.MethodPost,
		Path:          "/v1/blood-request/pool",
		Summary:       "Добавить питомца в пул поиска крови",
		Description:   "Создает новую заявку на поиск крови для питомца",
		Tags:          []string{"blood-request-v1"},
		DefaultStatus: http.StatusCreated,
	}, h.AddPetToBloodRequestPool)

	// // Получить список заявок на поиск крови
	// huma.Register(api, huma.Operation{
	// 	OperationID: "get-pets-from-blood-request-pool",
	// 	Method:      http.MethodPost,
	// 	Path:        "/v1/blood-request/pool/search",
	// 	Summary:     "Получить список заявок на поиск крови",
	// 	Description: "Возвращает список заявок по фильтрам",
	// 	Tags:        []string{"blood-request-v1"},
	// }, h.GetPetsFromBloodRequestPool)

	// Получить заявку по ID
	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Получить заявку по ID",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByID)

	// Получить заявку по ID питомца
	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-pet-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/pet/{id}",
		Summary:     "Получить заявку по ID питомца",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByPetID)

	// Откликнуться на заявку
	huma.Register(api, huma.Operation{
		OperationID:   "apply-for-blood-request", // More descriptive OperationID
		Method:        http.MethodPost,
		Path:          "/v1/blood-request/apply/{id}", // RESTful path for applying to a specific request
		Summary:       "Откликнуться на заявку на поиск крови",
		Description:   "Позволяет донору откликнуться на существующую заявку на поиск крови.",
		Tags:          []string{"blood-request-v1"},
		DefaultStatus: http.StatusCreated, // Applying usually creates a new application record
	}, h.ApplyForBloodRequest)

	// Обновить заявку
	huma.Register(api, huma.Operation{
		OperationID: "update-blood-request",
		Method:      http.MethodPatch,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Обновить заявку на поиск крови",
		Description: "Частично обновляет информацию о существующей заявке на поиск крови.",
		Tags:        []string{"blood-request-v1"},
	}, h.UpdateBloodRequest)

	// // Получить список доноров по ID заявки
	// huma.Register(api, huma.Operation{
	// 	OperationID: "get-donors-by-req-id",
	// 	Method:      http.MethodGet,
	// 	Path:        "/v1/blood-request/donors/{id}",
	// 	Summary:     "Получить список доноров по ID заявки",
	// 	Description: "Возвращает список доноров откликнувшихся на заявку",
	// 	Tags:        []string{"blood-request-v1"},
	// }, h.GetDonorsByID)

	// Удалить заявку
	huma.Register(api, huma.Operation{
		OperationID: "delete-blood-request",
		Method:      http.MethodDelete,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Удалить заявку",
		Description: "Удаляет заявку на поиск крови (soft delete)",
		Tags:        []string{"blood-request-v1"},
	}, h.DeleteBloodRequest)
}

// Handlers

func (h *BloodRequestHandler) AddPetToBloodRequestPool(ctx context.Context, input *dto.CreateBloodRequestInput) (*dto.CreateBloodRequestOutput, error) {
	slog.DebugContext(ctx, "adding pet to blood request pool", "pet_id", input.Body.PetID)

	bloodReq := h.bloodRequestMapper.FromCreate(input.Body)

	result, err := h.createHandler.Handle(ctx, bloodReq)
	if err != nil {
		return nil, err
	}

	return &dto.CreateBloodRequestOutput{Body: dto.CreateBloodRequestResult{
		ID:        result.ID,
		PetID:     result.PetID,
		Status:    dto.BloodRequestStatus(result.Status),
		CreatedAt: &result.CreatedAt,
	}}, nil
}

func (h *BloodRequestHandler) ApplyForBloodRequest(ctx context.Context, input *dto.ApplyForBloodRequestInput) (*dto.ApplyForBloodRequestOutput, error) {
	slog.DebugContext(ctx, "applying for blood request", "request_id", input.ID, "donor_id", input.Body.DonorID)

	resp, err := h.applyHandler.Handle(ctx, input.ID, input.Body.DonorID, input.Body.Conditions)
	if err != nil {
		return nil, err
	}

	return &dto.ApplyForBloodRequestOutput{
		Body: dto.DonorResponseResult{
			ID:        resp.ID,
			RequestID: resp.RequestID,
			DonorID:   resp.DonorID,
			Status:    dto.DonorResponseStatus(resp.Status),
			CreatedAt: &resp.CreatedAt,
		},
	}, nil
}

func (h *BloodRequestHandler) UpdateBloodRequest(ctx context.Context, input *dto.UpdateBloodRequestInput) (*dto.UpdateBloodRequestOutput, error) {
	slog.DebugContext(ctx, "updating blood request", "request_id", input.ID)

	// Получить текущий объект
	existing, err := h.getByIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	// Частично обновить поля
	if input.Body.BloodVolumeNeeded != nil {
		existing.BloodVolumeNeeded = *input.Body.BloodVolumeNeeded
	}
	if input.Body.BloodVolumeReserved != nil {
		existing.BloodVolumeReserved = *input.Body.BloodVolumeReserved
	}
	if len(input.Body.Regions) > 0 {
		existing.Regions = input.Body.Regions
	}
	if input.Body.SmallPetsNotifyAllowed != nil {
		existing.SmallPetsNotifyAllowed = *input.Body.SmallPetsNotifyAllowed
	}
	if input.Body.Description != nil {
		existing.Description = *input.Body.Description
	}
	if len(input.Body.BloodGroupNames) > 0 {
		existing.BloodGroupNames = input.Body.BloodGroupNames
	}
	if len(input.Body.BloodComponentIDs) > 0 {
		existing.BloodComponentIDs = input.Body.BloodComponentIDs
	}
	if len(input.Body.OnBoarding) > 0 {
		existing.OnBoarding = input.Body.OnBoarding
	}
	if input.Body.Status != nil {
		existing.Status = model.BloodRequestStatus(*input.Body.Status)
	}
	if input.Body.PrioritySearch != nil {
		existing.PrioritySearch = *input.Body.PrioritySearch
	}
	if input.Body.IncludeUnknownBloodGroup != nil {
		existing.IncludeUnknownBloodGroup = *input.Body.IncludeUnknownBloodGroup
	}

	result, err := h.updateHandler.Handle(ctx, input.ID, existing)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateBloodRequestOutput{Body: dto.UpdateBloodRequestResult{
		ID:        result.ID,
		UpdatedAt: &result.UpdatedAt,
	}}, nil
}

// func (h *BloodRequestHandler) GetPetsFromBloodRequestPool(ctx context.Context, input *struct {
// 	Body dto.BloodSearchFilterRequest
// }) (*dto.BloodRequestsResponse, error) {
// 	slog.DebugContext(ctx, "getting pets from blood request pool", "filters", input.Body)
// 	filters := make(map[string]any)
// 	if input.Body.PetID != "" {
// 		filters["pet_id"] = input.Body.PetID
// 	}
// 	if input.Body.Status != "" {
// 		filters["status"] = input.Body.Status
// 	}

// 	requests, err := h.svc.ListRequests(ctx, input.Body.Limit, input.Body.Offset, filters)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dtos := make([]dto.BloodSearchPetRequest, len(requests))
// 	for i, req := range requests {
// 		dtos[i] = mapBloodRequestToDTO(req)
// 	}

// 	return &dto.BloodRequestsResponse{Body: dtos}, nil
// }

func (h *BloodRequestHandler) GetBloodRequestByID(ctx context.Context, input *dto.GetBloodRequestByIDInput) (*dto.GetBloodRequestByIDOutput, error) {
	bloodReq, err := h.getByIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	zero := 0
	return &dto.GetBloodRequestByIDOutput{Body: h.bloodRequestMapper.ToResponse(bloodReq, &zero)}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByPetID(ctx context.Context, input *dto.GetBloodRequestByPetIDInput) (*dto.GetBloodRequestByPetIDOutput, error) {
	slog.DebugContext(ctx, "getting blood request by PET ID", "pet_id", input.ID)
	bloodReq, situatableDonors, err := h.getByPetIDHandler.Handle(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.GetBloodRequestByPetIDOutput{Body: h.bloodRequestMapper.ToResponse(bloodReq, &situatableDonors)}, nil
}

// func (h *BloodRequestHandler) GetDonorsByID(ctx context.Context, input *dto.IDPathStr) (*dto.PetsResponse, error) {
// 	slog.DebugContext(ctx, "getting donors by ID", "request_id", input.ID)

// 	// Create a slice of dto.Pet
// 	pets := []dto.Pet{mocks.Pet1, mocks.Pet2}

// 	return &dto.PetsResponse{Body: pets}, nil
// }

func (h *BloodRequestHandler) DeleteBloodRequest(ctx context.Context, input *dto.DeleteBloodRequestInput) (*dto.DeleteBloodRequestOutput, error) {
	slog.DebugContext(ctx, "deleting blood request", "request_id", input.ID)
	if err := h.deleteHandler.Handle(ctx, input.ID); err != nil {
		return nil, err
	}

	return &dto.DeleteBloodRequestOutput{Body: dto.DeleteBloodRequestResult{
		Message: "Заявка удалена",
	}}, nil
}
