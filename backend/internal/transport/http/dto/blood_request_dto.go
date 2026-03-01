package dto

import "time"

// ============================================
// Blood Request Status
// ============================================

// BloodRequestStatus представляет статус заявки на поиск крови
type BloodRequestStatus string

const (
	BloodRequestStatusActive BloodRequestStatus = "active"
	BloodRequestStatusClosed BloodRequestStatus = "closed"
	BloodRequestStatusDraft  BloodRequestStatus = "draft"
)

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
	PetID                  string   `json:"petId" doc:"ID питомца" minLength:"1" example:"PET-ABCDEABCDE"`
	BloodVolumeNeeded      int32    `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" minimum:"1" example:"100"`
	Regions                []string `json:"regions" doc:"Список ID регионов" example:"[\"MOSCOW\", \"SPB\"]"`
	SmallPetsNotifyAllowed bool     `json:"smallPetsNotifyAllowed" doc:"Разрешить уведомления для мелких питомцев"`
	Description            string   `json:"description,omitempty" doc:"Дополнительное описание" maxLength:"1000"`
	BloodGroupNames        []string `json:"bloodGroupNames" doc:"Список групп крови" enum:"DEA 1+,DEA 1-,A,B,AB" example:"[\"DEA 1+\", \"A\"]"`
	BloodComponentIDs      []string `json:"bloodComponentIds" doc:"Список ID компонентов крови"`
}

// CreateBloodRequestOutput представляет ответ на создание заявки
type CreateBloodRequestOutput struct {
	Body CreateBloodRequestResult
}

