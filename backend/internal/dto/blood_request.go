package dto

// BloodSearchRequestStatus представляет статус запроса на поиск крови
type BloodSearchRequestStatus string

// Возможные статусы запроса на поиск крови
const (
	BloodSearchRequestStatusActive BloodSearchRequestStatus = "active" // Активный запрос
	BloodSearchRequestStatusClosed BloodSearchRequestStatus = "closed" // Закрытый запрос
	BloodSearchRequestStatusDraft  BloodSearchRequestStatus = "draft"  // Черновик запроса
)

// BloodSearchPetRequest представляет запрос на добавление питомца в систему поиска крови
type BloodSearchPetRequest struct {
	PetID                  string   `json:"petId" validate:"required" doc:"ID питомца" example:"PET-aBcDeF1234"`
	BloodVolumeNeeded      int32    `json:"bloodVolumeNeeded" validate:"required,gt=0" doc:"Необходимый объем крови в мл" example:"100"`
	BloodVolumeReserved    int32    `json:"bloodVolumeReserved,omitempty" doc:"Зарезервированный объем крови в мл" example:"0"`
	Regions                []int32  `json:"regions" validate:"required,min=1" doc:"Список ID регионов, где требуется кровь" example:"[1, 2]"`
	SmallPetsNotifyAllowed bool     `json:"smallPetsNotifyAllowed" doc:"Разрешить уведомления для владельцев мелких питомцев" example:"true"`
	Description            string   `json:"description,omitempty" doc:"Дополнительное описание запроса" example:"Срочно нужна кровь для переливания"`
	PhotoUrls              []string `json:"photoUrls,omitempty" doc:"Список URL фотографий питомца" example:"[\"https://example.com/pet_photo1.jpg\", \"https://example.com/pet_photo2.jpg\"]"`
	BloodGroupNames        []string `json:"bloodGroupIds" doc:"Список названий групп крови, которые подходят" example:"[\"DEA 1.1\", \"DEA 1.2\"]"`
	BloodComponentIds      []int    `json:"bloodComponentIds" doc:"Список ID компонентов крови, которые требуются" example:"[1, 2]"`
}

// BloodSearchPetResponse представляет ответ после создания запроса на поиск крови
type BloodSearchPetResponse struct {
	ID     string                   `json:"id" doc:"ID запроса на поиск крови" example:"BSR-ABCDEABCDE" readOnly:"true"`
	PetID  string                   `json:"petId" doc:"ID питомца, для которого создан запрос" example:"PET-aBcDeF1234"`
	Status BloodSearchRequestStatus `json:"status" doc:"Текущий статус запроса" example:"active"`
}

// BloodSearchFilterRequest представляет параметры для фильтрации запросов на поиск крови
type BloodSearchFilterRequest struct {
	PetID  string                   `json:"petId,omitempty" doc:"ID питомца для фильтрации" example:"PET-aBcDeF1234"`
	Status BloodSearchRequestStatus `json:"status,omitempty" doc:"Статус запроса для фильтрации" example:"active"`
	Limit  int                      `json:"limit,omitempty" doc:"Максимальное количество результатов" example:"10"`
	Offset int                      `json:"offset,omitempty" doc:"Смещение для пагинации" example:"0"`
}

// BloodSearchPetsResponse представляет список запросов на поиск крови
type BloodSearchPetsResponse struct {
	Requests []BloodSearchPetRequest `json:"requests" doc:"Список запросов на поиск крови"`
}

// BloodSearchRequestIDPath представляет параметры пути с ID запроса на поиск крови
type BloodSearchRequestIDPath struct {
	ID string `path:"id" doc:"ID запроса на поиск крови" minLength:"1" example:"BSR-ABCDEABCDE"`
}

// BloodSearchRequestResponse представляет обертку для ответа с одним запросом на поиск крови для Huma
type BloodSearchRequestResponse struct {
	Body BloodSearchPetRequest
}

// BloodSearchRequestsResponse представляет обертку для ответа со списком запросов на поиск крови для Huma
type BloodSearchRequestsResponse struct {
	Body []BloodSearchPetRequest
}
