package dto

import "time"

// DonorResponseStatus представляет статус ответа донора на запрос на поиск крови
type DonorResponseStatus string

// Возможные статусы ответа донора на запрос
const (
	DonorResponseStatusPending  DonorResponseStatus = "pending"  // Ожидание ответа донора
	DonorResponseStatusAccepted DonorResponseStatus = "accepted" // Донор согласился на донацию
	DonorResponseStatusDeclined DonorResponseStatus = "declined" // Донор отказался
	DonorResponseStatusDonated  DonorResponseStatus = "donated"  // Кровь сдана
)

type DonorApplication struct {
	ID         string              `json:"id" doc:"ID ответа донора" example:"DR-ABCDEABCDE"`
	RequestID  string              `json:"requestId" doc:"ID запроса на кровь" example:"BSR-ABCDEABCDE"`
	DonorID    string              `json:"donorId" doc:"ID донора" example:"DON-ABCDEABCDE"`
	Conditions []string            `json:"conditions" doc:"Условия, при которых донор готов помочь" enum:"free,paid,food,taxi_compensation"`
	Status     DonorResponseStatus `json:"status" doc:"Статус ответа донора" enum:"pending,accepted,declined,donated" example:"pending"`
	CreatedAt  *time.Time          `json:"createdAt,omitempty" doc:"Дата создания ответа" example:"2023-10-01T12:00:00Z"`
	UpdatedAt  *time.Time          `json:"updatedAt,omitempty" doc:"Дата последнего обновления ответа" example:"2023-10-01T12:00:00Z"`
}

type DonorApplicationCreate struct {
	DonorID    string   `json:"donorId" doc:"ID донора" example:"DON-ABCDEABCDE"`
	Conditions []string `json:"conditions" doc:"Условия, при которых донор готов помочь" enum:"free,paid,food,taxi_compensation"`
}

type DonorApplicationResponse struct {
	ID        string              `json:"id" doc:"ID ответа донора" example:"DR-ABCDEABCDE"`
	ReqID     string              `json:"requestId" doc:"ID запроса на кровь" example:"BSR-ABCDEABCDE"`
	DonorID   string              `json:"donorId" doc:"ID донора" example:"DON-ABCDEABCDE"`
	Status    DonorResponseStatus `json:"status" doc:"Статус ответа донора" enum:"pending,accepted,declined,donated" example:"pending"`
	CreatedAt *time.Time          `json:"createdAt,omitempty" doc:"Дата создания ответа" example:"2023-10-01T12:00:00Z"`
}

type DonorApplicationCreateResponse struct {
	Body DonorApplicationResponse
}

// // DonorRequestsResponse представляет обертку для ответа со списком запросов доноров для Huma
// type DonorRequestsResponse struct {
// 	Body []BloodSearchDonorRequest
// }

// // DonorRequestResponse представляет обертку для ответа с одним запросом донора для Huma
// type DonorRequestResponse struct {
// 	Body BloodSearchDonorRequest
// }