// CreateBloodRequestResult представляет результат создания заявки
type CreateBloodRequestResult struct {
	ID        string             `json:"id" doc:"ID созданной заявки" example:"BLS-ABCDEABCDE"`
	PetID     string             `json:"petId" doc:"ID питомца" example:"PET-ABCDEABCDE"`
	Status    BloodRequestStatus `json:"status" doc:"Статус заявки" enum:"active,closed,draft"`
	CreatedAt *time.Time         `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
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
	BloodVolumeNeeded      *int32   `json:"bloodVolumeNeeded,omitempty" doc:"Необходимый объем крови в мл" minimum:"1"`
	BloodVolumeReserved    *int32   `json:"bloodVolumeReserved,omitempty" doc:"Зарезервированный объем крови в мл" minimum:"0"`
	Regions                []string `json:"regions,omitempty" doc:"Список ID регионов"`
	SmallPetsNotifyAllowed *bool    `json:"smallPetsNotifyAllowed,omitempty" doc:"Разрешить уведомления для мелких питомцев"`
	Description            *string  `json:"description,omitempty" doc:"Дополнительное описание" maxLength:"1000"`
	PhotoURLs              []string `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	BloodGroupNames        []string `json:"bloodGroupNames,omitempty" doc:"Список групп крови"`
	BloodComponentIDs      []string `json:"bloodComponentIds,omitempty" doc:"Список ID компонентов крови"`
	OnBoarding             []string `json:"onBoarding,omitempty" doc:"Список пройденных онбордингов"`
	Status                 *string  `json:"status,omitempty" doc:"Статус заявки"`
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
// Apply For Blood Request (Donor Response)
// ============================================

// ApplyForBloodRequestInput представляет запрос на отклик донора
type ApplyForBloodRequestInput struct {
	BloodRequestIDPath
	Body ApplyForBloodRequestBody
}

// ApplyForBloodRequestBody представляет тело запроса на отклик
type ApplyForBloodRequestBody struct {
	DonorID    string   `json:"donorId" doc:"ID питомца-донора" minLength:"1" example:"PET-ABCDEABCDE"`
	Conditions []string `json:"conditions,omitempty" doc:"Условия донации" enum:"free,paid,food,taxi_compensation"`
}

// ApplyForBloodRequestOutput представляет ответ на отклик
type ApplyForBloodRequestOutput struct {
	Body DonorResponseResult
}

// DonorResponseResult представляет результат создания отклика
type DonorResponseResult struct {
	ID        string              `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	RequestID string              `json:"requestId" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	DonorID   string              `json:"donorId" doc:"ID донора" example:"PET-ABCDEABCDE"`
	Status    DonorResponseStatus `json:"status" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	CreatedAt *time.Time          `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// List Blood Requests
// ============================================

// ListBloodRequestsInput представляет запрос на список заявок
type ListBloodRequestsInput struct {
	Body ListBloodRequestsFilter
}

// ListBloodRequestsFilter представляет фильтр для списка заявок
type ListBloodRequestsFilter struct {
	PetID  string             `json:"petId,omitempty" doc:"ID питомца для фильтрации"`
	Status BloodRequestStatus `json:"status,omitempty" doc:"Статус заявки" enum:"active,closed,draft"`
	Limit  int                `json:"limit,omitempty" doc:"Максимальное количество результатов" minimum:"1" maximum:"100"`
	Offset int                `json:"offset,omitempty" doc:"Смещение для пагинации" minimum:"0"`
}

// ListBloodRequestsOutput представляет ответ со списком заявок
type ListBloodRequestsOutput struct {
	Body BloodRequestsList
}

// BloodRequestsList представляет список заявок
type BloodRequestsList struct {
	Items []BloodRequestDetail `json:"items" doc:"Список заявок"`
	Total int                  `json:"total" doc:"Общее количество заявок"`
}

// ============================================
// Common Types
// ============================================

// BloodRequestDetail представляет полные данные заявки
type BloodRequestDetail struct {
	ID                     string             `json:"id" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	PetID                  string             `json:"petId" doc:"ID питомца" example:"PET-ABCDEABCDE"`
	BloodVolumeNeeded      int32              `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" example:"100"`
	BloodVolumeReserved    int32              `json:"bloodVolumeReserved" doc:"Зарезервированный объем крови в мл" example:"0"`
	Regions                []string           `json:"regions" doc:"Список регионов"`
	SmallPetsNotifyAllowed bool               `json:"smallPetsNotifyAllowed" doc:"Разрешить уведомления для мелких питомцев"`
	Description            string             `json:"description,omitempty" doc:"Дополнительное описание"`
	PhotoURLs              []string           `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	BloodGroupNames        []string           `json:"bloodGroupNames" doc:"Список групп крови"`
	BloodComponentIDs      []string           `json:"bloodComponentIds" doc:"Список ID компонентов крови"`
	OnBoarding             []string           `json:"onBoarding" doc:"Список пройденных онбордингов"`
	Status                 BloodRequestStatus `json:"status" doc:"Статус заявки" enum:"active,closed,draft"`
	Responses              []DonorApplication `json:"responses,omitempty" doc:"Отклики доноров"`
	SuitableDonors         int                `json:"suitableDonors" doc:"Количество подходящих доноров"`
	CreatedAt              *time.Time         `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt              *time.Time         `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt              *time.Time         `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

// DonorApplication представляет отклик донора
type DonorApplication struct {
	ID              string              `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	RequestID       string              `json:"requestId" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	DonorID         string              `json:"donorId" doc:"ID донора" example:"PET-ABCDEABCDE"`
	DonorName       string              `json:"donorName" doc:"Имя донора" example:"Барсик"`
	DonorPhotos     []string            `json:"donorPhotos,omitempty" doc:"Фотографии донора"`
	DonorBloodGroup string              `json:"donorBloodGroup" doc:"Группа крови донора" example:"DEA 1+"`
	Amount          int32               `json:"amount" doc:"Объем крови в мл" example:"450"`
	WarnFactors     []string            `json:"warnFactors,omitempty" doc:"Предупреждающие факторы"`
	Conditions      []string            `json:"conditions" doc:"Условия донации" enum:"free,paid,food,taxi_compensation"`
	Status          DonorResponseStatus `json:"status" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	CreatedAt       *time.Time          `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
	UpdatedAt       *time.Time          `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}

// DonorResponseStatus представляет статус отклика донора
type DonorResponseStatus string

const (
	DonorResponseStatusPending  DonorResponseStatus = "pending"
	DonorResponseStatusAccepted DonorResponseStatus = "accepted"
	DonorResponseStatusDeclined DonorResponseStatus = "declined"
	DonorResponseStatusDonated  DonorResponseStatus = "donated"
)
