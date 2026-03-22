package dto

import "time"

// ============================================
// Path Parameters
// ============================================

// DonorApplicationIDPath представляет параметр пути с ID отклика
type DonorApplicationIDPath struct {
	ID string `path:"id" doc:"ID отклика донора" minLength:"1" example:"RES-ABCDEABCDE"`
}

// DonorApplication представляет отклик донора
type DonorApplication struct {
	ID                string             `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	RequestID         string             `json:"requestId" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	DonorID           string             `json:"donorId" doc:"ID донора" example:"PET-ABCDEABCDE"`
	DonorName         string             `json:"donorName" doc:"Имя донора" example:"Барсик"`
	DonorPhotos       []string           `json:"donorPhotos,omitempty" doc:"Фотографии донора"`
	DonorBloodGroup   string             `json:"donorBloodGroup" doc:"Группа крови донора" example:"DEA 1+"`
	Amount            int32              `json:"amount" doc:"Объем крови в мл" example:"450"`
	DonorRestrictions *DonorRestrictions `json:"donorRestrictions,omitempty" doc:"Стоп-факторы и предупреждения"`
	CompensationType  string             `json:"compensationType" doc:"Условия донации" enum:"free,paid,food"`
	TaxiCompensation  bool               `json:"taxiCompensation" doc:"Компенсация такси" example:"true"`
	Status            string             `json:"status" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	CreatedAt         *time.Time         `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
	UpdatedAt         *time.Time         `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Create Donor Response
// ============================================

// CreateDonorApplicationInput представляет запрос на создание отклика донора
type CreateDonorApplicationInput struct {
	Body CreateDonorApplicationBody
}

// CreateDonorApplicationBody представляет тело запроса на создание отклика
type CreateDonorApplicationBody struct {
	RequestID        string `json:"requestId" doc:"ID заявки на поиск крови" minLength:"1" example:"BLS-ABCDEABCDE"`
	DonorID          string `json:"donorId" doc:"ID питомца-донора" minLength:"1" example:"PET-ABCDEABCDE"`
	CompensationType string `json:"compensationType,omitempty" doc:"Условия донации" enum:"free,paid,food"`
	TaxiCompensation bool   `json:"taxiCompensation,omitempty" doc:"Компенсация такси" example:"true"`
	Amount           int32  `json:"amount,omitempty" doc:"Объем крови в мл" minimum:"1" maximum:"500" example:"450"`
}

// CreateDonorApplicationOutput представляет ответ на создание отклика
type CreateDonorApplicationOutput struct {
	Body CreateDonorApplicationResult
}

// CreateDonorApplicationResult представляет результат создания отклика
type CreateDonorApplicationResult struct {
	ID        string     `json:"id" doc:"ID созданного отклика" example:"RES-ABCDEABCDE"`
	RequestID string     `json:"requestId" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	DonorID   string     `json:"donorId" doc:"ID донора" example:"PET-ABCDEABCDE"`
	Status    string     `json:"status" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	CreatedAt *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Update Donor Response Status
// ============================================

// UpdateDonorApplicationInput представляет запрос на обновление отклика
type UpdateDonorApplicationInput struct {
	DonorApplicationIDPath
	Body UpdateDonorApplicationBody
}

// UpdateDonorApplicationBody представляет тело запроса на обновление отклика
type UpdateDonorApplicationBody struct {
	Status string `json:"status" doc:"Новый статус отклика" enum:"pending,accepted,declined,donated"`
}

// UpdateDonorApplicationOutput представляет ответ на обновление отклика
type UpdateDonorApplicationOutput struct {
	Body UpdateDonorApplicationResult
}

// UpdateDonorApplicationResult представляет результат обновления отклика
type UpdateDonorApplicationResult struct {
	ID        string     `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	Status    string     `json:"status" doc:"Статус отклика"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// Get Donor Response By ID
// ============================================

// GetDonorApplicationByIDInput представляет запрос на получение отклика по ID
type GetDonorApplicationByIDInput struct {
	DonorApplicationIDPath
}

// GetDonorApplicationByIDOutput представляет ответ с данными отклика
type GetDonorApplicationByIDOutput struct {
	Body DonorApplication
}

// ============================================
// List Donor Responses
// ============================================

// ListDonorApplicationsInput представляет запрос на список откликов
type ListDonorApplicationsInput struct {
	Body ListDonorApplicationsFilter
}

// ListDonorApplicationsFilter представляет фильтр для списка откликов
type ListDonorApplicationsFilter struct {
	RequestID string `json:"requestId,omitempty" doc:"ID заявки для фильтрации"`
	DonorID   string `json:"donorId,omitempty" doc:"ID донора для фильтрации"`
	Status    string `json:"status,omitempty" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	Limit     int    `json:"limit,omitempty" doc:"Максимальное количество результатов" minimum:"1" maximum:"100"`
	Offset    int    `json:"offset,omitempty" doc:"Смещение для пагинации" minimum:"0"`
}

// ListDonorApplicationsOutput представляет ответ со списком откликов
type ListDonorApplicationsOutput struct {
	Body DonorApplicationsList
}

// DonorApplicationsList представляет список откликов
type DonorApplicationsList struct {
	Items []DonorApplication `json:"items" doc:"Список откликов"`
	Total int                `json:"total" doc:"Общее количество откликов"`
}

// ============================================
// Delete Donor Response
// ============================================

// DeleteDonorApplicationInput представляет запрос на удаление отклика
type DeleteDonorApplicationInput struct {
	DonorApplicationIDPath
}

// DeleteDonorApplicationOutput представляет ответ на удаление отклика
type DeleteDonorApplicationOutput struct {
	Body DeleteDonorApplicationResult
}

// DeleteDonorApplicationResult представляет результат удаления отклика
type DeleteDonorApplicationResult struct {
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
	DonorID          string `json:"donorId" doc:"ID питомца-донора" minLength:"1" example:"PET-ABCDEABCDE"`
	CompensationType string `json:"compensationType,omitempty" doc:"Условия донации" enum:"free,paid,food"`
	TaxiCompensation bool   `json:"taxiCompensation,omitempty" doc:"Компенсация такси" example:"true"`
}

// ApplyForBloodRequestOutput представляет ответ на отклик
type ApplyForBloodRequestOutput struct {
	Body DonorApplicationResult
}

// DonorApplicationResult представляет результат создания отклика
type DonorApplicationResult struct {
	ID        string     `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	RequestID string     `json:"requestId" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	DonorID   string     `json:"donorId" doc:"ID донора" example:"PET-ABCDEABCDE"`
	Status    string     `json:"status" doc:"Статус отклика" enum:"pending,accepted,declined,donated"`
	CreatedAt *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}
