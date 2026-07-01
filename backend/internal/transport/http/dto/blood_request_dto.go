package dto

import (
	"time"

	commondto "github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto/common"
)

// ============================================
// Path Parameters
// ============================================

// ============================================
// Create Blood Request
// ============================================

// CreateBloodRequestInput представляет запрос на создание заявки
type CreateBloodRequestInput struct {
	Body CreateBloodRequestBody
}

// CreateBloodRequestBody представляет тело запроса на создание заявки
type CreateBloodRequestBody struct {
	PetID                    string   `json:"petId" doc:"ID питомца" minLength:"1" example:"PET-ABCDEABCDE"`
	BloodVolumeNeeded        float64  `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" minimum:"1" example:"100"`
	Regions                  []string `json:"regions" doc:"Список ID регионов" example:"[\"MSK\", \"MO\"]"`
	SmallPetsNotifyAllowed   bool     `json:"smallPetsNotifyAllowed" doc:"Разрешить уведомления для мелких питомцев"`
	Description              string   `json:"description,omitempty" doc:"Дополнительное описание" maxLength:"1000"`
	BloodGroupNames          []string `json:"bloodGroupNames" doc:"Список групп крови" enum:"DEA 1+,DEA 1-,A,B,AB"`
	BloodComponentIDs        []string `json:"bloodComponentIds" doc:"Список ID компонентов крови"`
	PrioritySearch           bool     `json:"prioritySearch" doc:"Приоритетный поиск"`
	IncludeUnknownBloodGroup bool     `json:"includeUnknownBloodGroup" doc:"Включить неизвестную группу крови"`
}

// CreateBloodRequestOutput представляет ответ на создание заявки
type CreateBloodRequestOutput struct {
	Body CreateBloodRequestResult
}

