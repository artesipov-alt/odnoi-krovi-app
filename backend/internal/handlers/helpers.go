package handlers

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
	"github.com/gofiber/fiber/v2"
)

// ErrorResponse представляет ответ с ошибкой (для обратной совместимости с Swagger)
type ErrorResponse struct {
	Error string `json:"error"`
}

// ParseIDParam парсит ID из параметра пути
func ParseIDParam(c *fiber.Ctx, paramName string) (int, error) {
	idStr := c.Params(paramName)
	if idStr == "" {
		return 0, apperrors.BadRequest(paramName + " обязателен")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, apperrors.BadRequest("неверный формат " + paramName)
	}

	if id <= 0 {
		return 0, apperrors.BadRequest(paramName + " должен быть положительным числом")
	}

	return id, nil
}

// ParseInt64Query парсит int64 из query параметра
func ParseInt64Query(c *fiber.Ctx, paramName string) (int64, error) {
	str := c.Query(paramName)
	if str == "" {
		return 0, apperrors.BadRequest(paramName + " обязателен")
	}

	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, apperrors.BadRequest("неверный формат " + paramName)
	}

	if val <= 0 {
		return 0, apperrors.BadRequest(paramName + " должен быть положительным числом")
	}

	return val, nil
}

// ParseIntQuery парсит int из query параметра
func ParseIntQuery(c *fiber.Ctx, paramName string) (int, error) {
	str := c.Query(paramName)
	if str == "" {
		return 0, apperrors.BadRequest(paramName + " обязателен")
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, apperrors.BadRequest("неверный формат " + paramName)
	}

	if val <= 0 {
		return 0, apperrors.BadRequest(paramName + " должен быть положительным числом")
	}

	return val, nil
}

// ParseOptionalIntQuery парсит опциональный int из query параметра
func ParseOptionalIntQuery(c *fiber.Ctx, paramName string) (*int, error) {
	str := c.Query(paramName)
	if str == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return nil, apperrors.BadRequest("неверный формат " + paramName)
	}

	return &val, nil
}

// ParseFloatQuery парсит float64 из query параметра
func ParseFloatQuery(c *fiber.Ctx, paramName string) (float64, error) {
	str := c.Query(paramName)
	if str == "" {
		return 0, apperrors.BadRequest(paramName + " обязателен")
	}

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, apperrors.BadRequest("неверный формат " + paramName)
	}

	return val, nil
}

// ParseBody парсит тело запроса и валидирует структуру
func ParseBody(c *fiber.Ctx, target interface{}) error {
	if err := c.BodyParser(target); err != nil {
		return apperrors.BadRequest("неверное тело запроса")
	}

	if err := validation.ValidateStruct(target); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		// Конвертируем ValidationErrorResponse в map[string]interface{}
		details := map[string]interface{}{
			"errors": validationErrors.Errors,
		}
		return apperrors.Validation("ошибка валидации данных", details)
	}

	return nil
}

// SuccessResponse представляет успешный ответ
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendSuccess отправляет успешный JSON ответ
func SendSuccess(c *fiber.Ctx, message string) error {
	return c.JSON(SuccessResponse{
		Message: message,
	})
}

// SendSuccessWithData отправляет успешный JSON ответ с данными
func SendSuccessWithData(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// SendCreated отправляет ответ с кодом 201 Created
func SendCreated(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

// SendJSON отправляет JSON ответ
func SendJSON(c *fiber.Ctx, data interface{}) error {
	c.Set("Content-Type", "application/json; charset=utf-8")
	return c.JSON(data)
}
