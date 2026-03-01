package dto

import "time"

// ============================================
// Path Parameters
// ============================================

// DonorResponseIDPath представляет параметр пути с ID отклика
type DonorResponseIDPath struct {
	ID string `path:"id" doc:"ID отклика донора" minLength:"1" example:"RES-ABCDEABCDE"`
}

// ============================================
// Create Donor Response
// ============================================

// CreateDonorResponseInput представляет запрос на создание отклика донора
type CreateDonorResponseInput struct {
	Body CreateDonorResponseBody
}

// CreateDonorResponseBody представляет тело запроса на создание отклика
type CreateDonorResponseBody struct {
	RequestID  string   `json:"requestId" doc:"ID заявки на поиск крови" minLength:"1" example:"BLS-ABCDEABCDE"`
	DonorID    string   `json:"donorId" doc:"ID питомца-донора" minLength:"1" example:"PET-ABCDEABCDE"`
	Conditions []string `json:"conditions,omitempty" doc:"Условия донации" enum:"free,paid,food,taxi_compensation"`
	Amount     int32    `json:"amount,omitempty" doc:"Объем крови в мл" minimum:"1" maximum:"500" example:"450"`
}

// CreateDonorResponseOutput представляет ответ на создание отклика
type CreateDonorResponseOutput struct {
	Body CreateDonorResponseResult
}

// CreateDonorResponseResult представляет результат создания отклика
type CreateDonorResponseResult struct {
	ID        string              `json:"id" doc:"ID созданного отклика" example:"RES-ABCDEABCDE"`
	RequestID string              `json:"requestId" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	DonorID   string              `json:"donorId" doc:"ID донора" example:"PET-ABCDEABCDE"`
	Status    DonorResponseStatus `json:"status" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	CreatedAt *time.Time          `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Update Donor Response Status
// ============================================

// UpdateDonorResponseInput представляет запрос на обновление отклика
type UpdateDonorResponseInput struct {
	DonorResponseIDPath
	Body UpdateDonorResponseBody
}

// UpdateDonorResponseBody представляет тело запроса на обновление отклика
type UpdateDonorResponseBody struct {
	Status string `json:"status" doc:"Новый статус отклика" enum:"pending,accepted,declined,donated"`
}

// UpdateDonorResponseOutput представляет ответ на обновление отклика
type UpdateDonorResponseOutput struct {
	Body UpdateDonorResponseResult
}

// UpdateDonorResponseResult представляет результат обновления отклика
type UpdateDonorResponseResult struct {
	ID        string              `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	Status    DonorResponseStatus `json:"status" doc:"Статус отклика"`
	UpdatedAt *time.Time          `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Get Donor Response By ID
// ============================================

// GetDonorResponseByIDInput представляет запрос на получение отклика по ID
type GetDonorResponseByIDInput struct {
	DonorResponseIDPath
}

// GetDonorResponseByIDOutput представляет ответ с данными отклика
type GetDonorResponseByIDOutput struct {
	Body DonorResponseDetail
}

// ============================================
// List Donor Responses
// ============================================

// ListDonorResponsesInput представляет запрос на список откликов
type ListDonorResponsesInput struct {
	Body ListDonorResponsesFilter
}

// ListDonorResponsesFilter представляет фильтр для списка откликов
type ListDonorResponsesFilter struct {
	RequestID string `json:"requestId,omitempty" doc:"ID заявки для фильтрации"`
	DonorID   string `json:"donorId,omitempty" doc:"ID донора для фильтрации"`
	Status    string `json:"status,omitempty" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	Limit     int    `json:"limit,omitempty" doc:"Максимальное количество результатов" minimum:"1" maximum:"100"`
	Offset    int    `json:"offset,omitempty" doc:"Смещение для пагинации" minimum:"0"`
}

// ListDonorResponsesOutput представляет ответ со списком откликов
type ListDonorResponsesOutput struct {
	Body DonorResponsesList
}

// DonorResponsesList представляет список откликов
type DonorResponsesList struct {
	Items []DonorResponseDetail `json:"items" doc:"Список откликов"`
	Total int                   `json:"total" doc:"Общее количество откликов"`
}

// ============================================
// Delete Donor Response
// ============================================

// DeleteDonorResponseInput представляет запрос на удаление отклика
type DeleteDonorResponseInput struct {
	DonorResponseIDPath
}

// DeleteDonorResponseOutput представляет ответ на удаление отклика
type DeleteDonorResponseOutput struct {
	Body DeleteDonorResponseResult
}

// DeleteDonorResponseResult представляет результат удаления отклика
type DeleteDonorResponseResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции"`
}

// ============================================
// Common Types
// ============================================

// DonorResponseDetail представляет полные данные отклика донора
type DonorResponseDetail struct {
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
	CreatedAt       *time.Time          `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt       *time.Time          `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}
