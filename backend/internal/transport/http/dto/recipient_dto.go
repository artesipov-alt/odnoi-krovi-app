package dto

// DonorPreloadQuery представляет параметры для предзагрузки связанных данных
type DonorPreloadQuery struct {
	Status BloodRequestStatus `query:"status,omitempty" doc:"Статус заявки" enum:"active,closed,draft"`
	Limit  int                `query:"limit,omitempty" doc:"Максимальное количество результатов" minimum:"1" maximum:"100"`
	Offset int                `query:"offset,omitempty" doc:"Смещение для пагинации" minimum:"0"`
}

// GetRecipientInput представляет запрос на получение реципиента по ID
type GetRecipientsListInput struct {
	UserIDPath
	DonorPreloadQuery
}

// ListRecipientsOutput представляет ответ со списком реципиентов
type ListRecipientsOutput struct {
	Body RecipientsList
}

// RecipientsList представляет список реципиентов
type RecipientsList struct {
	Items []RecipientDetail `json:"items" doc:"Список реципиентов"`
	Total int               `json:"total" doc:"Общее количество реципиентов"`
}

// RecipientDetail представляет полные данные заявки
type RecipientDetail struct {
	ID                   string             `json:"id" doc:"ID заявки" example:"BLS-ABCDEABCDE"`
	PetID                string             `json:"petId" doc:"ID питомца" example:"PET-ABCDEABCDE"`
	PetName              string             `json:"petName" doc:"Имя питомца" example:"Шарик"`
	PetType              string             `json:"petType" doc:"Тип питомца" enum:"dog,cat"`
	OwnerName            string             `json:"ownerName,omitempty" doc:"Имя владельца" example:"Иван Иванов"`
	SearchRegions        []string           `json:"regions,omitempty" doc:"Список регионов" example:"[\"MSK\", \"MO\"]"`
	BloodVolumeNeeded    int32              `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" example:"200"`
	BloodVolumeReserved  int32              `json:"bloodVolumeReserved" doc:"Зарезервированный объем крови в мл" example:"50"`
	BloodVolumeRemaining int32              `json:"bloodVolumeRemaining,omitempty" doc:"Необходимый остаток объема крови в мл" example:"100"`
	SearchingBloodNames  []string           `json:"searchingBloodNames,omitempty" doc:"Список искомых групп крови" example:"[\"DEA 1+\", \"A\"]"`
	PhotoURLs            []string           `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	BloodGroupName       string             `json:"bloodGroupName" doc:"Группа крови реципиента"`
	PrioritySearch       bool               `json:"prioritySearch,omitempty" doc:"Приоритетный поиск"`
	Status               BloodRequestStatus `json:"status" doc:"Статус заявки" enum:"active,closed,draft"`
	AdvancedInfo         *AdvancedInfo      `json:"advancedInfo,omitempty" doc:"Дополнительная информация"`
	MatchingDonors       []MatchingDonor    `json:"matchingDonors,omitempty" doc:"Список ID подходящих доноров"`
	DefaultDonorPrefs    *DefaultDonorPrefs `json:"defaultPrefs,omitempty" doc:"Настройки донора по умолчанию"`
}

// DefaultPrefsпредставляет предпочтения реципиента по умолчанию
type DefaultDonorPrefs struct {
	CompensationType string   `json:"compensationType" doc:"Тип компенсации" enum:"free,paid,food"`
	Bonuses          []string `json:"bonuses" doc:"Бонусы за донорство"`
	TaxiCompensation bool     `json:"taxiCompensation" doc:"Компенсация такси"`
}

// AdvancedInfo представляет дополнительную информацию
type AdvancedInfo struct {
	Description string   `json:"description,omitempty" doc:"Дополнительное описание"`
	PhotoURLs   []string `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
}

// ListRecipientsOutput представляет ответ со списком реципиентов
type RecipientDetailsOutput struct {
	Body RecipientDetail
}

// MatchingDonor представляет информацию о подходящем доноре
type MatchingDonor struct {
	PetID           string   `json:"petId" doc:"ID питомца донора" example:"PET-ABCDEABCDE"`
	PetName         string   `json:"petName" doc:"Имя питомца донора" example:"Рекс"`
	Amount          int32    `json:"amount" doc:"Количество возможной крови для донорства" example:"1"`
	DonorBloodGroup string   `json:"donorBloodGroup" doc:"Группа крови донора" example:"DEA 1+"`
	PhotoURLs       []string `json:"photoUrls,omitempty" doc:"Список URL фотографий донора"`
}
