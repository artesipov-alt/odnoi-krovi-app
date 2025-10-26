package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// BloodStockHandler обрабатывает HTTP запросы для операций с запасами крови
type BloodStockHandler struct {
	bloodStockService services.BloodStockService
}

// NewBloodStockHandler создает новый обработчик запасов крови
func NewBloodStockHandler(bloodStockService services.BloodStockService) *BloodStockHandler {
	return &BloodStockHandler{
		bloodStockService: bloodStockService,
	}
}

// GetAllBloodStocksHandler godoc
// @Summary Получение всех запасов крови
// @Description Возвращает список всех запасов крови в системе
// @Tags blood-stocks
// @Produce json
// @Success 200 {array} models.BloodStock "Список запасов крови"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks [get]
func (h *BloodStockHandler) GetAllBloodStocksHandler(c *fiber.Ctx) error {
	logger.Log.Info("получение всех запасов крови")

	stocks, err := h.bloodStockService.GetAll(c.Context())
	if err != nil {
		return err
	}

	return SendJSON(c, stocks)
}

// GetBloodStockByIDHandler godoc
// @Summary Получение запаса крови по ID
// @Description Возвращает информацию о конкретном запасе крови
// @Tags blood-stocks
// @Produce json
// @Param id path int true "ID запаса крови"
// @Success 200 {object} models.BloodStock "Запас крови"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Запас крови не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks/{id} [get]
func (h *BloodStockHandler) GetBloodStockByIDHandler(c *fiber.Ctx) error {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение запаса крови", zap.Int("stockId", id))

	stock, err := h.bloodStockService.GetByID(c.Context(), id)
	if err != nil {
		return err
	}

	return SendJSON(c, stock)
}

// GetBloodStocksByClinicIDHandler godoc
// @Summary Получение запасов крови клиники
// @Description Возвращает все запасы крови для конкретной клиники
// @Tags blood-stocks
// @Produce json
// @Param clinic_id path int true "ID клиники"
// @Success 200 {array} models.BloodStock "Список запасов крови клиники"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Клиника не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks/clinic/{clinic_id} [get]
func (h *BloodStockHandler) GetBloodStocksByClinicIDHandler(c *fiber.Ctx) error {
	clinicID, err := ParseIDParam(c, "clinic_id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение запасов крови клиники", zap.Int("clinicId", clinicID))

	stocks, err := h.bloodStockService.GetByClinicID(c.Context(), clinicID)
	if err != nil {
		return err
	}

	return SendJSON(c, stocks)
}

// GetBloodStocksByBloodTypeIDHandler godoc
// @Summary Получение запасов крови по типу крови
// @Description Возвращает все запасы крови для конкретного типа крови
// @Tags blood-stocks
// @Produce json
// @Param blood_type_id path int true "ID типа крови"
// @Success 200 {array} models.BloodStock "Список запасов крови"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Тип крови не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks/blood-type/{blood_type_id} [get]
func (h *BloodStockHandler) GetBloodStocksByBloodTypeIDHandler(c *fiber.Ctx) error {
	bloodTypeID, err := ParseIDParam(c, "blood_type_id")
	if err != nil {
		return err
	}

	logger.Log.Info("получение запасов крови по типу", zap.Int("bloodTypeId", bloodTypeID))

	stocks, err := h.bloodStockService.GetByBloodTypeID(c.Context(), bloodTypeID)
	if err != nil {
		return err
	}

	return SendJSON(c, stocks)
}

