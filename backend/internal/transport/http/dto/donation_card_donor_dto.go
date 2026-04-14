package dto

import "time"

// =====================================================================
// 						 Основные структуры
// =====================================================================

type DonationCardForDonor struct {
	ApplicationData ApplicationShort  `json:"applicationData" doc:"Данные об отклике"`
	RecipientData   RecipientForDonor `json:"recipientData" doc:"Данные реципиента"`
}

type ApplicationShort struct {
	ID               string     `json:"id" doc:"ID отклика" example:"RES-ABCDEABCDE"`
	PetName          string     `json:"petName" doc:"Имя питомца" example:"Шарик"`
	Amount           float64    `json:"amount" doc:"Объем крови в мл" example:"450"`
	PhotoURLs        []string   `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	CompensationType string     `json:"compensationType" doc:"Условия донации" enum:"free,paid,food"`
	TaxiCompensation bool       `json:"taxiCompensation" doc:"Компенсация такси" example:"true"`
	Bonuses          []string   `json:"bonuses" doc:"Бонусы портала"`
	IsConfirmed      bool       `json:"isConfirmed" doc:"Подтверждение отклика от реципиента" example:"false"`
	RejectedReason   string     `json:"rejectedReason,omitempty" doc:"Причина отказа от донации реципиентом"`
	Status           string     `json:"status" doc:"Статус отклика" enum:"pending,accepted,rejected,cancelled,completed,failed"`
	CreatedAt        *time.Time `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

type RecipientForDonor struct {
	ID                  string           `json:"id" doc:"Уникальный идентификатор питомца" example:"PET-aBcDeF1234" readOnly:"true"`
	PetName             string           `json:"petName" doc:"Имя питомца" example:"Шарик"`
	PetType             string           `json:"petType" doc:"Тип питомца" enum:"dog,cat"`
	BloodVolumeNeeded   float64          `json:"bloodVolumeNeeded" doc:"Необходимый объем крови в мл" example:"200"`
	BloodVolumeReserved float64          `json:"bloodVolumeReserved" doc:"Зарезервированный объем крови в мл" example:"50"`
	BloodVolumeDonated  float64          `json:"bloodVolumeDonated" doc:"Фактически проведённый объем донации крови в мл" example:"50"`
	AdvancedInfo        *AdvancedInfoDTO `json:"advancedInfo,omitempty" doc:"Дополнительная информация"`
	PhotoURLs           []string         `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	SearchingBloodNames []string         `json:"searchingBloodNames" doc:"Список групп крови"`
	BloodGroup          string           `json:"bloodGroup" doc:"Группа крови реципиента"`
	Regions             []string         `json:"regions,omitempty" doc:"Список ID регионов"`
	OwnerID             string           `json:"ownerID" doc:"Идентификатор владельца"`
	OwnerName           string           `json:"ownerName" doc:"Имя владельца"`
	Status              string           `json:"status" doc:"Статус заявки" enum:"active,closed,reserved_full,draft"`
	CreatedAt           *time.Time       `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt           *time.Time       `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt           *time.Time       `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

type AdvancedInfoDTO struct {
	PhotoURLs   []string `json:"photoUrls,omitempty" doc:"Список URL фотографий"`
	Description string   `json:"description,omitempty" doc:"Дополнительное описание"`
}

// =====================================================================
// 						 Тела для запросов в HUMA
// =====================================================================

type PlannedDonationsList struct {
	Items []DonationCardForDonor `json:"items" doc:"Список планируемых донаций"`
	Total int                    `json:"total" doc:"Общее количество донаций"`
}

type ListPlannedDonationsOutput struct {
	Body PlannedDonationsList
}

type CompletedDonationsList struct {
	Items          []DonationCardForDonor `json:"items" doc:"Список завершенных донаций"`
	Total          int                    `json:"total" doc:"Общее количество донаций"`
	TotalDonations int                    `json:"totalDonations" doc:"Общее количество завершенных донаций"`
	TotalVolume    float64                `json:"totalVolume" doc:"Общий объем сданной крови в мл"`
}

type ListCompletedDonationsOutput struct {
	Body CompletedDonationsList
}
