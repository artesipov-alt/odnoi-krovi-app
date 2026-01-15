package dto

import "time"

// UserIDPath представляет параметры пути с ID пользователя
type UserIDPath struct {
	ID string `path:"id" doc:"ID пользователя" minLength:"1" example:"USR-ABCDEABCDE"`
}

// TelegramIDQuery представляет параметры запроса с Telegram ID
type TelegramIDQuery struct {
	TelegramID int64 `query:"telegram_id" doc:"Telegram ID пользователя" minimum:"1" example:"123456789"`
}

// UserRegistrationSimple представляет запрос на простую регистрацию через Telegram
type UserRegistrationSimple struct {
	TelegramID int64  `json:"telegramId" doc:"Telegram ID пользователя" format:"int64" example:"123456789" minimum:"1"`
	FullName   string `json:"fullName,omitempty" doc:"Полное имя пользователя" maxLength:"255" example:"Иван Иванов"`
}

// UserUpdate представляет структуру для обновления данных пользователя
type UserUpdate struct {
	FullName   *string `json:"fullName,omitempty" doc:"Полное имя" minLength:"2" maxLength:"255"`
	Phone      *string `json:"phone,omitempty" doc:"Номер телефона" pattern:"^\\+?[1-9]\\d{1,14}$"`
	Email      *string `json:"email,omitempty" doc:"Email адрес" format:"email"`
	AllowGeo   *bool   `json:"allowGeo,omitempty" doc:"Разрешение использовать геоданные"`
	OnBoarding *bool   `json:"onBoarding,omitempty" doc:"Статус онбординга"`
	LocationID *int    `json:"locationId,omitempty" doc:"ID локации" minimum:"1"`
}

// User представляет данные пользователя для ответа API
type User struct {
	ID               string     `json:"id" doc:"Внутренний ID пользователя" example:"USR-ABCDEABCDE"`
	TelegramID       int64      `json:"telegramId" doc:"Telegram ID" example:"123456789"`
	FullName         string     `json:"fullName" doc:"Полное имя" example:"Иван Иванов"`
	Phone            string     `json:"phone,omitempty" doc:"Телефон" example:"+79991234567"`
	Email            string     `json:"email,omitempty" doc:"Email" example:"user@example.com"`
	OrganizationName string     `json:"organizationName,omitempty" doc:"Название организации"`
	ConsentPd        bool       `json:"consentPd" doc:"Согласие на ПД"`
	OnBoarding       bool       `json:"onBoarding" doc:"Статус онбординга"`
	AllowGeo         bool       `json:"allowGeo" doc:"Разрешение использовать геоданные"`
	LocationID       int        `json:"locationId,omitempty" doc:"ID локации"`
	Role             string     `json:"role" doc:"Роль"`
	Pets             []Pet      `json:"pets,omitempty" doc:"Список питомцев"`
	CreatedAt        *time.Time `json:"createdAt" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
	UpdatedAt        *time.Time `json:"updatedAt" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z"`
}

// UserResponse представляет обертку для ответа с одним пользователем для Huma
type UserResponse struct {
	Body User
}

// UsersDeletedBody представляет тело ответа со списком удаленных пользователей
type UsersDeletedBody struct {
	Message string `json:"message" doc:"Информационное сообщение"`
	Users   []User `json:"users" doc:"Список удаленных пользователей"`
}

// UsersDeletedResponse представляет ответ со списком удаленных пользователей
type UsersDeletedResponse struct {
	Body UsersDeletedBody
}

type UserPreloadQuery struct {
	WithPets bool `query:"with_pets" doc:"Включить данные о питомцах"`
}
