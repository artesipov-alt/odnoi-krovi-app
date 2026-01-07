package apperrors

import (
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/enums"
)

// Domain-specific errors для переиспользования

// User domain errors
var (
	ErrUserNotFound         = NotFound("пользователь не найден")
	ErrUserAlreadyExists    = AlreadyExists("пользователь с этим Telegram ID уже существует")
	ErrInvalidTelegramID    = BadRequest("неверный Telegram ID")
	ErrUserPhoneRequired    = BadRequest("номер телефона обязателен")
	ErrUserEmailInvalid     = BadRequest("неверный формат email")
	ErrUserConsentRequired  = BadRequest("требуется согласие на обработку персональных данных")
	ErrUserLocationRequired = BadRequest("местоположение обязательно")
	ErrUserInvalidRole      = BadRequest(fmt.Sprintf("неверная роль пользователя. Доступные роли: %v", enums.GetAllEntUserRoles()))
)

// Pet domain errors
var (
	ErrPetNotFound            = NotFound("питомец не найден")
	ErrPetNameRequired        = BadRequest("имя питомца обязательно")
	ErrInvalidPetType         = BadRequest(fmt.Sprintf("неверный тип питомца. Доступные типы: %v", enums.GetAllEntPetTypes()))
	ErrPetInvalidStatus       = BadRequest(fmt.Sprintf("неверный статус питомца. Доступные статусы: %v", enums.GetAllEntPetStatuses()))
	ErrPetInvalidRole         = BadRequest("неверная роль питомца")
	ErrInvalidGender          = BadRequest(fmt.Sprintf("неверный пол животного. Доступные значения: %v", enums.GetAllEntGenders()))
	ErrInvalidLivingCondition = BadRequest(fmt.Sprintf("неверные условия проживания. Доступные значения: %v", enums.GetAllEntLivingConditions()))
	ErrInvalidWeight          = BadRequest("вес должен быть положительным числом")
	ErrInvalidAge             = BadRequest("возраст должен быть положительным числом")
	ErrInvalidAgeMonths       = BadRequest("месяцы должны быть от 0 до 11")
)

// Pet Health & Analysis errors
var (
	ErrInvalidHealthStatus       = BadRequest(fmt.Sprintf("неверный статус здоровья. Доступные статусы: %v", enums.GetAllEntHealthStatuses()))
	ErrInvalidReproductiveStatus = BadRequest("неверный репродуктивный статус")
	ErrInvalidAnalysisType       = BadRequest("неверный метод проведения анализа")
)

// BloodType domain errors
var (
	ErrBloodTypeNotFound     = NotFound("тип крови не найден")
	ErrBloodTypeNameRequired = BadRequest("название типа крови обязательно")
)

// Location domain errors
var (
	ErrLocationNotFound     = NotFound("местоположение не найдено")
	ErrLocationNameRequired = BadRequest("название местоположения обязательно")
	ErrInvalidCoordinates   = BadRequest("неверные координаты")
)

// ==========Helper functions для создания ошибок с контекстом=============

// NewUserNotFoundError создает ошибку с ID пользователя
func NewUserNotFoundError(userID string) *AppError {
	return NotFound("пользователь не найден").WithDetails(map[string]any{
		"user_id": userID,
	})
}

// NewUserAlreadyExistsError создает ошибку с Telegram ID
func NewUserAlreadyExistsError(telegramID int64) *AppError {
	return AlreadyExists("пользователь с этим Telegram ID уже существует").WithDetails(map[string]any{
		"telegram_id": telegramID,
	})
}

// NewPetNotFoundError создает ошибку с ID питомца
func NewPetNotFoundError(petID string) *AppError {
	return NotFound("питомец не найден").WithDetails(map[string]any{
		"pet_id": petID,
	})
}

// NewBloodTypeNotFoundError создает ошибку с ID типа крови
func NewBloodTypeNotFoundError(bloodTypeID int) *AppError {
	return NotFound("тип крови не найден").WithDetails(map[string]any{
		"blood_type_id": bloodTypeID,
	})
}
