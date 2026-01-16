package dto

import (
	"time"
)

// PetHealth представляет информацию о здоровье питомца
type PetHealth struct {
	ReproductiveStatus    *string    `json:"reproductiveStatus,omitempty" doc:"Репродуктивный статус питомца" enum:"pregnancy,lactation,estrus"`
	HealthStatus          *string    `json:"healthStatus,omitempty" doc:"Общее состояние здоровья питомца" enum:"healthy,ill,unknown"`
	LastDonation          *time.Time `json:"lastDonation,omitempty" doc:"Дата последней сдачи крови" example:"2023-10-01T12:00:00Z"`
	Transfused            *bool      `json:"transfused,omitempty" doc:"Были ли переливания крови" example:"false"`
	Medications           *string    `json:"medications,omitempty" doc:"Текущие лекарства" example:"Antibiotics"`
	SurgicalInterventions *string    `json:"surgicalInterventions,omitempty" doc:"Хирургические вмешательства" example:"Sterilization"`
}

// PetTreatment представляет информацию о лечении питомца
type PetTreatment struct {
	RabiesVaccinationDate     *time.Time `json:"rabiesVaccinationDate,omitempty" doc:"Дата вакцинации от бешенства" example:"2023-10-01T12:00:00Z"`
	InfectionVaccinationDate  *time.Time `json:"infectionVaccinationDate,omitempty" doc:"Дата вакцинации от инфекций" example:"2023-10-01T12:00:00Z"`
	EctoparasiteTreatmentDate *time.Time `json:"ectoparasiteTreatmentDate,omitempty" doc:"Дата обработки от эктопаразитов" example:"2023-10-01T12:00:00Z"`
	DewormingDate             *time.Time `json:"dewormingDate,omitempty" doc:"Дата дегельминтизации" example:"2023-10-01T12:00:00Z"`
}

// PetAnalysis представляет информацию об анализах питомца
type PetAnalysis struct {
	ID           *string    `json:"id,omitempty" doc:"ID анализа в системе" example:"PAN-aBcD1aBcD1" readOnly:"true"`
	AnalysisName *string    `json:"analysisName,omitempty" doc:"Название анализа" enum:"leukemia,immunodeficiency,hemoplasmosis,bartonellosis,babesiosis,dirofilaria,ehrlichiosis,anaplasmosis" example:"leukemia"`
	AnalysisType *string    `json:"analysisType,omitempty" doc:"Тип анализа" enum:"PCR,ELISA,ICA,Microscopy,Express" example:"PCR"`
	AnalysisDate *time.Time `json:"analysisDate,omitempty" doc:"Дата проведения анализа" example:"2023-10-01T12:00:00Z"`
}

// PetBonus представляет дополнительную информацию о питомце
type PetBonus struct {
	IsArtist      bool `json:"isArtist,omitempty" doc:"Является ли питомец артистом" example:"false"`
	IsTherapist   bool `json:"isTherapist,omitempty" doc:"Является ли питомец терапевтом" example:"false"`
	IsFormerDonor bool `json:"isFormerDonor,omitempty" doc:"Был ли питомец донором ранее" example:"true"`
	IsGuideDog    bool `json:"isGuideDog,omitempty" doc:"Является ли питомец собакой-проводником" example:"false"`
}

