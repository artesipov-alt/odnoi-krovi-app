package handlers

// "github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
// "github.com/artesipov-alt/odnoi-krovi-app/internal/services"
// "github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
// "github.com/gofiber/fiber/v2"
// "go.uber.org/zap"

// // LocationHandler обрабатывает HTTP запросы для операций с локациями
// type LocationHandler struct {
// 	locationService services.LocationService
// }

// // LocationCreateRequest представляет запрос на создание локации
// type LocationCreateRequest struct {
// 	Name string `json:"name" validate:"required,min=2,max=255" example:"Москва"`
// }

// // LocationUpdateRequest представляет запрос на обновление локации
// type LocationUpdateRequest struct {
// 	Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=255" example:"Санкт-Петербург"`
// }

// // NewLocationHandler создает новый обработчик локаций
// func NewLocationHandler(locationService services.LocationService) *LocationHandler {
// 	return &LocationHandler{
// 		locationService: locationService,
// 	}
// }

// // CreateLocationHandler godoc
// // @Summary Создание новой локации
// // @Description Создает новую локацию в системе
// // @Tags locations
// // @Accept json
// // @Produce json
// // @Param request body LocationCreateRequest true "Данные для создания локации"
// // @Success 201 {object} models.Location "Созданная локация"
// // @Failure 400 {object} ErrorResponse "Неверный запрос"
// // @Failure 409 {object} ErrorResponse "Локация с таким названием уже существует"
// // @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// // @Router /locations [post]
// func (h *LocationHandler) CreateLocationHandler(c *fiber.Ctx) error {
// 	var request LocationCreateRequest
// 	if err := ParseBody(c, &request); err != nil {
// 		return err
// 	}

// 	logger.Log.Info("создание локации", zap.String("name", request.Name))

// 	locationData := services.LocationCreate{
// 		Name: request.Name,
// 	}

// 	location, err := h.locationService.CreateLocation(c.Context(), locationData)
// 	if err != nil {
// 		return err
// 	}

// 	return SendCreated(c, location)
// }

// // GetLocationHandler godoc
// // @Summary Получение локации по ID
// // @Description Возвращает информацию о локации по её идентификатору
// // @Tags locations
// // @Produce json
// // @Param id path int true "ID локации"
// // @Success 200 {object} models.Location "Данные локации"
// // @Failure 400 {object} ErrorResponse "Неверный запрос"
// // @Failure 404 {object} ErrorResponse "Локация не найдена"
// // @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// // @Router /locations/{id} [get]
// func (h *LocationHandler) GetLocationHandler(c *fiber.Ctx) error {
// 	id, err := ParseIDParam(c, "id")
// 	if err != nil {
// 		return err
// 	}

// 	logger.Log.Info("получение локации", zap.Int("locationId", id))

// 	location, err := h.locationService.GetLocationByID(c.Context(), id)
// 	if err != nil {
// 		return err
// 	}

// 	return SendJSON(c, location)
// }

// // GetAllLocationsHandler godoc
// // @Summary Получение всех локаций
// // @Description Возвращает список всех локаций в системе
// // @Tags locations
// // @Produce json
// // @Success 200 {array} models.Location "Список локаций"
// // @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// // @Router /locations [get]
// func (h *LocationHandler) GetAllLocationsHandler(c *fiber.Ctx) error {
// 	logger.Log.Info("получение всех локаций")

// 	locations, err := h.locationService.GetAllLocations(c.Context())
// 	if err != nil {
// 		return err
// 	}

// 	return SendJSON(c, locations)
// }

// // UpdateLocationHandler godoc
// // @Summary Обновление данных локации
// // @Description Обновляет информацию о локации
// // @Tags locations
// // @Accept json
// // @Produce json
// // @Param id path int true "ID локации"
// // @Param request body LocationUpdateRequest true "Данные для обновления"
// // @Success 200 {object} SuccessResponse "Данные успешно обновлены"
// // @Failure 400 {object} ErrorResponse "Неверный запрос"
// // @Failure 404 {object} ErrorResponse "Локация не найдена"
// // @Failure 409 {object} ErrorResponse "Локация с таким названием уже существует"
// // @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// // @Router /locations/{id} [put]
// func (h *LocationHandler) UpdateLocationHandler(c *fiber.Ctx) error {
// 	id, err := ParseIDParam(c, "id")
// 	if err != nil {
// 		return err
// 	}

// 	var updateData LocationUpdateRequest
// 	if err := ParseBody(c, &updateData); err != nil {
// 		return err
// 	}

// 	logger.Log.Info("обновление локации", zap.Int("locationId", id))

// 	updateServiceData := services.LocationUpdate{
// 		Name: updateData.Name,
// 	}

// 	if err := h.locationService.UpdateLocation(c.Context(), id, updateServiceData); err != nil {
// 		return err
// 	}

// 	return SendSuccess(c, "Локация успешно обновлена")
// }

// // DeleteLocationHandler godoc
// // @Summary Удаление локации по ID
// // @Description Удаляет локацию из системы
// // @Tags locations
// // @Produce json
// // @Param id path int true "ID локации"
// // @Success 200 {object} SuccessResponse "Локация успешно удалена"
// // @Failure 400 {object} ErrorResponse "Неверный запрос"
// // @Failure 404 {object} ErrorResponse "Локация не найдена"
// // @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// // @Router /locations/{id} [delete]
// func (h *LocationHandler) DeleteLocationHandler(c *fiber.Ctx) error {
// 	id, err := ParseIDParam(c, "id")
// 	if err != nil {
// 		return err
// 	}

// 	logger.Log.Info("удаление локации", zap.Int("locationId", id))

// 	if err := h.locationService.DeleteLocation(c.Context(), id); err != nil {
// 		return err
// 	}

// 	return SendSuccess(c, "Локация успешно удалена")
// }

// // GetLocationsReferenceHandler godoc
// // @Summary Получение справочника локаций
// // @Description Возвращает список всех локаций в формате справочника для выбора на фронтенде
// // @Tags reference
// // @Produce json
// // @Success 200 {object} ReferenceResponseDB "Список локаций в формате справочника"
// // @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// // @Router /reference/locations [get]
// func (h *LocationHandler) GetLocationsReferenceHandler(c *fiber.Ctx) error {
// 	logger.Log.Info("получение справочника локаций")

// 	locations, err := h.locationService.GetAllLocations(c.Context())
// 	if err != nil {
// 		logger.Log.Error("не удалось получить локации из БД", zap.Error(err))
// 		return apperrors.Internal(err, "не удалось получить список локаций")
// 	}

// 	items := make([]ReferenceItemDB, len(locations))
// 	for i, location := range locations {
// 		items[i] = ReferenceItemDB{
// 			Value: location.ID,
// 			Label: location.Name,
// 		}
// 	}

// 	c.Set("Content-Type", "application/json; charset=utf-8")
// 	return c.JSON(ReferenceResponseDB{Data: items})
// }
