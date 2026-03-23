package dto

import "time"

// ============================================
// Path Parameters
// ============================================

// BloodRequestIDPath представляет параметр пути с ID заявки
type BloodRequestIDPath struct {
	ID string `path:"id" doc:"ID заявки на поиск крови" minLength:"1" example:"BLS-ABCDEABCDE"`
}

// PetIDPath представляет параметр пути с ID питомца
type PetIDPath struct {
	ID string `path:"id" doc:"ID питомца" minLength:"1" example:"PET-ABCDEABCDE"`
}

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
	BloodVolumeNeeded        int32    `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" minimum:"1" example:"100"`
	Regions                  []string `json:"regions" doc:"Список ID регионов" example:"[\"MOSCOW\", \"SPB\"]"`
	SmallPetsNotifyAllowed   bool     `json:"smallPetsNotifyAllowed" doc:"Разрешить уведомления для мелких питомцев"`
	Description              string   `json:"description,omitempty" doc:"Дополнительное описание" maxLength:"1000"`
	BloodGroupNames          []string `json:"bloodGroupNames" doc:"Список групп крови" enum:"DEA 1+,DEA 1-,A,B,AB" example:"[\"DEA 1+\", \"A\"]"`
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
	Status    string     `json:"status" doc:"Статус заявки" enum:"active,closed,draft"`
	CreatedAt *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Update Blood Request
// ============================================

// UpdateBloodRequestInput представляет запрос на обновление заявки
type UpdateBloodRequestInput struct {
	BloodRequestIDPath
	Body UpdateBloodRequestBody
}

// UpdateBloodRequestBody представляет тело запроса на обновление заявки
type UpdateBloodRequestBody struct {
	BloodVolumeNeeded        *int32   `json:"bloodVolumeNeeded,omitempty" doc:"Необходимый объем крови в мл" minimum:"1"`
	BloodVolumeReserved      *int32   `json:"bloodVolumeReserved,omitempty" doc:"Зарезервированный объем крови в мл" minimum:"0"`
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

// GetBloodRequestByIDInput представляет запрос на получение заявки по ID
type GetBloodRequestByIDInput struct {
	BloodRequestIDPath
}

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

// GetBloodRequestByPetIDInput представляет запрос на получение заявки по ID питомца
type GetBloodRequestByPetIDInput struct {
	PetIDPath
}

// GetBloodRequestByPetIDOutput представляет ответ с данными заявки
type GetBloodRequestByPetIDOutput struct {
	Body BloodRequestDetail
}

// ============================================
// Apply Response
// ============================================

// DonorResponseIDPath представляет параметр пути с ID отклика донора
type DonorResponseIDPath struct {
	ID string `path:"id" doc:"ID отклика донора" minLength:"1" example:"RES-ABCDEABCDE"`
}

// ApplyResponseInput представляет запрос на применение отклика донора
type ApplyResponseInput struct {
	DonorResponseIDPath
}

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

// DeleteBloodRequestInput представляет запрос на удаление заявки
type DeleteBloodRequestInput struct {
	BloodRequestIDPath
}

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
	BloodVolumeNeeded        int32              `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" example:"100"`
	BloodVolumeReserved      int32              `json:"bloodVolumeReserved" doc:"Зарезервированный объем крови в мл" example:"0"`
	Regions                  []string           `json:"regions" doc:"Список регионов"`
	SmallPetsNotifyAllowed   bool               `json:"smallPetsNotifyAllowed" doc:"Разрешить уведомления для мелких питомцев"`
	Description              string             `json:"description,omitempty" doc:"Дополнительное описание"`
	PhotoURLs                []string           `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	BloodGroupNames          []string           `json:"bloodGroupNames" doc:"Список групп крови"`
	BloodComponentIDs        []string           `json:"bloodComponentIds" doc:"Список ID компонентов крови"`
	OnBoarding               []string           `json:"onBoarding" doc:"Список пройденных онбордингов"`
	PrioritySearch           bool               `json:"prioritySearch" doc:"Приоритетный поиск"`
	IncludeUnknownBloodGroup bool               `json:"includeUnknownBloodGroup" doc:"Включить неизвестную группу крови"`
	Status                   string             `json:"status" doc:"Статус заявки" enum:"active,closed,draft"`
	Responses                []DonorApplication `json:"responses,omitempty" doc:"Отклики доноров"`
	SuitableDonors           int                `json:"suitableDonors" doc:"Количество подходящих доноров"`
	CreatedAt                *time.Time         `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt                *time.Time         `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt                *time.Time         `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
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
