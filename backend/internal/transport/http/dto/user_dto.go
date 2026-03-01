package dto

import (
	"time"
)

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
	FullName   *string   `json:"fullName,omitempty" doc:"Полное имя" minLength:"2" maxLength:"255"`
	Phone      *string   `json:"phone,omitempty" doc:"Номер телефона" pattern:"^\\+?[1-9]\\d{1,14}$"`
	Email      *string   `json:"email,omitempty" doc:"Email адрес" format:"email"`
	PhotoURLs  []string  `json:"photoUrls,omitempty" doc:"URLs фотографий пользователя" validate:"omitempty,dive,max=255"`
	AllowGeo   *bool     `json:"allowGeo,omitempty" doc:"Разрешение использовать геоданные"`
	OnBoarding *[]string `json:"onBoarding,omitempty" doc:"Статусы онбординга" enum:"START,FIND_BLOOD"`
	LocationID *string   `json:"locationId,omitempty" doc:"ID локации" minimum:"1"`
}

// User представляет данные пользователя для ответа API
type User struct {
	ID               string     `json:"id" doc:"Внутренний ID пользователя" example:"USR-ABCDEABCDE" readOnly:"true"`
	TelegramID       int64      `json:"telegramId" doc:"Telegram ID" example:"123456789"`
	FullName         string     `json:"fullName" doc:"Полное имя" example:"Иван Иванов"`
	Phone            string     `json:"phone,omitempty" doc:"Телефон" example:"+79991234567"`
	Email            string     `json:"email,omitempty" doc:"Email" example:"user@example.com"`
	PhotoURLs        []string   `json:"photoUrls,omitempty" doc:"URLs фотографий пользователя" example:"https://example.com/photo.jpg"`
	OrganizationName string     `json:"organizationName,omitempty" doc:"Название организации"`
	ConsentPd        bool       `json:"consentPd" doc:"Согласие на ПД"`
	OnBoarding       []string   `json:"onBoarding" doc:"Статусы онбординга" enum:"START,FIND_BLOOD"`
	AllowGeo         bool       `json:"allowGeo" doc:"Разрешение использовать геоданные"`
	LocationID       string     `json:"locationId,omitempty" doc:"ID локации"`
	Role             string     `json:"role" doc:"Роль"`
	Pets             []Pet      `json:"pets,omitempty" doc:"Список питомцев"`
	CreatedAt        *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

// UserResponse представляет обертку для ответа с одним пользователем для Huma
type UserResponse struct {
	Body User
}

// IDPath представляет параметры пути с ID объекта
type UserPathID struct {
	ID string `path:"id" doc:"ID сущности (заявки/питомца/пользователя)" minLength:"1" example:"ENT-ABCDEABCDE"`
}

// IDPath представляет параметры пути с ID объекта
type IDPathInt struct {
	ID int64 `path:"id" doc:"ID Телеграм" minLength:"1" example:"12345678"`
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

// MessageBody представляет тело простого текстового ответа
// type MessageBody struct {
// 	Message string `json:"message" doc:"Сообщение об успехе или ошибке"`
// }

// // MessageResponse представляет простой текстовый ответ
// type MessageResponse struct {
// 	Body MessageBody
// }
