package handlers

import (
	"context"
	"errors"
	"log/slog" // Import slog
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/danielgtaylor/huma/v2"
)

// BloodRequestHandler обрабатывает HTTP запросы для операций с заявками на поиск крови
type BloodRequestHandler struct {
	service services.BloodSearchService
}

// NewBloodRequestHandler создает новый обработчик для заявок на поиск крови
func NewBloodRequestHandler(service services.BloodSearchService) *BloodRequestHandler {
	return &BloodRequestHandler{
		service: service,
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

	// Получить список заявок на поиск крови
	huma.Register(api, huma.Operation{
		OperationID: "get-pets-from-blood-request-pool",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/pool/search",
		Summary:     "Получить список заявок на поиск крови",
		Description: "Возвращает список заявок по фильтрам",
		Tags:        []string{"blood-request-v1"},
	}, h.GetPetsFromBloodRequestPool)

	// Получить заявку по ID
	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Получить заявку по ID",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByID)

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

// Вспомогательные структуры для Huma

// mapBloodRequestToDTO преобразует ENT модель заявки в DTO
func mapBloodRequestToDTO(req *ent.BloodSearchRequest) dto.BloodSearchPetRequest {
	return dto.BloodSearchPetRequest{
		PetID:                  req.PetID,
		BloodVolumeNeeded:      req.BloodVolumeNeeded,
		BloodVolumeReserved:    req.BloodVolumeReserved,
		Regions:                req.Regions,
		SmallPetsNotifyAllowed: req.SmallPetsNotifyAllowed,
		Description:            req.Description,
		PhotoUrls:              req.PhotoUrls,
		BloodGroupNames:        req.BloodGroupNames,
		BloodComponentIds:      req.BloodComponentIds,
		Status:                 dto.BloodSearchRequestStatus(req.Status),
	}
}

// mapDTOToBloodRequest преобразует DTO создания заявки в ENT модель
func mapDTOToBloodRequest(d dto.BloodSearchPetRequest) *ent.BloodSearchRequest {
	return &ent.BloodSearchRequest{
		PetID:                  d.PetID,
		BloodVolumeNeeded:      d.BloodVolumeNeeded,
		BloodVolumeReserved:    d.BloodVolumeReserved,
		Regions:                d.Regions,
		SmallPetsNotifyAllowed: d.SmallPetsNotifyAllowed,
		Description:            d.Description,
		PhotoUrls:              d.PhotoUrls,
		BloodGroupNames:        d.BloodGroupNames,
		BloodComponentIds:      d.BloodComponentIds,
		Status:                 bloodsearchrequest.StatusActive,
	}
}

// Handlers

func (h *BloodRequestHandler) AddPetToBloodRequestPool(ctx context.Context, input *struct {
	Body dto.BloodSearchPetRequest
}) (*dto.BloodRequestCreateResponse, error) {
	bloodReq := mapDTOToBloodRequest(input.Body)

	result, err := h.service.CreateRequest(ctx, bloodReq)
	if err != nil {
		if errors.Is(err, apperrors.ErrBloodRequestAlreadyExists) {
			slog.DebugContext(ctx, "blood request already exists for pet", "pet_id", input.Body.PetID, "error", err.Error())
			return nil, huma.Error409Conflict("Заявка на поиск крови для этого питомца уже существует")
		}
		slog.ErrorContext(ctx, "failed to create blood request", "pet_id", input.Body.PetID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &dto.BloodRequestCreateResponse{Body: dto.BloodSearchPetResponse{
		ID:     result.ID,
		PetID:  result.PetID,
		Status: dto.BloodSearchRequestStatus(result.Status),
	}}, nil
}

func (h *BloodRequestHandler) GetPetsFromBloodRequestPool(ctx context.Context, input *struct {
	Body dto.BloodSearchFilterRequest
}) (*dto.BloodRequestsResponse, error) {
	filters := make(map[string]any)
	if input.Body.PetID != "" {
		filters["pet_id"] = input.Body.PetID
	}
	if input.Body.Status != "" {
		filters["status"] = input.Body.Status
	}

	requests, err := h.service.ListRequests(ctx, input.Body.Limit, input.Body.Offset, filters)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list blood requests", "filters", filters, "error", err.Error())
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	dtos := make([]dto.BloodSearchPetRequest, len(requests))
	for i, req := range requests {
		dtos[i] = mapBloodRequestToDTO(req)
	}

	return &dto.BloodRequestsResponse{Body: dtos}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByID(ctx context.Context, input *dto.IDPath) (*dto.BloodRequestResponse, error) {
	result, err := h.service.GetRequestByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			slog.DebugContext(ctx, "blood request not found", "request_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Заявка не найдена")
		}
		slog.ErrorContext(ctx, "failed to get blood request by ID", "request_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	return &dto.BloodRequestResponse{Body: mapBloodRequestToDTO(result)}, nil
}

func (h *BloodRequestHandler) DeleteBloodRequest(ctx context.Context, input *dto.IDPath) (*dto.MessageResponse, error) {
	if err := h.service.DeleteRequest(ctx, input.ID); err != nil {
		if errors.Is(err, apperrors.ErrBloodRequestNotFound) {
			slog.DebugContext(ctx, "blood request not found for deletion", "request_id", input.ID, "error", err.Error())
			return nil, huma.Error404NotFound("Заявка не найдена")
		}
		slog.ErrorContext(ctx, "failed to delete blood request", "request_id", input.ID, "error", err.Error())
		return nil, huma.Error500InternalServerError("Ошибка сервера")
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "Заявка успешно удалена"
	return resp, nil
}
