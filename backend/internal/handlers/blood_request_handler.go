package handlers

import (
	"context"
	"log/slog" // Import slog
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/mocks"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jinzhu/copier"
)

// BloodSearchService определяет интерфейс для бизнес-логики заявок на поиск крови
type BloodSearchService interface {
	// CreateRequest создает новую заявку на поиск крови
	CreateRequest(ctx context.Context, bloodReq *ent.CreateBloodSearchRequestInput) (*ent.BloodSearchRequest, error)

	// GetRequestByID получает заявку по её ID
	GetRequestByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error)

	// GetRequestByPetID получает активную заявку для конкретного питомца
	GetRequestByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error)

	// UpdateRequest обновляет информацию о заявке
	UpdateRequest(ctx context.Context, id string, bloodReq *ent.UpdateBloodSearchRequestInput) (*ent.BloodSearchRequest, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status string) error

	// ExistsByID проверяет существование заявки по её ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// DeleteRequest удаляет заявку (soft delete)
	DeleteRequest(ctx context.Context, id string) error

	// ListRequests возвращает список заявок с фильтрацией
	ListRequests(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error)

	// buildFullPhotoURLs преобразует пути к фото в полные публичные URL
	BuildFullPhotoURLs(paths []string) []string
}

// BloodRequestHandler обрабатывает HTTP запросы для операций с заявками на поиск крови
type BloodRequestHandler struct {
	svc BloodSearchService
}

// NewBloodRequestHandler создает новый обработчик для заявок на поиск крови
func NewBloodRequestHandler(service BloodSearchService) *BloodRequestHandler {
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

	// Получить заявку по ID питомца
	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-pet-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/pet/{id}",
		Summary:     "Получить заявку по ID питомца",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request-v1"},
	}, h.GetBloodRequestByPetID)

	// Получить список доноров по ID заявки
	huma.Register(api, huma.Operation{
		OperationID: "get-donors-by-req-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/donors/{id}",
		Summary:     "Получить список доноров по ID заявки",
		Description: "Возвращает список доноров откликнувшихся на заявку",
		Tags:        []string{"blood-request-v1"},
	}, h.GetDonorsByID)

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
		CreatedAt:              &req.CreatedAt,
		UpdatedAt:              &req.UpdatedAt,
		DeletedAt:              req.DeletedAt,
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

func (h *BloodRequestHandler) GetPetsFromBloodRequestPool(ctx context.Context, input *struct {
	Body dto.BloodSearchFilterRequest
}) (*dto.BloodRequestsResponse, error) {
	slog.DebugContext(ctx, "getting pets from blood request pool", "filters", input.Body)
	filters := make(map[string]any)
	if input.Body.PetID != "" {
		filters["pet_id"] = input.Body.PetID
	}
	if input.Body.Status != "" {
		filters["status"] = input.Body.Status
	}

	requests, err := h.svc.ListRequests(ctx, input.Body.Limit, input.Body.Offset, filters)
	if err != nil {
		return nil, err
	}

	dtos := make([]dto.BloodSearchPetRequest, len(requests))
	for i, req := range requests {
		dtos[i] = mapBloodRequestToDTO(req)
	}

	return &dto.BloodRequestsResponse{Body: dtos}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByID(ctx context.Context, input *dto.IDPathStr) (*dto.BloodRequestResponse, error) {
	bloodReq, err := h.svc.GetRequestByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.BloodRequestResponse{Body: mapBloodRequestToDTO(bloodReq)}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByPetID(ctx context.Context, input *dto.IDPathStr) (*dto.BloodRequestResponse, error) {
	slog.DebugContext(ctx, "getting blood request by PET ID", "pet_id", input.ID)
	bloodReq, err := h.svc.GetRequestByPetID(ctx, input.ID)
	if err != nil {
		return nil, err
	}

	return &dto.BloodRequestResponse{Body: mapBloodRequestToDTO(bloodReq)}, nil
}

func (h *BloodRequestHandler) GetDonorsByID(ctx context.Context, input *dto.IDPathStr) (*dto.PetsResponse, error) {
	slog.DebugContext(ctx, "getting donors by ID", "request_id", input.ID)

	// Create a slice of dto.Pet
	pets := []dto.Pet{mocks.Pet1, mocks.Pet2}

	return &dto.PetsResponse{Body: pets}, nil
}

func (h *BloodRequestHandler) DeleteBloodRequest(ctx context.Context, input *dto.IDPathStr) (*dto.MessageResponse, error) {
	slog.DebugContext(ctx, "deleting blood request", "request_id", input.ID)
	if err := h.svc.DeleteRequest(ctx, input.ID); err != nil {
		return nil, err
	}

	resp := &dto.MessageResponse{}
	resp.Body.Message = "Заявка удалена"
	return resp, nil
}
