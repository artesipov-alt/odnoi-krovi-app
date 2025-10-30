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
	ErrInvalidGender          = BadRequest("неверный пол животного")
	ErrInvalidLivingCondition = BadRequest("неверные условия проживания")
	ErrInvalidWeight          = BadRequest("вес должен быть положительным числом")
	ErrInvalidAge             = BadRequest("возраст должен быть положительным числом")
	ErrInvalidAgeMonths       = BadRequest("месяцы должны быть от 0 до 11")
)

// VetClinic domain errors
var (
	ErrClinicNotFound           = NotFound("ветеринарная клиника не найдена")
	ErrClinicAlreadyExists      = AlreadyExists("клиника с таким названием уже существует")
	ErrClinicNameRequired       = BadRequest("название клиники обязательно")
	ErrClinicAddressRequired    = BadRequest("адрес клиники обязателен")
	ErrClinicPhoneInvalid       = BadRequest("неверный формат номера телефона клиники")
	ErrClinicCoordinatesInvalid = BadRequest("неверные координаты клиники")
)

// BloodStock domain errors
var (
	ErrBloodStockNotFound     = NotFound("запас крови не найден")
	ErrInvalidBloodType       = BadRequest("неверный тип крови")
	ErrInvalidVolume          = BadRequest("объем должен быть положительным числом")
	ErrInvalidExpirationDate  = BadRequest("неверная дата истечения срока")
	ErrInvalidBloodStatus     = BadRequest("неверный статус запаса крови")
	ErrBloodStockExpired      = BadRequest("запас крови просрочен")
	ErrBloodStockNotAvailable = BadRequest("запас крови недоступен")
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

// Donation domain errors
var (
	ErrDonationNotFound         = NotFound("донация не найдена")
	ErrDonorNotFound            = NotFound("донор не найден")
	ErrRecipientNotFound        = NotFound("реципиент не найден")
	ErrInvalidDonationStatus    = BadRequest("неверный статус донации")
	ErrInvalidDonationDate      = BadRequest("неверная дата донации")
	ErrDonationAlreadyCompleted = BadRequest("донация уже завершена")
)

// BloodSearch domain errors
var (
	ErrBloodSearchNotFound = NotFound("поиск крови не найден")
	ErrInvalidSearchStatus = BadRequest("неверный статус поиска")
	ErrSearchAlreadyClosed = BadRequest("поиск крови уже закрыт")
)

// Helper functions для создания ошибок с контекстом

// NewUserNotFoundError создает ошибку с ID пользователя
func NewUserNotFoundError(userID int) *AppError {
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
func NewPetNotFoundError(petID int) *AppError {
	return NotFound("питомец не найден").WithDetails(map[string]any{
		"pet_id": petID,
	})
}

// NewClinicNotFoundError создает ошибку с ID клиники
func NewClinicNotFoundError(clinicID int) *AppError {
	return NotFound("ветеринарная клиника не найдена").WithDetails(map[string]any{
		"clinic_id": clinicID,
	})
}

// NewBloodStockNotFoundError создает ошибку с ID запаса крови
func NewBloodStockNotFoundError(stockID int) *AppError {
	return NotFound("запас крови не найден").WithDetails(map[string]any{
		"stock_id": stockID,
	})
}

// NewBloodTypeNotFoundError создает ошибку с ID типа крови
func NewBloodTypeNotFoundError(bloodTypeID int) *AppError {
	return NotFound("тип крови не найден").WithDetails(map[string]any{
		"blood_type_id": bloodTypeID,
	})
}
