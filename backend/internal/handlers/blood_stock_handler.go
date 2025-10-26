package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
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
		logger.Log.Error("не удалось получить запасы крови", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(stocks)
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
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID запаса крови обязателен",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID запаса крови",
		})
	}

	logger.Log.Info("получение запаса крови", zap.Int("stockId", id))

	stock, err := h.bloodStockService.GetByID(c.Context(), id)
	if err != nil {
		logger.Log.Error("не удалось получить запас крови", zap.Error(err), zap.Int("stockId", id))

		if err.Error() == "запас крови не найден" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(stock)
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
	clinicIDStr := c.Params("clinic_id")
	if clinicIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID клиники обязателен",
		})
	}

	clinicID, err := strconv.Atoi(clinicIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID клиники",
		})
	}

	logger.Log.Info("получение запасов крови клиники", zap.Int("clinicId", clinicID))

	stocks, err := h.bloodStockService.GetByClinicID(c.Context(), clinicID)
	if err != nil {
		logger.Log.Error("не удалось получить запасы крови клиники", zap.Error(err), zap.Int("clinicId", clinicID))

		if err.Error() == "клиника не найдена" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(stocks)
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
	bloodTypeIDStr := c.Params("blood_type_id")
	if bloodTypeIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID типа крови обязателен",
		})
	}

	bloodTypeID, err := strconv.Atoi(bloodTypeIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID типа крови",
		})
	}

	logger.Log.Info("получение запасов крови по типу", zap.Int("bloodTypeId", bloodTypeID))

	stocks, err := h.bloodStockService.GetByBloodTypeID(c.Context(), bloodTypeID)
	if err != nil {
		logger.Log.Error("не удалось получить запасы крови по типу", zap.Error(err), zap.Int("bloodTypeId", bloodTypeID))

		if err.Error() == "тип крови не найден" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(stocks)
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

	// Парсим clinic_id
	if clinicIDStr := c.Query("clinic_id"); clinicIDStr != "" {
		clinicID, err := strconv.Atoi(clinicIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Неверный формат clinic_id",
			})
		}
		filters.ClinicID = &clinicID
	}

	// Парсим pet_type
	if petTypeStr := c.Query("pet_type"); petTypeStr != "" {
		petType := models.PetType(petTypeStr)
		filters.PetType = &petType
	}

	// Парсим blood_type_id
	if bloodTypeIDStr := c.Query("blood_type_id"); bloodTypeIDStr != "" {
		bloodTypeID, err := strconv.Atoi(bloodTypeIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Неверный формат blood_type_id",
			})
		}
		filters.BloodTypeID = &bloodTypeID
	}

	// Парсим status
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.BloodStockStatus(statusStr)
		filters.Status = &status
	}

	// Парсим min_volume
	if minVolumeStr := c.Query("min_volume"); minVolumeStr != "" {
		minVolume, err := strconv.Atoi(minVolumeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Неверный формат min_volume",
			})
		}
		filters.MinVolume = &minVolume
	}

	// Парсим max_volume
	if maxVolumeStr := c.Query("max_volume"); maxVolumeStr != "" {
		maxVolume, err := strconv.Atoi(maxVolumeStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Неверный формат max_volume",
			})
		}
		filters.MaxVolume = &maxVolume
	}

	// Парсим min_price
	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Неверный формат min_price",
			})
		}
		filters.MinPrice = &minPrice
	}

	// Парсим max_price
	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: "Неверный формат max_price",
			})
		}
		filters.MaxPrice = &maxPrice
	}

	logger.Log.Info("поиск запасов крови с фильтрами")

	stocks, err := h.bloodStockService.Search(c.Context(), filters)
	if err != nil {
		logger.Log.Error("не удалось выполнить поиск запасов крови", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(stocks)
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
	if err := c.BodyParser(&stockData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверное тело запроса",
		})
	}

	// Валидация данных
	if err := validation.ValidateStruct(stockData); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	logger.Log.Info("создание запаса крови", zap.String("petType", string(stockData.PetType)))

	stock, err := h.bloodStockService.CreateBloodStock(c.Context(), stockData)
	if err != nil {
		logger.Log.Error("не удалось создать запас крови", zap.Error(err))

		if err.Error() == "клиника не найдена" || err.Error() == "тип крови не найден" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.Status(fiber.StatusCreated).JSON(stock)
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
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID запаса крови обязателен",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID запаса крови",
		})
	}

	var updateData services.BloodStockUpdate
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверное тело запроса",
		})
	}

	// Валидация данных
	if err := validation.ValidateStruct(updateData); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(validationErrors)
	}

	logger.Log.Info("обновление запаса крови", zap.Int("stockId", id))

	if err := h.bloodStockService.UpdateBloodStock(c.Context(), id, updateData); err != nil {
		logger.Log.Error("не удалось обновить запас крови", zap.Error(err), zap.Int("stockId", id))

		if err.Error() == "запас крови не найден" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(SuccessResponse{
		Message: "Запас крови успешно обновлен",
	})
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
	idStr := c.Params("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "ID запаса крови обязателен",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Неверный формат ID запаса крови",
		})
	}

	logger.Log.Info("удаление запаса крови", zap.Int("stockId", id))

	if err := h.bloodStockService.DeleteBloodStock(c.Context(), id); err != nil {
		logger.Log.Error("не удалось удалить запас крови", zap.Error(err), zap.Int("stockId", id))

		if err.Error() == "запас крови не найден" {
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: err.Error(),
		})
	}

	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(SuccessResponse{
		Message: "Запас крови успешно удален",
	})
}
