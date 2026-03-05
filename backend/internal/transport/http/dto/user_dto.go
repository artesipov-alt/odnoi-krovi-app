package dto

import (
	"time"
)

// ============================================
// Path Parameters
// ============================================

// UserIDPath представляет параметр пути с ID пользователя
type UserIDPath struct {
	ID string `path:"id" doc:"ID пользователя" minLength:"1" example:"USR-ABCDEABCDE"`
}

// TelegramIDPath представляет параметр пути с Telegram ID
type TelegramIDPath struct {
	ID int64 `path:"id" doc:"Telegram ID пользователя" minimum:"1" example:"123456789"`
}

// ============================================
// Query Parameters
// ============================================

// UserPreloadQuery представляет параметры для предзагрузки связанных данных
type UserPreloadQuery struct {
	WithPets            bool `query:"with_pets" doc:"Включить данные о питомцах"`
	WithDonorPreference bool `query:"with_donor_preference" doc:"Включить данные о предпочтениях донора"`
}

// ============================================
// Create User
// ============================================

// CreateUserInput представляет запрос на создание пользователя
type CreateUserInput struct {
	Body CreateUserBody
}

// CreateUserBody представляет тело запроса на создание пользователя
type CreateUserBody struct {
	ProviderID   int64              `json:"providerId" doc:"ID пользователя в мессенджере" format:"int64" example:"123456789" minimum:"1"`
	ProviderName string             `json:"providerName" doc:"Название мессенджера" minLength:"1" maxLength:"50" enum:"telegram_bot,max_bot"`
	FullName     string             `json:"fullName" doc:"Полное имя пользователя" minLength:"2" maxLength:"255" example:"Иван Иванов"`
	MetaData     *map[string]string `json:"metaData,omitempty" doc:"Метаданные пользователя"`
}

// CreateUserOutput представляет ответ на создание пользователя
type CreateUserOutput struct {
	Body CreateUserResult
}

// CreateUserResult представляет результат создания пользователя
type CreateUserResult struct {
	ID        string     `json:"id" doc:"ID созданного пользователя" example:"USR-ABCDEABCDE"`
	CreatedAt *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Update User
// ============================================

// UpdateUserInput представляет запрос на обновление пользователя
type UpdateUserInput struct {
	UserIDPath
	Body UpdateUserBody
}

// UpdateUserBody представляет тело запроса на обновление пользователя
type UpdateUserBody struct {
	FullName        *string                `json:"fullName,omitempty" doc:"Полное имя" minLength:"2" maxLength:"255"`
	Phone           *string                `json:"phone,omitempty" doc:"Номер телефона" pattern:"^\\+?[1-9]\\d{1,14}$" example:"+79991234567"`
	Email           *string                `json:"email,omitempty" doc:"Email адрес" format:"email" example:"user@example.com"`
	AllowGeo        *bool                  `json:"allowGeo,omitempty" doc:"Разрешение использовать геоданные"`
	OnBoarding      *[]string              `json:"onBoarding,omitempty" doc:"Статусы онбординга" enum:"START,FIND_BLOOD"`
	LocationID      *string                `json:"locationId,omitempty" doc:"ID локации"`
	DonorPreference *DonorPreferenceParams `json:"donorPreference,omitempty" doc:"Параметры донора"`
}

// UpdateUserOutput представляет ответ на обновление пользователя
type UpdateUserOutput struct {
	Body UpdateUserResult
}

// UpdateUserResult представляет результат обновления пользователя
type UpdateUserResult struct {
	ID        string     `json:"id" doc:"ID обновленного пользователя" example:"USR-ABCDEABCDE"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Auth User
// ============================================

// AuthUserInput представляет запрос на создание пользователя
type AuthUserInput struct {
	Body AuthUserBody
}

// AuthUserBody представляет тело запроса на создание пользователя
type AuthUserBody struct {
	ProviderID   int64  `json:"providerId" doc:"ID пользователя в мессенджере" format:"int64" example:"123456789" minimum:"1"`
	ProviderName string `json:"providerName" doc:"Название мессенджера" minLength:"1" maxLength:"50" enum:"telegram_bot,max_bot"`
	AuthBotToken string `json:"authBotToken,omitempty" doc:"Зашифрованный токен бота для сверки"`
}

// AuthUserOutput представляет ответ на создание пользователя
type AuthUserOutput struct {
	Body AuthUserResult
}

// AuthUserResult представляет результат создания пользователя
type AuthUserResult struct {
	UserID    string     `json:"userId" doc:"ID пользователя на портале" example:"USR-ABCDEABCDE"`
	XBToken   string     `json:"xbToken,omitempty" doc:"JWT токен для аутентификации" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	CreatedAt *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Get User By ID
// ============================================

// GetUserByIDInput представляет запрос на получение пользователя по ID
type GetUserByIDInput struct {
	UserIDPath
	UserPreloadQuery
}

// GetUserByIDOutput представляет ответ с данными пользователя
type GetUserByIDOutput struct {
	Body UserDetail
}

// ============================================
// Get User By Telegram
// ============================================

// GetUserByTelegramInput представляет запрос на получение пользователя по Telegram ID
type GetUserByTelegramInput struct {
	TelegramIDPath
	UserPreloadQuery
}

// GetUserByTelegramOutput представляет ответ с данными пользователя
type GetUserByTelegramOutput struct {
	Body UserDetail
}

// ============================================
// Delete User
// ============================================

// DeleteUserInput представляет запрос на удаление пользователя
type DeleteUserInput struct {
	UserIDPath
}

// DeleteUserOutput представляет ответ на удаление пользователя
type DeleteUserOutput struct {
	Body DeleteUserResult
}

// DeleteUserResult представляет результат удаления пользователя
type DeleteUserResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}

// ============================================
// Reset User
// ============================================

// ResetUserInput представляет запрос на сброс пользователя
type ResetUserInput struct {
	UserIDPath
}

// ResetUserOutput представляет ответ на сброс пользователя
type ResetUserOutput struct {
	Body ResetUserResult
}

// ResetUserResult представляет результат сброса пользователя
type ResetUserResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}

