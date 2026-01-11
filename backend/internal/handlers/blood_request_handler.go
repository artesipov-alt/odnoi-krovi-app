package handlers

// import (
// 	"net/http"

// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils"
// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/logger"
// 	"github.com/labstack/echo/v4"
// 	"go.uber.org/zap"
// )

// // BloodRequestHandler обрабатывает HTTP запросы для операций с заявками на поиск крови
// type BloodRequestHandler struct {
// 	service services.BloodSearchService
// }

// // NewBloodRequestHandler создает новый обработчик для заявок на поиск крови
// func NewBloodRequestHandler(service services.BloodSearchService) *BloodRequestHandler {
// 	return &BloodRequestHandler{
// 		service: service,
// 	}
// }

// // AddPetToBloodRequestPool godoc
// // @Summary Добавить питомца в пул поиска крови
// // @Description Создает новую заявку на поиск крови для питомца
// // @Tags blood-request
// // @Accept json
// // @Produce json
// // @Param request body dto.BloodSearchPetRequest true "Данные заявки"
// // @Success 201 {object} dto.BloodSearchPetResponse "Созданная заявка"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /v1/blood-request/pool [post]
// func (h *BloodRequestHandler) AddPetToBloodRequestPool(c echo.Context) error {
// 	var req dto.BloodSearchPetRequest
// 	if err := utils.ParseBody(c, &req); err != nil {
// 		return err
// 	}

// 	logger.Log.Info("запрос на создание заявки на поиск крови", zap.String("petId", req.PetID))

// 	result, err := h.service.CreateRequest(c.Request().Context(), req)
// 	if err != nil {
// 		return err // apperrors handled by middleware
// 	}

// 	return utils.SendCreated(c, result)
// }

// // GetPetsFromBloodRequestPool godoc
// // @Summary Получить список заявок на поиск крови
// // @Description Возвращает список заявок по фильтрам
// // @Tags blood-request
// // @Accept json
// // @Produce json
// // @Param request body dto.BloodSearchFilterRequest true "Фильтры поиска"
// // @Success 200 {object} dto.BloodSearchPetsResponse "Список заявок"
// // @Failure 400 {object} utils.ErrorResponse "Неверный запрос"
// // @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// // @Router /v1/blood-request/pool/search [post]
// func (h *BloodRequestHandler) GetPetsFromBloodRequestPool(c echo.Context) error {
// 	var filterReq dto.BloodSearchFilterRequest
// 	if err := utils.ParseBody(c, &filterReq); err != nil {
// 		return err
// 	}

// 	filters := make(map[string]interface{})
// 	if filterReq.PetID != "" {
// 		filters["pet_id"] = filterReq.PetID
// 	}
// 	if filterReq.Status != "" {
// 		filters["status"] = filterReq.Status
// 	}

// 	logger.Log.Info("получение списка заявок на поиск крови", zap.Any("filters", filters))

// 	requests, err := h.service.ListRequests(c.Request().Context(), filterReq.Limit, filterReq.Offset, filters)
// 	if err != nil {
// 		return err
// 	}

// 	return utils.SendJSON(c, dto.BloodSearchPetsResponse{
// 		Requests: requests,
// 	})
// }

// // GetBloodRequestByID godoc
// // @Summary Получить заявку по ID
// // @Description Возвращает информацию о конкретной заявке
// // @Tags blood-request
// // @Produce json
// // @Param id path string true "ID заявки"
// // @Success 200 {object} dto.BloodSearchRequestDTO "Данные заявки"
// // @Router /v1/blood-request/{id} [get]
// func (h *BloodRequestHandler) GetBloodRequestByID(c echo.Context) error {
// 	id := c.Param("id")
// 	if id == "" {
// 		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
// 	}

// 	result, err := h.service.GetRequestByID(c.Request().Context(), id)
// 	if err != nil {
// 		return err
// 	}

// 	return utils.SendJSON(c, result)
// }

// // DeleteBloodRequest godoc
// // @Summary Удалить заявку
// // @Description Удаляет заявку на поиск крови (soft delete)
// // @Tags blood-request
// // @Param id path string true "ID заявки"
// // @Success 200 {object} utils.SuccessResponse
// // @Router /v1/blood-request/{id} [delete]
// func (h *BloodRequestHandler) DeleteBloodRequest(c echo.Context) error {
// 	id := c.Param("id")
// 	if id == "" {
// 		return echo.NewHTTPError(http.StatusBadRequest, "ID is required")
// 	}

// 	if err := h.service.DeleteRequest(c.Request().Context(), id); err != nil {
// 		return err
// 	}

// 	return utils.SendSuccess(c, "Заявка успешно удалена")
// }
