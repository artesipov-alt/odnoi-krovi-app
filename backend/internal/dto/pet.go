package dto

import (
	"time"
)

// PetHealth представляет информацию о здоровье питомца
type PetHealth struct {
	ReproductiveStatus    *string    `json:"reproductiveStatus,omitempty" doc:"Репродуктивный статус питомца (беременность, лактация, течка или отсутствие)" enum:"беременность,лактация,течка,отсутствие"`
	HealthStatus          *string    `json:"healthStatus,omitempty" doc:"Общее состояние здоровья питомца (здоров, болен или неизвестно)" enum:"здоров,болен,неизвестно"`
	LastDonation          *time.Time `json:"lastDonation,omitempty" doc:"Дата последней сдачи крови" example:"2023-10-01T12:00:00Z"`
	Transfused            *bool      `json:"transfused,omitempty" doc:"Были ли переливания крови" example:"false"`
	Medications           *string    `json:"medications,omitempty" doc:"Текущие лекарства" example:"Антибиотики"`
	SurgicalInterventions *string    `json:"surgicalInterventions,omitempty" doc:"Хирургические вмешательства" example:"Стерилизация"`
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
	LeukemiaDate         *time.Time `json:"leukemiaDate,omitempty" doc:"Дата анализа на лейкемию" example:"2023-10-01T12:00:00Z"`
	LeukemiaType         *string    `json:"leukemiaType,omitempty" doc:"Тип лейкемии" example:"Вирусная"`
	ImmunodeficiencyDate *time.Time `json:"immunodeficiencyDate,omitempty" doc:"Дата анализа на иммунодефицит" example:"2023-10-01T12:00:00Z"`
	ImmunodeficiencyType *string    `json:"immunodeficiencyType,omitempty" doc:"Тип иммунодефицита" example:"FIV"`
	HemoplasmosisDate    *time.Time `json:"hemoplasmosisDate,omitempty" doc:"Дата анализа на гемоплазмоз" example:"2023-10-01T12:00:00Z"`
	HemoplasmosisType    *string    `json:"hemoplasmosisType,omitempty" doc:"Тип гемоплазмоза" example:"Mycoplasma haemofelis"`
	BartonellosisDate    *time.Time `json:"bartonellosisDate,omitempty" doc:"Дата анализа на бартонеллез" example:"2023-10-01T12:00:00Z"`
	BartonellosisType    *string    `json:"bartonellosisType,omitempty" doc:"Тип бартонеллеза" example:"Bartonella henselae"`
	BabesiosisDate       *time.Time `json:"babesiosisDate,omitempty" doc:"Дата анализа на бабезиоз" example:"2023-10-01T12:00:00Z"`
	BabesiosisType       *string    `json:"babesiosisType,omitempty" doc:"Тип бабезиоза" example:"Babesia canis"`
	DirofilariaDate      *time.Time `json:"dirofilariaDate,omitempty" doc:"Дата анализа на дирофиляриоз" example:"2023-10-01T12:00:00Z"`
	DirofilariaType      *string    `json:"dirofilariaType,omitempty" doc:"Тип дирофиляриоза" example:"Dirofilaria immitis"`
	EhrlichiosisDate     *time.Time `json:"ehrlichiosisDate,omitempty" doc:"Дата анализа на эрлихиоз" example:"2023-10-01T12:00:00Z"`
	EhrlichiosisType     *string    `json:"ehrlichiosisType,omitempty" doc:"Тип эрлихиоза" example:"Ehrlichia canis"`
	AnaplasmosisDate     *time.Time `json:"anaplasmosisDate,omitempty" doc:"Дата анализа на анаплазмоз" example:"2023-10-01T12:00:00Z"`
	AnaplasmosisType     *string    `json:"anaplasmosisType,omitempty" doc:"Тип анаплазмоза" example:"Anaplasma phagocytophilum"`
}

