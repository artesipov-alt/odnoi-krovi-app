package dto

// =====================================================================
// 						 Основные структуры
// =====================================================================

type DonationCard struct {
	DonorData     PetWithApplication `json:"donorData" doc:"Данные донора"`
	RecipientData RecipientShort     `json:"recipientData" doc:"Данные реципиента"`
}

type PetWithApplication struct {
	PetDetail
	Application CoreApplicationData `json:"application" doc:"Данные отклика"`
}

type CoreApplicationData struct {
	ID               string `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	Amount           int32  `json:"amount" doc:"Объем крови в мл" example:"450"`
	CompensationType string `json:"compensationType" doc:"Условия донации" enum:"free,paid,food"`
	TaxiCompensation bool   `json:"taxiCompensation" doc:"Компенсация такси" example:"true"`
	Status           string `json:"status" doc:"Статус отклика" enum:"pending,accepted,rejected,cancelled,completed,failed"`
}

type RecipientShort struct {
	PetName             string   `json:"petName" doc:"Имя питомца" example:"Шарик"`
	PetType             string   `json:"petType" doc:"Тип питомца" enum:"dog,cat"`
	BloodVolumeNeeded   int32    `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" example:"200"`
	BloodVolumeReserved int32    `json:"bloodVolumeReserved" doc:"Зарезервированный объем крови в мл" example:"50"`
	BloodVolumeDonated  int32    `json:"bloodVolumeDonated" doc:"Фактически проведённый объем донации крови в мл" example:"50"`
	PhotoURLs           []string `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
}

// =====================================================================
// 						 Тела для запросов в HUMA
// =====================================================================

type GetDonationOutput struct {
	Body DonationCard
}