// PetUpdate представляет структуру для обновления существующего питомца
type PetUpdate struct {
	Name            *string        `json:"name,omitempty" validate:"omitempty,min=1,max=100" doc:"Имя питомца" example:"Шарик"`
	ChipNumber      *string        `json:"chipNumber,omitempty" validate:"omitempty,len=15" doc:"Номер чипа" example:"123456789012345"`
	PhotoURL        *string        `json:"photoUrl,omitempty" validate:"omitempty,url,max=255" doc:"URL фотографии питомца" example:"https://example.com/photo.jpg"`
	BreedID         *int           `json:"breedId,omitempty" validate:"omitempty,min=1" doc:"ID породы" example:"1"`
	WeightKg        *float64       `json:"weightKg,omitempty" validate:"omitempty,min=0" doc:"Вес в килограммах" example:"15.5"`
	AgeYears        *int           `json:"ageYears,omitempty" validate:"omitempty,min=0" doc:"Возраст в годах" example:"3"`
	AgeMonths       *int           `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11" doc:"Возраст в месяцах" example:"6"`
	BirthDate       *time.Time     `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	LivingCondition *string        `json:"livingCondition,omitempty" doc:"Условия проживания" enum:"indoor,leashWalking,selfOutdoor" example:"indoor"`
	Gender          *string        `json:"gender,omitempty" doc:"Пол питомца" enum:"male,female" example:"male"`
	Type            *string        `json:"type,omitempty" doc:"Тип животного" enum:"dog,cat" example:"dog"`
	BloodGroup      *string        `json:"bloodGroup,omitempty" doc:"Группа крови" example:"DEA 1.1"`
	PetStatus       *string        `json:"petStatus,omitempty" doc:"Статус питомца" enum:"donor,recipient" example:"donor"`
	Health          *PetHealth     `json:"health,omitempty" doc:"Информация о здоровье"`
	Treatments      *PetTreatment  `json:"treatments,omitempty" doc:"Информация о лечении"`
	Analyses        []*PetAnalysis `json:"analyses,omitempty" doc:"Список анализов"`
	Bonuses         *PetBonus      `json:"bonuses,omitempty" doc:"Дополнительная информация"`
}

// Pet представляет ответ с информацией о питомце
type Pet struct {
	ID              string         `json:"id,omitempty" doc:"Уникальный идентификатор питомца" example:"PET-aBcDeF1234" readOnly:"true"`
	Name            string         `json:"name" doc:"Имя питомца" example:"Шарик"`
	ChipNumber      string         `json:"chipNumber,omitempty" doc:"Номер чипа" example:"123456789012345"`
	PhotoURL        string         `json:"photoUrl,omitempty" doc:"URL фотографии питомца" example:"https://example.com/photo.jpg"`
	BreedID         int            `json:"breedId" doc:"ID породы" example:"1"`
	WeightKg        float64        `json:"weightKg,omitempty" doc:"Вес в килограммах" example:"15.5"`
	AgeYears        int            `json:"ageYears,omitempty" doc:"Возраст в годах" example:"3"`
	AgeMonths       int            `json:"ageMonths,omitempty" doc:"Возраст в месяцах" example:"6"`
	BirthDate       *time.Time     `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	LivingCondition string         `json:"livingCondition,omitempty" doc:"Условия проживания" enum:"indoor,leashWalking,selfOutdoor" example:"indoor"`
	Gender          string         `json:"gender,omitempty" doc:"Пол питомца" enum:"male,female" example:"male"`
	Type            string         `json:"type" doc:"Тип животного" enum:"dog,cat" example:"dog"`
	BloodGroup      string         `json:"bloodGroup,omitempty" doc:"Группа крови" example:"DEA 1.1"`
	PetStatus       string         `json:"petStatus" doc:"Статус питомца" enum:"donor,recipient,none" example:"donor"`
	Health          *PetHealth     `json:"health,omitempty" doc:"Информация о здоровье"`
	Treatments      *PetTreatment  `json:"treatments,omitempty" doc:"Информация о лечении"`
	Analyses        []*PetAnalysis `json:"analyses,omitempty" doc:"Список анализов"`
	Bonuses         *PetBonus      `json:"bonuses,omitempty" doc:"Дополнительная информация"`
	CreatedAt       *time.Time     `json:"createdAt,omitempty" doc:"Дата создания записи" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt       *time.Time     `json:"updatedAt,omitempty" doc:"Дата последнего обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt       *time.Time     `json:"deletedAt,omitempty" doc:"Дата удаления записи" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

// Вспомогательные структуры для Huma

type PetIDPath struct {
	ID string `path:"id" doc:"ID питомца" minLength:"1" example:"PET-aBcDeF1234"`
}

type PetUserIDPath struct {
	ID string `path:"user_id" doc:"ID пользователя" minLength:"1" example:"1"`
}

type PetPreloadQuery struct {
	WithHealth     bool `query:"with_health" doc:"Включить данные о здоровье"`
	WithTreatments bool `query:"with_treatments" doc:"Включить данные о ветеринарных обработках"`
	WithAnalysis   bool `query:"with_analysis" doc:"Включить данные об анализах"`
	WithBonuses    bool `query:"with_bonuses" doc:"Включить данные о бонусах"`
	WithAll        bool `query:"with_all" doc:"Включить все связанные данные"`
}

type AvatarPathParam struct {
	Path string `path:"path" doc:"Путь к аватарке питомца" example:"pets/PET-aBcDeF1234/avatar.jpg"`
}

type PetResponse struct {
	Body Pet
}

type PetsResponse struct {
	Body []Pet
}

type UploadURLResponse struct {
	Body struct {
		URL  string `json:"url"`
		Path string `json:"path"`
	}
}

type ConfirmUploadResponse struct {
	Body struct {
		PublicURL string `json:"publicUrl"`
	}
}