// ============================================
// Restore User
// ============================================

// RestoreUserInput представляет запрос на восстановление пользователя
type RestoreUserInput struct {
	UserIDPath
}

// RestoreUserOutput представляет ответ на восстановление пользователя
type RestoreUserOutput struct {
	Body RestoreUserResult
}

// RestoreUserResult представляет результат восстановления пользователя
type RestoreUserResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}

// ============================================
// Get Deleted Users
// ============================================

// GetDeletedUsersInput представляет запрос на получение удаленных пользователей
type GetDeletedUsersInput struct{}

// GetDeletedUsersOutput представляет ответ со списком удаленных пользователей
type GetDeletedUsersOutput struct {
	Body DeletedUsersList
}

// DeletedUsersList представляет список удаленных пользователей
type DeletedUsersList struct {
	Message string       `json:"message" doc:"Информационное сообщение"`
	Users   []UserDetail `json:"users" doc:"Список удаленных пользователей"`
}

// ============================================
// Common Types
// ============================================

// UserDetail представляет полные данные пользователя
type UserDetail struct {
	ID               string           `json:"id" doc:"Внутренний ID пользователя" example:"USR-ABCDEABCDE" readOnly:"true"`
	TelegramID       int64            `json:"telegramId" doc:"ID пользователя в мессенджере" format:"int64" example:"123456789" minimum:"1" deprecated:"true"`
	FullName         string           `json:"fullName" doc:"Полное имя" example:"Иван Иванов"`
	Phone            string           `json:"phone,omitempty" doc:"Телефон" example:"+79991234567"`
	Email            string           `json:"email,omitempty" doc:"Email" example:"user@example.com"`
	PhotoURLs        []string         `json:"photoUrls,omitempty" doc:"URLs фотографий пользователя"`
	OrganizationName string           `json:"organizationName,omitempty" doc:"Название организации"`
	ConsentPd        bool             `json:"consentPd" doc:"Согласие на обработку персональных данных"`
	OnBoarding       []string         `json:"onBoarding" doc:"Статусы онбординга" enum:"START,FIND_BLOOD"`
	AllowGeo         bool             `json:"allowGeo" doc:"Разрешение использовать геоданные"`
	LocationID       string           `json:"locationId,omitempty" doc:"ID локации"`
	Role             string           `json:"role" doc:"Роль пользователя"`
	Pets             []PetDetail      `json:"pets,omitempty" doc:"Список питомцев"`
	DonorPreference  *DonorPreference `json:"donorPreference,omitempty" doc:"Параметры донора"`
	CreatedAt        *time.Time       `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt        *time.Time       `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt        *time.Time       `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

// SimpleMessage представляет простое текстовое сообщение
type SimpleMessage struct {
	Message string `json:"message" doc:"Сообщение об успехе или ошибке"`
}

// SimpleMessageOutput представляет обертку для простого текстового ответа
type SimpleMessageOutput struct {
	Body SimpleMessage
}

// DonorPreference represents donor's default preferences for blood donation responses
type DonorPreference struct {
	ID                    string     `json:"id" doc:"ID параметров донора" example:"DPR-ABCDEABCDE" readOnly:"true"`
	UserID                string     `json:"userId" doc:"ID пользователя" example:"USR-ABCDEABCDE" readOnly:"true"`
	PreferredLocationIDs  []string   `json:"preferredLocationIds" doc:"Предпочитаемые ID локаций"`
	RecoveryPeriodMonths  int        `json:"recoveryPeriodMonths" doc:"Период восстановления в месяцах" minimum:"2" example:"3"`
	CompensationType      string     `json:"compensationType" doc:"Тип компенсации" enum:"free,paid,food"`
	TaxiCompensation      bool       `json:"taxiCompensation" doc:"Компенсация такси"`
	NotificationFrequency string     `json:"notificationFrequency" doc:"Частота уведомлений" enum:"immediately,daily,weekly,never"`
	CreatedAt             *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt             *time.Time `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt             *time.Time `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

// DonorPreferenceParams holds the parameters for creating or updating a DonorPreference.
type DonorPreferenceParams struct {
	PreferredLocationIDs  []string `json:"preferredLocationIds,omitempty" doc:"Предпочитаемые ID локаций"`
	RecoveryPeriodMonths  *int     `json:"recoveryPeriodMonths,omitempty" doc:"Период восстановления в месяцах" minimum:"2" example:"3"`
	CompensationType      *string  `json:"compensationType,omitempty" doc:"Тип компенсации" enum:"free,paid,food"`
	TaxiCompensation      *bool    `json:"taxiCompensation,omitempty" doc:"Компенсация такси"`
	NotificationFrequency *string  `json:"notificationFrequency,omitempty" doc:"Частота уведомлений" enum:"immediately,daily,weekly,never"`
}