// PetBonus представляет дополнительную информацию о питомце
type PetBonus struct {
	IsArtist      bool `json:"isArtist,omitempty" doc:"Является ли питомец артистом" example:"false"`
	IsTherapist   bool `json:"isTherapist,omitempty" doc:"Является ли питомец терапевтом" example:"false"`
	IsFormerDonor bool `json:"isFormerDonor,omitempty" doc:"Был ли питомец донором ранее" example:"true"`
	IsGuideDog    bool `json:"isGuideDog,omitempty" doc:"Является ли питомец собакой-проводником" example:"false"`
}

// PetCreate представляет структуру для создания нового питомца
type PetCreate struct {
	Name            string         `json:"name" validate:"required,min=1,max=100" doc:"Имя питомца" example:"Шарик"`
	ChipNumber      string         `json:"chipNumber,omitempty" validate:"omitempty,len=15" doc:"Номер чипа (15 символов)" example:"123456789012345"`
	PhotoURL        string         `json:"photoUrl,omitempty" validate:"omitempty,url,max=255" doc:"URL фотографии питомца" example:"https://example.com/photo.jpg"`
	BreedID         int            `json:"breedId,omitempty" validate:"omitempty,min=1" doc:"ID породы" example:"1"`
	WeightKg        float64        `json:"weightKg,omitempty" validate:"omitempty,min=0" doc:"Вес в килограммах" example:"15.5"`
	AgeYears        int            `json:"ageYears,omitempty" validate:"omitempty,min=0" doc:"Возраст в годах" example:"3"`
	AgeMonths       int            `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11" doc:"Возраст в месяцах (0-11)" example:"6"`
	BirthDate       *time.Time     `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	LivingCondition string         `json:"livingCondition,omitempty" doc:"Условия проживания (домашний, выгул на шлейке, самовыгул)" enum:"indoor,leashed,free" example:"indoor"`
	Gender          string         `json:"gender,omitempty" doc:"Пол питомца (самец, самка)" enum:"male,female" example:"male"`
	Type            string         `json:"type" validate:"required" doc:"Тип животного (собака, кошка)" enum:"dog,cat" example:"dog"`
	BloodGroup      string         `json:"bloodGroup,omitempty" doc:"Группа крови" example:"DEA 1.1"`
	PetStatus       string         `json:"petStatus" validate:"required,oneof=donor recipient unknown" doc:"Статус питомца (донор, реципиент, неизвестно)" enum:"donor,recipient,unknown" example:"donor"`
	Health          *PetHealth     `json:"health,omitempty" doc:"Информация о здоровье"`
	Treatments      *PetTreatment  `json:"treatments,omitempty" doc:"Информация о лечении"`
	Analyses        []*PetAnalysis `json:"analyses,omitempty" doc:"Список анализов"`
	Bonuses         *PetBonus      `json:"bonuses,omitempty" doc:"Дополнительная информация"`
}

// PetUpdate представляет структуру для обновления существующего питомца
type PetUpdate struct {
	Name            *string        `json:"name,omitempty" validate:"omitempty,min=1,max=100" doc:"Имя питомца" example:"Шарик"`
	ChipNumber      *string        `json:"chipNumber,omitempty" validate:"omitempty,len=15" doc:"Номер чипа (15 символов)" example:"123456789012345"`
	PhotoURL        *string        `json:"photoUrl,omitempty" validate:"omitempty,url,max=255" doc:"URL фотографии питомца" example:"https://example.com/photo.jpg"`
	BreedID         *int           `json:"breedId,omitempty" validate:"omitempty,min=1" doc:"ID породы" example:"1"`
	WeightKg        *float64       `json:"weightKg,omitempty" validate:"omitempty,min=0" doc:"Вес в килограммах" example:"15.5"`
	AgeYears        *int           `json:"ageYears,omitempty" validate:"omitempty,min=0" doc:"Возраст в годах" example:"3"`
	AgeMonths       *int           `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11" doc:"Возраст в месяцах (0-11)" example:"6"`
	BirthDate       *time.Time     `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	LivingCondition *string        `json:"livingCondition,omitempty" doc:"Условия проживания (домашний, выгул на шлейке, самовыгул)" enum:"indoor,leashed,free" example:"indoor"`
	Gender          *string        `json:"gender,omitempty" doc:"Пол питомца (самец, самка)" enum:"male,female" example:"male"`
	Type            *string        `json:"type,omitempty" doc:"Тип животного (собака, кошка)" enum:"dog,cat" example:"dog"`
	BloodGroup      *string        `json:"bloodGroup,omitempty" doc:"Группа крови" example:"DEA 1.1"`
	PetStatus       *string        `json:"petStatus,omitempty" doc:"Статус питомца (донор, реципиент, неизвестно)" enum:"donor,recipient,unknown" example:"donor"`
	Health          *PetHealth     `json:"health,omitempty" doc:"Информация о здоровье"`
	Treatments      *PetTreatment  `json:"treatments,omitempty" doc:"Информация о лечении"`
	Analyses        []*PetAnalysis `json:"analyses,omitempty" doc:"Список анализов"`
	Bonuses         *PetBonus      `json:"bonuses,omitempty" doc:"Дополнительная информация"`
}

// PetResponse представляет ответ с информацией о питомце
type PetResponse struct {
	ID              string         `json:"id" doc:"Уникальный идентификатор питомца" example:"1"`
	Name            string         `json:"name" doc:"Имя питомца" example:"Шарик"`
	ChipNumber      string         `json:"chipNumber,omitempty" doc:"Номер чипа (15 символов)" example:"123456789012345"`
	PhotoURL        string         `json:"photoUrl,omitempty" doc:"URL фотографии питомца" example:"https://example.com/photo.jpg"`
	BreedID         int            `json:"breedId,omitempty" doc:"ID породы" example:"1"`
	WeightKg        float64        `json:"weightKg,omitempty" doc:"Вес в килограммах" example:"15.5"`
	AgeYears        int            `json:"ageYears,omitempty" doc:"Возраст в годах" example:"3"`
	AgeMonths       int            `json:"ageMonths,omitempty" doc:"Возраст в месяцах (0-11)" example:"6"`
	BirthDate       *time.Time     `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	LivingCondition string         `json:"livingCondition,omitempty" doc:"Условия проживания (домашний, выгул на шлейке, самовыгул)" enum:"indoor,leashed,free" example:"indoor"`
	Gender          string         `json:"gender,omitempty" doc:"Пол питомца (самец, самка)" enum:"male,female" example:"male"`
	Type            string         `json:"type" doc:"Тип животного (собака, кошка)" enum:"dog,cat" example:"dog"`
	BloodGroup      string         `json:"bloodGroup,omitempty" doc:"Группа крови" example:"DEA 1.1"`
	PetStatus       string         `json:"petStatus" doc:"Статус питомца (донор, реципиент, неизвестно)" enum:"donor,recipient,unknown" example:"donor"`
	Health          *PetHealth     `json:"health" doc:"Информация о здоровье"`
	Treatments      *PetTreatment  `json:"treatments" doc:"Информация о лечении"`
	Analyses        []*PetAnalysis `json:"analyses" doc:"Список анализов"`
	Bonuses         *PetBonus      `json:"bonuses" doc:"Дополнительная информация"`
	CreatedAt       string         `json:"createdAt" doc:"Дата создания записи" example:"2023-10-01T12:00:00Z"`
	UpdatedAt       string         `json:"updatedAt" doc:"Дата последнего обновления" example:"2023-10-01T12:00:00Z"`
}

// Вспомогательные структуры для Huma

type PetIDPath struct {
	ID string `path:"id" doc:"ID питомца" minLength:"1" example:"PET-25-000001"`
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
	Path string `path:"path" doc:"Путь к аватарке питомца" example:"pets/PET-25-000001/avatar.jpg"`
}

type PetResponseWrapper struct {
	Body PetResponse
}

type PetsResponseWrapper struct {
	Body []PetResponse
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
