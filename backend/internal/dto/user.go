package dto

// UserRegistrationSimple представляет запрос на простую регистрацию через Telegram
type UserRegistrationSimple struct {
	TelegramID int64  `json:"telegramId" doc:"Telegram ID пользователя" format:"int64" example:"123456789" minimum:"1"`
	FullName   string `json:"fullName,omitempty" doc:"Полное имя пользователя" maxLength:"255" example:"Иван Иванов"`
}

// UserRegistration представляет полную структуру регистрации пользователя
type UserRegistrationFull struct {
	FullName   string `json:"fullName" doc:"Полное имя" minLength:"2" maxLength:"255" example:"Иван Иванов"`
	Phone      string `json:"phone" doc:"Номер телефона в формате E.164" pattern:"^\\+?[1-9]\\d{1,14}$" example:"+79991234567"`
	Email      string `json:"email,omitempty" doc:"Email адрес" format:"email" example:"user@example.com"`
	ConsentPD  bool   `json:"consentPd" doc:"Согласие на обработку персональных данных" example:"true"`
	LocationID int    `json:"locationId" doc:"ID локации" minimum:"1" example:"1"`
	Role       string `json:"role" doc:"Роль пользователя" enum:"user,admin" default:"user" example:"user"`
}

// UserUpdate представляет структуру для обновления данных пользователя
type UserUpdate struct {
	FullName   *string `json:"fullName,omitempty" doc:"Полное имя" minLength:"2" maxLength:"255"`
	Phone      *string `json:"phone,omitempty" doc:"Номер телефона" pattern:"^\\+?[1-9]\\d{1,14}$"`
	Email      *string `json:"email,omitempty" doc:"Email адрес" format:"email"`
	AllowGeo   *bool   `json:"allowGeo,omitempty" doc:"Разрешение на использование геопозиции"`
	OnBoarding *bool   `json:"onBoarding,omitempty" doc:"Флаг прохождения онбординга"`
	LocationID *int    `json:"locationId,omitempty" doc:"ID локации" minimum:"1"`
}

// UserResponseDTO представляет данные пользователя для ответа API (без внутренних полей БД)
type UserResponseDTO struct {
	ID               string `json:"id" doc:"Внутренний ID пользователя" example:"01HGWZ..."`
	TelegramID       int64  `json:"telegramId" doc:"Telegram ID" example:"123456789"`
	FullName         string `json:"fullName" doc:"Полное имя" example:"Иван Иванов"`
	Phone            string `json:"phone,omitempty" doc:"Телефон" example:"+79991234567"`
	Email            string `json:"email,omitempty" doc:"Email" example:"user@example.com"`
	OrganizationName string `json:"organizationName,omitempty" doc:"Название организации"`
	ConsentPd        bool   `json:"consentPd" doc:"Согласие на ПД"`
	OnBoarding       bool   `json:"onBoarding" doc:"Статус онбординга"`
	AllowGeo         bool   `json:"allowGeo" doc:"Разрешение гео"`
	LocationID       int    `json:"locationId,omitempty" doc:"ID локации"`
	Role             string `json:"role" doc:"Роль"`
	CreatedAt        string `json:"createdAt" doc:"Дата создания" format:"date-time"`
	UpdatedAt        string `json:"updatedAt" doc:"Дата обновления" format:"date-time"`
}
