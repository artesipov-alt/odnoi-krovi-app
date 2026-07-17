package dto

import (
	"time"

	commondto "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto/common"
)

// ============================================
// Path Parameters
// ============================================

// ============================================
// Query Parameters
// ============================================

// UserPreloadQuery представляет параметры для предзагрузки связанных данных
type UserPreloadQuery struct {
	WithPets            bool `query:"with_pets" doc:"Включить данные о питомцах"`
	WithDonorPreference bool `query:"with_donor_preference" doc:"Включить данные о предпочтениях донора"`
	WithIdentities      bool `query:"with_identities" doc:"Включить данные об идентификаторах пользователя"`
}

// ContactPreloadQuery представляет параметры для предзагрузки связанных данных
type ContactPreloadQuery struct {
	Provider string `query:"provider" doc:"Месенджер для получения контакта" enum:"telegram_bot,max_bot" minLength:"1"`
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
	FullName     string             `json:"fullName" doc:"Полное имя пользователя" minLength:"1" maxLength:"255" example:"Иван Иванов"`
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
	commondto.UserIDPath
	Body UpdateUserBody
}

// UpdateUserBody представляет тело запроса на обновление пользователя
type UpdateUserBody struct {
	FullName        *string                `json:"fullName,omitempty" doc:"Полное имя" minLength:"2" maxLength:"255"`
	Phone           *string                `json:"phone,omitempty" doc:"Номер телефона" pattern:"^\\+?[1-9]\\d{1,14}$" example:"+79991234567" deprecated:"true"`
	Email           *string                `json:"email,omitempty" doc:"Email адрес" format:"email" example:"user@example.com"`
	AllowGeo        *bool                  `json:"allowGeo,omitempty" doc:"Разрешение использовать геоданные"`
	OnBoarding      *[]string              `json:"onBoarding,omitempty" doc:"Статусы онбординга" enum:"START,FIND_BLOOD,RECIPIENT_LIST"`
	LocationID      *string                `json:"locationId,omitempty" doc:"ID локации"`
	DonorPreference *DonorPreferenceParams `json:"donorPreference,omitempty" doc:"Параметры донора"`
}

// UpdateUserBody представляет тело запроса на обновление пользователя
type ChangeUserPhoneBody struct {
	Phone string `json:"phone" doc:"Номер телефона" pattern:"^\\+?[1-9]\\d{1,14}$" example:"+79991234567"`
}

type ChangeUserPhoneInput struct {
	commondto.UserIDPath
	Body ChangeUserPhoneBody
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
// Get User By ID
// ============================================

// GetUserByIDInput представляет запрос на получение пользователя по ID
type GetUserByIDInput struct {
	commondto.UserIDPath
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
	commondto.TelegramIDPath
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
	commondto.UserIDPath
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
	commondto.UserIDPath
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
	commondto.UserIDPath
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
	TelegramID       int64            `json:"telegramId,omitempty" doc:"ID пользователя в мессенджере" format:"int64" example:"123456789" minimum:"1" deprecated:"true"`
	FullName         string           `json:"fullName" doc:"Полное имя" example:"Иван Иванов"`
	Phone            string           `json:"phone,omitempty" doc:"Телефон" example:"+79991234567"`
	Verified         bool             `json:"verified" doc:"Верифицирован ли пользователь"`
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
	Identities       []Identity       `json:"identities" doc:"Список идентификаторов пользователя в различных системах"`
	CreatedAt        *time.Time       `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt        *time.Time       `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt        *time.Time       `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	LastSeenAt       *time.Time       `json:"lastSeenAt,omitempty" doc:"Время последнего посещения" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

// VerifyUserPhoneInput представляет запрос на верификацию номера телефона
type VerifyUserPhoneInput struct {
	commondto.UserIDPath
	Body VerifyUserPhoneBody
}

// VerifyUserPhoneBody представляет тело запроса на верификацию номера
type VerifyUserPhoneBody struct {
	Code string `json:"code" doc:"Код подтверждения из SMS" minLength:"4" maxLength:"4" example:"2026"`
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
	OpenForContact        bool       `json:"openForContact" doc:"Разрешает реципиентам находить донора как потенциального"`
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
	OpenForContact        *bool    `json:"openForContact,omitempty" doc:"Разрешить реципиентам находить донора как потенциального"`
}

type Identity struct {
	ProviderName string `json:"providerName" doc:"Название провайдера идентификации (например, telegram, max)"`
	ProviderID   string `json:"providerId" doc:"ID пользователя у провайдера" example:"123456789"`
	RefURL       string `json:"refUrl,omitempty" doc:"URL для ссылки на профиль пользователя у провайдера с меткой" example:"https://t.me/username"`
}
