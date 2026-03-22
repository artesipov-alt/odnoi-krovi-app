package dto

import (
	"time"
)

// ============================================
// Общие типы (вспомогательные)
// ============================================

// PetHealth представляет информацию о здоровье питомца
type PetHealth struct {
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

// PetAnalysis представляет информацию об анализе питомца
type PetAnalysis struct {
	ID           *string    `json:"id,omitempty" doc:"ID анализа в системе" example:"PAN-aBcD1aBcD1" readOnly:"true"`
	AnalysisName *string    `json:"analysisName,omitempty" doc:"Название анализа" enum:"leukemia,immunodeficiency,hemoplasmosis,bartonellosis,babesiosis,dirofilaria,ehrlichiosis,anaplasmosis" example:"leukemia"`
	AnalysisType *string    `json:"analysisType,omitempty" doc:"Тип анализа" enum:"PCR,ELISA,ICA,Microscopy,Express" example:"PCR"`
	AnalysisDate *time.Time `json:"analysisDate,omitempty" doc:"Дата проведения анализа" example:"2023-10-01T12:00:00Z"`
}

// PetAnalysisGroup группирует анализы по типам
type PetAnalysisGroup struct {
	Leukemia         []*PetAnalysis `json:"leukemia,omitempty" doc:"Анализы на лейкемию"`
	Immunodeficiency []*PetAnalysis `json:"immunodeficiency,omitempty" doc:"Анализ на иммунодефицит"`
	Hemoplasmosis    []*PetAnalysis `json:"hemoplasmosis,omitempty" doc:"Анализ на гемоплазмоз"`
	Bartonellosis    []*PetAnalysis `json:"bartonellosis,omitempty" doc:"Анализ на бартонеллез"`
	Babesiosis       []*PetAnalysis `json:"babesiosis,omitempty" doc:"Анализ на бабезиоз"`
	Dirofilaria      []*PetAnalysis `json:"dirofilaria,omitempty" doc:"Анализ на дирофиляриоз"`
	Ehrlichiosis     []*PetAnalysis `json:"ehrlichiosis,omitempty" doc:"Анализ на эрлихиоз"`
	Anaplasmosis     []*PetAnalysis `json:"anaplasmosis,omitempty" doc:"Анализ на анаплазмоз"`
}

// RestrictionFactor представляет фактор ограничения донорства
type RestrictionFactor struct {
	Code           string `json:"code" doc:"Код фактора" example:"STOP_TOO_OLD"`
	Description    string `json:"description" doc:"Описание фактора" example:"Питомцу больше 8 лет"`
	SubDescription string `json:"subDescription,omitempty" doc:"Дополнительное описание фактора"`
}

// DonorRestrictions представляет стоп-факторы и варн-факторы
type DonorRestrictions struct {
	StopFactors []RestrictionFactor `json:"stopFactors,omitempty" doc:"Стоп-факторы"`
	WarnFactors []RestrictionFactor `json:"warnFactors,omitempty" doc:"Предупреждающие факторы"`
}

// PetPathParam представляет параметр пути с ID питомца
type PetPathParam struct {
	ID string `path:"id" doc:"ID питомца" minLength:"1" example:"PET-aBcDeF1234"`
}

// PetUserPathParam представляет параметр пути с ID пользователя
type PetUserPathParam struct {
	UserID string `path:"user_id" doc:"ID пользователя" minLength:"1" example:"USR-aBcDeF1234"`
}

// PetPreloadQuery представляет параметры запроса для подгрузки связанных данных
type PetPreloadQuery struct {
	WithHealth     bool `query:"with_health" doc:"Включить данные о здоровье"`
	WithTreatments bool `query:"with_treatments" doc:"Включить данные о ветеринарных обработках"`
	WithAnalysis   bool `query:"with_analysis" doc:"Включить данные об анализах"`
	WithBonuses    bool `query:"with_bonuses" doc:"Включить данные о бонусах"`
	WithAll        bool `query:"with_all" doc:"Включить все связанные данные"`
}

// ============================================
// CreatePet - Создание питомца
// ============================================

// CreatePetInput представляет запрос на создание питомца
type CreatePetInput struct {
	PetUserPathParam
	Body CreatePetBody
}

// CreatePetBody представляет тело запроса на создание питомца
type CreatePetBody struct {
	Name               string            `json:"name" validate:"required,min=1,max=100" doc:"Имя питомца" example:"Шарик"`
	Type               string            `json:"type" validate:"required,oneof=dog cat" doc:"Тип животного" enum:"dog,cat" example:"dog"`
	WeightKg           float64           `json:"weightKg" validate:"required,gt=0" doc:"Вес в килограммах" example:"15.5"`
	Gender             string            `json:"gender,omitempty" validate:"omitempty,oneof=male female" doc:"Пол питомца" enum:"male,female" example:"male"`
	BirthDate          *time.Time        `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	AgeYears           int               `json:"ageYears,omitempty" validate:"omitempty,min=0,max=30" doc:"Возраст в годах (альтернатива birthDate)" example:"3"`
	AgeMonths          int               `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11" doc:"Возраст в месяцах (дополнение к ageYears)" example:"6"`
	ChipNumber         string            `json:"chipNumber,omitempty" validate:"omitempty,len=15" doc:"Номер чипа" example:"123456789012345"`
	LivingCondition    string            `json:"livingCondition,omitempty" doc:"Условия проживания" enum:"indoor,leash_walking,self_outdoor" example:"indoor"`
	ReproductiveStatus string            `json:"reproductiveStatus,omitempty" doc:"Репродуктивный статус" enum:"pregnancy,lactation,estrus"`
	BreedID            string            `json:"breedId,omitempty" doc:"ID породы" example:"MIX"`
	BloodGroup         string            `json:"bloodGroup,omitempty" doc:"Группа крови" enum:"DEA 1+,DEA 1-,A,B,AB" example:"DEA 1+"`
	Health             *PetHealth        `json:"health,omitempty" doc:"Информация о здоровье"`
	Treatments         *PetTreatment     `json:"treatments,omitempty" doc:"Информация о лечении"`
	Analyses           *PetAnalysisGroup `json:"analyses,omitempty" doc:"Группированные анализы"`
	Bonuses            []string          `json:"bonuses,omitempty" doc:"Дополнительная информация"`
	PetStatus          string            `json:"petStatus,omitempty" doc:"Статус питомца" enum:"donor,recipient,none" deprecated:"true" example:"none"`
}

// CreatePetOutput представляет ответ на создание питомца
type CreatePetOutput struct {
	Body CreatePetResult
}

// CreatePetResult представляет результат создания питомца (только ID и CreatedAt)
type CreatePetResult struct {
	ID        string     `json:"id" doc:"ID созданного питомца" example:"PET-aBcDeF1234"`
	CreatedAt *time.Time `json:"createdAt" doc:"Дата создания" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// UpdatePet - Обновление питомца
// ============================================

// UpdatePetInput представляет запрос на обновление питомца
type UpdatePetInput struct {
	PetPathParam
	Body UpdatePetBody
}

// UpdatePetBody представляет тело запроса на обновление питомца
// Все поля опциональные (указатели) для частичного обновления
type UpdatePetBody struct {
	Name               *string           `json:"name,omitempty" validate:"omitempty,min=1,max=100" doc:"Имя питомца" example:"Шарик"`
	Type               *string           `json:"type,omitempty" validate:"omitempty,oneof=dog cat" doc:"Тип животного" enum:"dog,cat" example:"dog"`
	WeightKg           *float64          `json:"weightKg,omitempty" validate:"omitempty,gt=0" doc:"Вес в килограммах" example:"15.5"`
	Gender             *string           `json:"gender,omitempty" validate:"omitempty,oneof=male female" doc:"Пол питомца" enum:"male,female" example:"male"`
	BirthDate          *time.Time        `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	AgeYears           *int              `json:"ageYears,omitempty" validate:"omitempty,min=0,max=30" doc:"Возраст в годах" example:"3"`
	AgeMonths          *int              `json:"ageMonths,omitempty" validate:"omitempty,min=0,max=11" doc:"Возраст в месяцах" example:"6"`
	ChipNumber         *string           `json:"chipNumber,omitempty" validate:"omitempty,len=15" doc:"Номер чипа" example:"123456789012345"`
	LivingCondition    *string           `json:"livingCondition,omitempty" doc:"Условия проживания" enum:"indoor,leash_walking,self_outdoor" example:"indoor"`
	ReproductiveStatus *string           `json:"reproductiveStatus,omitempty" doc:"Репродуктивный статус" enum:"pregnancy,lactation,estrus"`
	BreedID            *string           `json:"breedId,omitempty" doc:"ID породы" example:"MIX"`
	BloodGroup         *string           `json:"bloodGroup,omitempty" doc:"Группа крови" enum:"DEA 1+,DEA 1-,A,B,AB" example:"DEA 1+"`
	Health             *PetHealth        `json:"health,omitempty" doc:"Информация о здоровье"`
	Treatments         *PetTreatment     `json:"treatments,omitempty" doc:"Информация о лечении"`
	Analyses           *PetAnalysisGroup `json:"analyses,omitempty" doc:"Группированные анализы"`
	Bonuses            *[]string         `json:"bonuses,omitempty" doc:"Дополнительная информация"`
	PetStatus          *string           `json:"petStatus,omitempty" doc:"Статус питомца" enum:"donor,recipient,none" example:"none"`
}

// UpdatePetOutput представляет ответ на обновление питомца
type UpdatePetOutput struct {
	Body UpdatePetResult
}

// UpdatePetResult представляет результат обновления питомца (только ID и UpdatedAt)
type UpdatePetResult struct {
	ID        string     `json:"id" doc:"ID обновленного питомца" example:"PET-aBcDeF1234"`
	UpdatedAt *time.Time `json:"updatedAt" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}

// ============================================
// GetPetByID - Получение питомца по ID
// ============================================

// GetPetByIDInput представляет запрос на получение питомца по ID
type GetPetByIDInput struct {
	PetPathParam
	PetPreloadQuery
}

// GetPetByIDOutput представляет ответ с полными данными питомца
type GetPetByIDOutput struct {
	Body PetDetail
}

type PetDetail struct {
	ID                   string             `json:"id" doc:"Уникальный идентификатор питомца" example:"PET-aBcDeF1234" readOnly:"true"`
	Name                 string             `json:"name" doc:"Имя питомца" example:"Шарик"`
	OwnerName            string             `json:"ownerName,omitempty" doc:"Имя владельца питомца" example:"Иван" readOnly:"true"`
	Type                 string             `json:"type" doc:"Тип животного" enum:"dog,cat" example:"dog"`
	WeightKg             float64            `json:"weightKg" doc:"Вес в килограммах" example:"15.5"`
	Gender               string             `json:"gender" doc:"Пол питомца" enum:"male,female" example:"male"`
	BirthDate            *time.Time         `json:"birthDate,omitempty" doc:"Дата рождения" example:"2020-05-15T00:00:00Z"`
	ChipNumber           string             `json:"chipNumber,omitempty" doc:"Номер чипа" example:"123456789012345"`
	PhotoURLs            []string           `json:"photoUrls,omitempty" doc:"URLs фотографий" example:"https://example.com/photo.jpg"`
	LivingCondition      string             `json:"livingCondition,omitempty" doc:"Условия проживания" enum:"indoor,leash_walking,self_outdoor" example:"indoor"`
	ReproductiveStatus   string             `json:"reproductiveStatus,omitempty" doc:"Репродуктивный статус" enum:"pregnancy,lactation,estrus"`
	BreedID              string             `json:"breedId,omitempty" doc:"ID породы" example:"MIX"`
	BloodGroup           string             `json:"bloodGroup,omitempty" doc:"Группа крови" enum:"DEA 1+,DEA 1-,A,B,AB" example:"DEA 1+"`
	PetStatus            string             `json:"petStatus" doc:"Статус питомца" enum:"none,donor,recipient,blood_found,recovering,planned_donation" example:"donor"`
	AvailableBloodAmount int32              `json:"availableBloodAmount,omitempty" doc:"Доступный объем крови для донации в мл" example:"450"`
	DonorRestrictions    *DonorRestrictions `json:"donorRestrictions,omitempty" doc:"Стоп-факторы и предупреждения"`
	Health               *PetHealth         `json:"health,omitempty" doc:"Информация о здоровье"`
	Treatments           *PetTreatment      `json:"treatments,omitempty" doc:"Информация о лечении"`
	Analyses             *PetAnalysisGroup  `json:"analyses,omitempty" doc:"Группированные анализы"`
	Bonuses              []string           `json:"bonuses,omitempty" doc:"Дополнительная информация"`
	CreatedAt            *time.Time         `json:"createdAt,omitempty" doc:"Дата создания" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	UpdatedAt            *time.Time         `json:"updatedAt,omitempty" doc:"Дата обновления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
	DeletedAt            *time.Time         `json:"deletedAt,omitempty" doc:"Дата удаления" example:"2023-10-01T12:00:00Z" readOnly:"true"`
}

// ============================================
// GetPetsByUser - Получение питомцев пользователя
// ============================================

// GetPetsByUserInput представляет запрос на получение питомцев пользователя
type GetPetsByUserInput struct {
	PetUserPathParam
	PetPreloadQuery
}

// GetPetsByUserOutput представляет ответ со списком питомцев и общей информацией
type GetPetsByUserOutput struct {
	Body GetPetsByUserResult
}

// GetPetsByUserResult представляет результат получения питомцев пользователя
type GetPetsByUserResult struct {
	Pets             []PetDetail `json:"pets" doc:"Список питомцев"`
	PlannedDonations []any       `json:"plannedDonations"`
	TotalPets        int         `json:"totalPets" doc:"Общее количество питомцев у пользователя" example:"5"`
	TotalDonations   int         `json:"totalDonations,omitempty" doc:"Общее количество планируемых донаций питомцев пользователя" example:"12"` // Добавлено по запросу
}

// ============================================
// DeletePet - Удаление питомца
// ============================================

// DeletePetInput представляет запрос на удаление питомца
type DeletePetInput struct {
	PetPathParam
}

// DeletePetOutput представляет ответ на удаление питомца
type DeletePetOutput struct {
	Body DeletePetResult
}

// DeletePetResult представляет результат удаления питомца
type DeletePetResult struct {
	Message string `json:"message" doc:"Сообщение о результате операции" example:"Питомец успешно удален"`
}

// ============================================
// ValidateDonor - Валидация донора
// ============================================

// ValidateDonorInput представляет запрос на валидацию донора
type ValidateDonorInput struct {
	PetPathParam
}

// ValidateDonorOutput представляет ответ на валидацию донора
type ValidateDonorOutput struct {
	Body ValidateDonorResult
}

// ValidateDonorResult представляет результат валидации донора
type ValidateDonorResult struct {
	ID        string     `json:"id" doc:"ID питомца" example:"PET-aBcDeF1234"`
	UpdatedAt *time.Time `json:"updatedAt" doc:"Дата обновления" example:"2023-10-01T12:00:00Z"`
}
