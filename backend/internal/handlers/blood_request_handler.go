package handlers

import (
	"context"
	"log/slog"
	"net/http"

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
		Tags:          []string{"blood-request"},
		DefaultStatus: http.StatusCreated,
	}, h.AddPetToBloodRequestPool)

	// Получить список заявок на поиск крови
	huma.Register(api, huma.Operation{
		OperationID: "get-pets-from-blood-request-pool",
		Method:      http.MethodPost,
		Path:        "/v1/blood-request/pool/search",
		Summary:     "Получить список заявок на поиск крови",
		Description: "Возвращает список заявок по фильтрам",
		Tags:        []string{"blood-request"},
	}, h.GetPetsFromBloodRequestPool)

	// Получить заявку по ID
	huma.Register(api, huma.Operation{
		OperationID: "get-blood-request-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Получить заявку по ID",
		Description: "Возвращает информацию о конкретной заявке",
		Tags:        []string{"blood-request"},
	}, h.GetBloodRequestByID)

	// Удалить заявку
	huma.Register(api, huma.Operation{
		OperationID: "delete-blood-request",
		Method:      http.MethodDelete,
		Path:        "/v1/blood-request/{id}",
		Summary:     "Удалить заявку",
		Description: "Удаляет заявку на поиск крови (soft delete)",
		Tags:        []string{"blood-request"},
	}, h.DeleteBloodRequest)
}

// Вспомогательные структуры для Huma

type BloodRequestIDPath struct {
	ID string `path:"id" doc:"ID заявки" minLength:"1" example:"BR-25-000001"`
}

type BloodSearchPetResponseWrapper struct {
	Body dto.BloodSearchPetResponse
}

type BloodSearchPetsResponseWrapper struct {
	Body dto.BloodSearchPetsResponse
}

type BloodSearchRequestDTOWrapper struct {
	Body dto.BloodSearchRequestDTO
}

// Handlers

func (h *BloodRequestHandler) AddPetToBloodRequestPool(ctx context.Context, input *struct {
	Body dto.BloodSearchPetRequest
}) (*BloodSearchPetResponseWrapper, error) {
	slog.InfoContext(ctx, "Создание заявки на поиск крови", "pet_id", input.Body.PetID)

	result, err := h.service.CreateRequest(ctx, input.Body)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка создания заявки на поиск крови", "pet_id", input.Body.PetID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Заявка на поиск крови успешно создана", "pet_id", input.Body.PetID, "request_id", result.ID)
	return &BloodSearchPetResponseWrapper{Body: result}, nil
}

func (h *BloodRequestHandler) GetPetsFromBloodRequestPool(ctx context.Context, input *struct {
	Body dto.BloodSearchFilterRequest
}) (*BloodSearchPetsResponseWrapper, error) {
	filters := make(map[string]interface{})
	if input.Body.PetID != "" {
		filters["pet_id"] = input.Body.PetID
	}
	if input.Body.Status != "" {
		filters["status"] = input.Body.Status
	}

	slog.InfoContext(ctx, "Получение списка заявок на поиск крови", "filters", filters)

	requests, err := h.service.ListRequests(ctx, input.Body.Limit, input.Body.Offset, filters)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка получения списка заявок на поиск крови", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Список заявок на поиск крови успешно получен", "count", len(requests))
	return &BloodSearchPetsResponseWrapper{Body: dto.BloodSearchPetsResponse{
		Requests: requests,
	}}, nil
}

func (h *BloodRequestHandler) GetBloodRequestByID(ctx context.Context, input *BloodRequestIDPath) (*BloodSearchRequestDTOWrapper, error) {
	slog.InfoContext(ctx, "Получение заявки по ID", "request_id", input.ID)

	result, err := h.service.GetRequestByID(ctx, input.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка получения заявки по ID", "request_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Заявка по ID успешно получена", "request_id", input.ID)
	return &BloodSearchRequestDTOWrapper{Body: result}, nil
}

func (h *BloodRequestHandler) DeleteBloodRequest(ctx context.Context, input *BloodRequestIDPath) (*MessageResponse, error) {
	slog.InfoContext(ctx, "Удаление заявки", "request_id", input.ID)

	if err := h.service.DeleteRequest(ctx, input.ID); err != nil {
		slog.ErrorContext(ctx, "Ошибка удаления заявки", "request_id", input.ID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "Заявка успешно удалена", "request_id", input.ID)
	resp := &MessageResponse{}
	resp.Body.Message = "Заявка успешно удалена"
	return resp, nil
}
