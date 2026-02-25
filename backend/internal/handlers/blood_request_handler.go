package handlers

import (
	"context"
	"log/slog" // Import slog
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jinzhu/copier"
)

// calculateDonationAmount вычисляет максимальный объем донации крови для питомца (до 20% циркулирующей крови, но не более лимита)
// Для собак: не более 17.6 мл/кг
// Для кошек: не более 13.2 мл/кг
func calculateDonationAmount(pet *ent.Pet) int32 {
	weight := pet.WeightKg
	var limitPerKg float64
	switch pet.Type {
	case "dog":
		limitPerKg = 17.6
	case "cat":
		limitPerKg = 13.2
	default:
		return 0
	}
	amount := limitPerKg * weight
	return int32(amount)
}

// BloodRequestHandler обрабатывает HTTP запросы для операций с заявками на поиск крови
type BloodRequestHandler struct {
	svc services.BloodSearchService
}

// NewBloodRequestHandler создает новый обработчик для заявок на поиск крови
func NewBloodRequestHandler(service services.BloodSearchService) *BloodRequestHandler {
	return &BloodRequestHandler{
		svc: service,
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

// mapBloodRequestToDTO преобразует ENT модель заявки в DTO
func mapBloodRequestToDTO(req *ent.BloodSearchRequest, situatableDonors *int) dto.BloodSearchRequest {
	var applications []*dto.DonorApplication
	for _, response := range req.Edges.Responses {
		donor := response.Edges.Donor
		applications = append(applications, &dto.DonorApplication{
			ID:              response.ID,
			RequestID:       req.ID,
			DonorID:         donor.ID,
			DonorName:       donor.Name,
			DonorPhotos:     donor.PhotoUrls,
			DonorBloodGroup: donor.Edges.BloodGroupRef.BloodGroup,
			Amount:          calculateDonationAmount(donor),
			WarnFactors:     donor.WarnFactors,
			Conditions:      response.Conditions,
			Status:          dto.DonorResponseStatus(response.Status),
			CreatedAt:       &response.CreatedAt,
			UpdatedAt:       &response.UpdatedAt,
		})
	}
	return dto.BloodSearchRequest{
		ID:                     req.ID,
		PetID:                  req.PetID,
		BloodVolumeNeeded:      req.BloodVolumeNeeded,
		BloodVolumeReserved:    req.BloodVolumeReserved,
		Regions:                req.Regions,
		SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
		Description:            req.Description,
		PhotoUrls:              req.PhotoUrls,
		BloodGroupNames:        req.BloodGroupNames,
		BloodComponentIds:      req.BloodComponentIds,
		OnBoarding:             req.OnBoarding,
		Status:                 dto.BloodSearchRequestStatus(req.Status),
		Responses:              applications,
		SuitableDonors:         *situatableDonors,
		CreatedAt:              &req.CreatedAt,
		UpdatedAt:              &req.UpdatedAt,
		DeletedAt:              req.DeletedAt,
	}
}

// mapDTOToBloodRequest преобразует DTO создания заявки в ENT модель
func mapDTOToBloodRequest(d dto.CreateBloodSearchRequest) *ent.BloodSearchRequest {
	return &ent.BloodSearchRequest{
		PetID:                  d.PetID,
		BloodVolumeNeeded:      d.BloodVolumeNeeded,
		Regions:                d.Regions,
		SmallPetsNotifyAllowed: d.SmallPetsNotifyAllowed,
		Description:            d.Description,
		BloodGroupNames:        d.BloodGroupNames,
		BloodComponentIds:      d.BloodComponentIds,
		Status:                 bloodsearchrequest.StatusActive,
	}
}

// Handlers

func (h *BloodRequestHandler) AddPetToBloodRequestPool(ctx context.Context, input *struct {
	Body dto.CreateBloodSearchRequest
}) (*dto.BloodRequestCreateResponse, error) {
	slog.DebugContext(ctx, "adding pet to blood request pool", "pet_id", input.Body.PetID)

	bloodReq := new(ent.CreateBloodSearchRequestInput)
	if err := copier.Copy(bloodReq, input.Body); err != nil {
		return nil, err
	}

	result, err := h.svc.CreateRequest(ctx, bloodReq)
	if err != nil {
		return nil, err
	}

	return &dto.BloodRequestCreateResponse{Body: dto.BloodSearchPetResponse{
		ID:     result.ID,
		PetID:  result.PetID,
		Status: dto.BloodSearchRequestStatus(result.Status),
	}}, nil
}

func (h *BloodRequestHandler) ApplyForBloodRequest(ctx context.Context, input *struct {
	dto.IDPathStr
	Body dto.DonorApplicationCreate
}) (*dto.DonorApplicationCreateResponse, error) {
	slog.DebugContext(ctx, "applying for blood request", "request_id", input.IDPathStr.ID, "donor_id", input.Body.DonorID)

	resp, err := h.svc.ApplyForBloodRequest(ctx, input.IDPathStr.ID, input.Body.DonorID, input.Body.Conditions)
	if err != nil {
		return nil, err
	}

	return &dto.DonorApplicationCreateResponse{
		Body: dto.DonorApplicationResponse{
			ID:      resp.ID,
			ReqID:   resp.Edges.Request.ID,
			DonorID: resp.Edges.Donor.ID,
			Status:  dto.DonorResponseStatus(resp.Status),
		},
	}, nil
}

func (h *BloodRequestHandler) UpdateBloodRequest(ctx context.Context, input *struct {
	dto.IDPathStr
	Body dto.UpdateBloodRequestDTO
}) (*dto.BloodRequestUpdateResponse, error) {
	slog.DebugContext(ctx, "updating blood request", "request_id", input.IDPathStr.ID)

	// Получить текущий объект
	existing, err := h.svc.GetRequestByID(ctx, input.IDPathStr.ID)
	if err != nil {
		return nil, err
	}

	// Частично обновить поля
	if input.Body.PetID != nil {
		existing.PetID = *input.Body.PetID
	}
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
	if len(input.Body.PhotoUrls) > 0 {
		existing.PhotoUrls = input.Body.PhotoUrls
	}
	if len(input.Body.BloodGroupNames) > 0 {
		existing.BloodGroupNames = input.Body.BloodGroupNames
	}
	if len(input.Body.BloodComponentIds) > 0 {
		existing.BloodComponentIds = input.Body.BloodComponentIds
	}
	if len(input.Body.OnBoarding) > 0 {
		existing.OnBoarding = input.Body.OnBoarding
	}
	if input.Body.Status != nil {
		existing.Status = bloodsearchrequest.Status(*input.Body.Status)
	}

	// Создать update input из обновленного existing
	updateReq := new(ent.UpdateBloodSearchRequestInput)
	if err := copier.Copy(updateReq, existing); err != nil {
		return nil, err
	}

	result, err := h.svc.UpdateRequest(ctx, input.IDPathStr.ID, updateReq)
	if err != nil {
		return nil, err
	}

	return &dto.BloodRequestUpdateResponse{Body: dto.BloodRequestUpdateResponseBody{
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

func (h *BloodRequestHandler) GetBloodRequestByID(ctx context.Context, input *dto.IDPathStr) (*dto.BloodRequestResponse, error) {
	bloodReq, err := h.svc.GetRequestByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	//TODO Метод GetRequestByID Также нужно исправить.
	return &dto.BloodRequestResponse{Body: mapBloodRequestToDTO(bloodReq, new(0))}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByPetID(ctx context.Context, input *dto.IDPathStr) (*dto.BloodRequestResponse, error) {
	slog.DebugContext(ctx, "getting blood request by PET ID", "pet_id", input.ID)
	bloodReq, situatableDonors, err := h.svc.GetRequestByPetID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.BloodRequestResponse{Body: mapBloodRequestToDTO(bloodReq, &situatableDonors)}, nil
}

// func (h *BloodRequestHandler) GetDonorsByID(ctx context.Context, input *dto.IDPathStr) (*dto.PetsResponse, error) {
// 	slog.DebugContext(ctx, "getting donors by ID", "request_id", input.ID)

// 	// Create a slice of dto.Pet
// 	pets := []dto.Pet{mocks.Pet1, mocks.Pet2}

// 	return &dto.PetsResponse{Body: pets}, nil
// }

func (h *BloodRequestHandler) DeleteBloodRequest(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "deleting blood request", "request_id", input.ID)
	if err := h.svc.DeleteRequest(ctx, input.ID); err != nil {
		return nil, err
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "Заявка удалена"
	return resp, nil
}