// CreateBloodRequestResult представляет результат создания заявки
type CreateBloodRequestResult struct {
	ID        string     `json:"id" doc:"ID созданной заявки" example:"BLS-ABCDEABCDE"`
	PetID     string     `json:"petId" doc:"ID питомца" example:"PET-ABCDEABCDE"`
	Status    string     `json:"status" doc:"Статус заявки" enum:"active,closed,reserved_full,draft"`
	CreatedAt *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Update Blood Request
// ============================================

// UpdateBloodRequestInput представляет запрос на обновление заявки
type UpdateBloodRequestInput struct {
	commondto.BloodRequestIDPath
	Body UpdateBloodRequestBody
}

// UpdateBloodRequestBody представляет тело запроса на обновление заявки
type UpdateBloodRequestBody struct {
	BloodVolumeNeeded        *float64 `json:"bloodVolumeNeeded,omitempty" doc:"Необходимый объем крови в мл" minimum:"1"`
	BloodVolumeReserved      *float64 `json:"bloodVolumeReserved,omitempty" doc:"Зарезервированный объем крови в мл" minimum:"0"`
	Regions                  []string `json:"regions,omitempty" doc:"Список ID регионов"`
	SmallPetsNotifyAllowed   *bool    `json:"smallPetsNotifyAllowed,omitempty" doc:"Разрешить уведомления для мелких питомцев"`
	Description              *string  `json:"description,omitempty" doc:"Дополнительное описание" maxLength:"1000"`
	BloodGroupNames          []string `json:"bloodGroupNames,omitempty" doc:"Список групп крови"`
	BloodComponentIDs        []string `json:"bloodComponentIds,omitempty" doc:"Список ID компонентов крови"`
	OnBoarding               []string `json:"onBoarding,omitempty" doc:"Список пройденных онбордингов"`
	Status                   *string  `json:"status,omitempty" doc:"Статус заявки"`
	PrioritySearch           *bool    `json:"prioritySearch,omitempty" doc:"Приоритетный поиск"`
	IncludeUnknownBloodGroup *bool    `json:"includeUnknownBloodGroup,omitempty" doc:"Включить неизвестную группу крови"`
}

// UpdateBloodRequestOutput представляет ответ на обновление заявки
type UpdateBloodRequestOutput struct {
	Body UpdateBloodRequestResult
}

// UpdateBloodRequestResult представляет результат обновления заявки
type UpdateBloodRequestResult struct {
	ID        string     `json:"id" doc:"ID обновленной заявки" example:"BLS-ABCDEABCDE"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Get Blood Request By ID
// ============================================

// GetBloodRequestByIDOutput представляет ответ с данными заявки
type GetBloodRequestByIDOutput struct {
	Body BloodRequestDetail
}

// GetDonorByIDOutput представляет ответ с данными заявки
type GetDonorByIDOutput struct {
	Body DonorDetail
}

// ============================================
// Get Blood Request By Pet ID
// ============================================

// GetBloodRequestByPetIDOutput представляет ответ с данными заявки
type GetBloodRequestByPetIDOutput struct {
	Body BloodRequestDetail
}

// ============================================
// Apply Response
// ============================================

// ApplyResponseOutput представляет ответ на применение отклика донора
type ApplyResponseOutput struct {
	Body ApplyResponseResult
}

// ApplyResponseResult представляет результат применения отклика донора
type ApplyResponseResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}

// ============================================
// Delete Blood Request
// ============================================

// DeleteBloodRequestOutput представляет ответ на удаление заявки
type DeleteBloodRequestOutput struct {
	Body DeleteBloodRequestResult
}

// DeleteBloodRequestResult представляет результат удаления заявки
type DeleteBloodRequestResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}

// ============================================
// Common Types
// ============================================

// BloodRequestDetail представляет полные данные заявки
type BloodRequestDetail struct {
	ID                       string             `json:"id" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	PetID                    string             `json:"petId" doc:"ID питомца" example:"PET-ABCDEABCDE"`
	BloodVolumeNeeded        float64            `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" example:"100"`
	BloodVolumeReserved      float64            `json:"bloodVolumeReserved" doc:"Зарезервированный объем крови в мл" example:"0"`
	BloodVolumeDonated       float64            `json:"bloodVolumeDonated" doc:"Фактически проведённый объем донации крови в мл" example:"50"`
	Regions                  []string           `json:"regions" doc:"Список регионов"`
	SmallPetsNotifyAllowed   bool               `json:"smallPetsNotifyAllowed" doc:"Разрешить уведомления для мелких питомцев"`
	Description              string             `json:"description,omitempty" doc:"Дополнительное описание"`
	PhotoURLs                []string           `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	BloodGroupNames          []string           `json:"bloodGroupNames" doc:"Список групп крови"`
	BloodComponentIDs        []string           `json:"bloodComponentIds" doc:"Список ID компонентов крови"`
	OnBoarding               []string           `json:"onBoarding" doc:"Список пройденных онбордингов"`
	PrioritySearch           bool               `json:"prioritySearch" doc:"Приоритетный поиск"`
	IncludeUnknownBloodGroup bool               `json:"includeUnknownBloodGroup" doc:"Включить неизвестную группу крови"`
	Status                   string             `json:"status" doc:"Статус заявки" enum:"active,closed,reserved_full,draft"`
	Responses                []DonorApplication `json:"responses,omitempty" doc:"Отклики доноров"`
	AcceptedDonors           []DonorApplication `json:"acceptedDonors,omitempty" doc:"Принятые отклики доноров"`
	CompletedDonations       []DonorApplication `json:"completedDonations,omitempty" doc:"Завершенные донации"`
	SuitableDonors           int                `json:"suitableDonors" doc:"Количество подходящих доноров"`
	CreatedAt                *time.Time         `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt                *time.Time         `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt                *time.Time         `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

type ConfirmDonorApplicationInput struct {
	commondto.DonorApplicationIDPath
	Body ConfirmData
}

type ConfirmData struct {
	Amount float64 `json:"amount,omitempty" doc:"Фактический объем донации в мл" minimum:"1"`
}

type RejectData struct {
	Reason string `json:"reason,omitempty" doc:"Причина отклонения"`
}

type RejectDonorApplicationInput struct {
	commondto.DonorApplicationIDPath
	Body RejectData
}

// DonorDetail представляет полные данные донора
type DonorDetail struct {
	ResponseID string `json:"responseId" doc:"ID отклика донора" example:"RES-ABCDEABCDE"`
	PetDetail
	Compensation
}

type Compensation struct {
	CompensationType string `json:"compensationType" doc:"Тип компенсации" enum:"free,paid,food"`
	Taxi             bool   `json:"taxi" doc:"Компенсация такси"`
}
