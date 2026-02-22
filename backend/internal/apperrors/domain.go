package apperrors

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
	ErrUserInvalidRole      = BadRequest("неверная роль пользователя")
)

// Pet domain errors
var (
	ErrPetNotFound            = NotFound("питомец не найден")
	ErrPetNameRequired        = BadRequest("имя питомца обязательно")
	ErrInvalidPetType         = BadRequest("неверный тип питомца")
	ErrPetInvalidRole         = BadRequest("неверная роль пользователя")
	ErrInvalidGender          = BadRequest("неверный пол животного")
	ErrInvalidLivingCondition = BadRequest("неверные условия проживания")
	ErrInvalidWeight          = BadRequest("вес должен быть положительным числом")
	ErrInvalidAge             = BadRequest("возраст должен быть положительным числом")
	ErrInvalidAgeMonths       = BadRequest("месяцы должны быть от 0 до 11")
	ErrInvalidPetStatus       = BadRequest("неверный статус питомца")
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

// BloodRequest domain errors
var (
	ErrBloodRequestNotFound      = NotFound("заявка на поиск крови не найдена")
	ErrBloodRequestAlreadyExists = AlreadyExists("заявка на поиск крови уже существует для этого питомца")
	ErrInvalidBloodRequestStatus = BadRequest("неверный статус заявки")
)

// DonorResponse domain errors
var (
	ErrDonorResponseAlreadyExists = AlreadyExists("отклик донора уже существует")
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
