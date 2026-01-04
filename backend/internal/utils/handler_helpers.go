package utils

import (
	"strconv"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/validation"
	"github.com/labstack/echo/v4"
)

// ErrorResponse представляет ответ с ошибкой (для обратной совместимости с Swagger)
type ErrorResponse struct {
	Message string `json:"message"`
}

// ParseIDParam парсит ID из параметра пути
func ParseIDParam(c echo.Context, paramName string) (int, error) {
	idStr := c.Param(paramName)
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

// ParseStringParam парсит строку из параметра пути
func ParseStringParam(c echo.Context, paramName string) (string, error) {
	str := c.Param(paramName)
	if str == "" {
		return "", apperrors.BadRequest(paramName + " обязателен")
	}

	return str, nil
}

// ParseInt64Query парсит int64 из query параметра
func ParseInt64Query(c echo.Context, paramName string) (int64, error) {
	str := c.QueryParam(paramName)
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
func ParseIntQuery(c echo.Context, paramName string) (int, error) {
	str := c.QueryParam(paramName)
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
func ParseOptionalIntQuery(c echo.Context, paramName string) (*int, error) {
	str := c.QueryParam(paramName)
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
func ParseFloatQuery(c echo.Context, paramName string) (float64, error) {
	str := c.QueryParam(paramName)
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
func ParseBody(c echo.Context, target any) error {
	if err := c.Bind(target); err != nil {
		return apperrors.BadRequest("неверное тело запроса")
	}

	if err := validation.ValidateStruct(target); err != nil {
		validationErrors := validation.GetValidationErrors(err)
		// Конвертируем ValidationErrorResponse в map[string]interface{}
		details := map[string]any{
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
func SendSuccess(c echo.Context, message string) error {
	return c.JSON(200, SuccessResponse{
		Message: message,
	})
}

// SendSuccessWithData отправляет успешный JSON ответ с данными
func SendSuccessWithData(c echo.Context, message string, data any) error {
	return c.JSON(200, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

// SendCreated отправляет ответ с кодом 201 Created
func SendCreated(c echo.Context, data any) error {
	return c.JSON(201, data)
}

// SendJSON отправляет JSON ответ
func SendJSON(c echo.Context, data any) error {
	return c.JSON(200, data)
}