// SearchBloodStocksHandler godoc
// @Summary Поиск запасов крови с фильтрами
// @Description Выполняет поиск запасов крови по различным параметрам (клиника, тип животного, тип крови, статус, объем, цена)
// @Tags blood-stocks
// @Accept json
// @Produce json
// @Param clinic_id query int false "ID клиники"
// @Param pet_type query string false "Тип животного (dog/cat)"
// @Param blood_type_id query int false "ID типа крови"
// @Param status query string false "Статус (active/reserved/used/expired)"
// @Param min_volume query int false "Минимальный объем (мл)"
// @Param max_volume query int false "Максимальный объем (мл)"
// @Param min_price query number false "Минимальная цена (руб)"
// @Param max_price query number false "Максимальная цена (руб)"
// @Success 200 {array} models.BloodStock "Список найденных запасов крови"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks/search [get]
func (h *BloodStockHandler) SearchBloodStocksHandler(c *fiber.Ctx) error {
	filters := repositories.BloodStockFilters{}

	// Парсим опциональные параметры
	if clinicID, err := ParseOptionalIntQuery(c, "clinic_id"); err != nil {
		return err
	} else if clinicID != nil {
		filters.ClinicID = clinicID
	}

	if petTypeStr := c.Query("pet_type"); petTypeStr != "" {
		petType := models.PetType(petTypeStr)
		filters.PetType = &petType
	}

	if bloodTypeID, err := ParseOptionalIntQuery(c, "blood_type_id"); err != nil {
		return err
	} else if bloodTypeID != nil {
		filters.BloodTypeID = bloodTypeID
	}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.BloodStockStatus(statusStr)
		filters.Status = &status
	}

	if minVolume, err := ParseOptionalIntQuery(c, "min_volume"); err != nil {
		return err
	} else if minVolume != nil {
		filters.MinVolume = minVolume
	}

	if maxVolume, err := ParseOptionalIntQuery(c, "max_volume"); err != nil {
		return err
	} else if maxVolume != nil {
		filters.MaxVolume = maxVolume
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return apperrors.BadRequest("неверный формат min_price")
		}
		filters.MinPrice = &minPrice
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return apperrors.BadRequest("неверный формат max_price")
		}
		filters.MaxPrice = &maxPrice
	}

	logger.Log.Info("поиск запасов крови с фильтрами")

	stocks, err := h.bloodStockService.Search(c.Context(), filters)
	if err != nil {
		return err
	}

	return SendJSON(c, stocks)
}

// CreateBloodStockHandler godoc
// @Summary Создание нового запаса крови
// @Description Создает новый запас крови в системе
// @Tags blood-stocks
// @Accept json
// @Produce json
// @Param request body services.BloodStockCreate true "Данные запаса крови"
// @Success 201 {object} models.BloodStock "Созданный запас крови"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Клиника или тип крови не найдены"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks [post]
func (h *BloodStockHandler) CreateBloodStockHandler(c *fiber.Ctx) error {
	var stockData services.BloodStockCreate
	if err := ParseBody(c, &stockData); err != nil {
		return err
	}

	logger.Log.Info("создание запаса крови", zap.String("petType", string(stockData.PetType)))

	stock, err := h.bloodStockService.CreateBloodStock(c.Context(), stockData)
	if err != nil {
		return err
	}

	return SendCreated(c, stock)
}

// UpdateBloodStockHandler godoc
// @Summary Обновление запаса крови
// @Description Обновляет информацию о запасе крови
// @Tags blood-stocks
// @Accept json
// @Produce json
// @Param id path int true "ID запаса крови"
// @Param request body services.BloodStockUpdate true "Данные для обновления"
// @Success 200 {object} SuccessResponse "Запас крови успешно обновлен"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Запас крови не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks/{id} [put]
func (h *BloodStockHandler) UpdateBloodStockHandler(c *fiber.Ctx) error {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	var updateData services.BloodStockUpdate
	if err := ParseBody(c, &updateData); err != nil {
		return err
	}

	logger.Log.Info("обновление запаса крови", zap.Int("stockId", id))

	if err := h.bloodStockService.UpdateBloodStock(c.Context(), id, updateData); err != nil {
		return err
	}

	return SendSuccess(c, "Запас крови успешно обновлен")
}

// DeleteBloodStockHandler godoc
// @Summary Удаление запаса крови
// @Description Удаляет запас крови из системы
// @Tags blood-stocks
// @Produce json
// @Param id path int true "ID запаса крови"
// @Success 200 {object} SuccessResponse "Запас крови успешно удален"
// @Failure 400 {object} ErrorResponse "Неверный запрос"
// @Failure 404 {object} ErrorResponse "Запас крови не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /blood-stocks/{id} [delete]
func (h *BloodStockHandler) DeleteBloodStockHandler(c *fiber.Ctx) error {
	id, err := ParseIDParam(c, "id")
	if err != nil {
		return err
	}

	logger.Log.Info("удаление запаса крови", zap.Int("stockId", id))

	if err := h.bloodStockService.DeleteBloodStock(c.Context(), id); err != nil {
		return err
	}

	return SendSuccess(c, "Запас крови успешно удален")
}
