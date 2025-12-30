package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Представление структуры питомца в системе
type Pet struct {
	//Сигнатура ID питомца включает в себя префикс питомца, год и шестизначный номер
	ID string `gorm:"primaryKey;" json:"id" example:"PET-25-000001"`
	//==============Простая регистрация для поиска крови включает в себя=====================
	OwnerID    string  `json:"ownerId,omitempty" example:"USR-25-0001"`
	Name       string  `gorm:"size:100;not null" json:"name" example:"Бобик"`
	Type       PetType `json:"type,omitempty" example:"dog"`
	PetStatus  PetRole `json:"petStatus" example:"donor"`
	WeightKg   float64 `gorm:"type:numeric" json:"weightKg,omitempty" example:"25.5"`
	BloodGroup string  `json:"bloodGroup,omitempty" example:"DEA 1+"`
	//==============Дополнительная регистрация для донорства включает в себя=================
	Gender          Gender          `json:"gender,omitempty" example:"male"`
	AgeYears        int             `json:"ageYears,omitempty" example:"3"`
	AgeMonths       int             `json:"ageMonths,omitempty" example:"6"`
	ChipNumber      string          `gorm:"size:15" json:"chipNumber,omitempty" example:"123456789012345"`
	PhotoURL        string          `gorm:"size:255" json:"photoUrl,omitempty" example:"https://example.com/photo.jpg"`
	Breed           string          `gorm:"size:100" json:"breed,omitempty" example:"Лабрадор"`
	LivingCondition LivingCondition `json:"livingCondition,omitempty" example:"indoor"`

	// --- СВЯЗИ ---
	// Has One: У питомца есть одна запись о здоровье.
	// foreignKey:ID говорит, что в таблице PetHealth ключом является ID, который ссылается на Pet.ID
	Health     PetHealth     `gorm:"foreignKey:ID" json:"health"`
	Treatments PetTreatments `gorm:"foreignKey:ID" json:"treatments"`
	Analysis   PetAnalysis   `gorm:"foreignKey:ID" json:"analysis"`

	CreatedAt time.Time       `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time       `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

// Информация о здоровье питомца
type PetHealth struct {
	// ID питомца (соответствует ID в структуре Pet)
	ID string `gorm:"primaryKey;" json:"id" example:"PET-25-000001"`
	// Репродуктивный статус (беременность, лактация и т.д.)
	ReproductiveStatus ReproductiveStatus `json:"reproductiveStatus,omitempty" example:"pregnancy"`
	// Текущее состояние здоровья
	HealthStatus HealthStatus `json:"healthStatus,omitempty" example:"healthy"`
	// Дата последней донации
	LastDonation string `json:"lastDonation,omitempty" example:"2023-01-01T00:00:00Z"`
	// Было ли раннее переливание крови?
	Transfused bool `json:"transfused,omitempty" example:"true"`
	// Принимаемые лекарственные препараты
	Medications string `json:"medications,omitempty" example:"antibiotics"`
	// Хирургические вмешательства перечисление
	SurgicalInterventions string          `json:"surgicalInterventions,omitempty" example:"spaying"`
	CreatedAt             time.Time       `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt             time.Time       `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	DeletedAt             *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

// Информация о ветеринарных обработках питомца
type PetTreatments struct {
	// ID питомца (соответствует ID в структуре Pet)
	ID string `gorm:"primaryKey;" json:"id" example:"PET-25-000001"`
	// Дата последней вакцинации от бешенства
	RabiesVaccinationDate string `json:"rabiesVaccinationDate,omitempty" example:"2023-01-01T00:00:00Z"`
	// Дата последней вакцинации от инфекций
	InfectionVaccinationDate string `json:"infectionVaccinationDate,omitempty" example:"2023-01-01T00:00:00Z"`
	// Дата последней обработки от эктопаразитов (блохи, клещи)
	EctoparasiteTreatmentDate string `json:"ectoparasiteTreatmentDate,omitempty" example:"2023-01-01T00:00:00Z"`
	// Дата последней дегельминтизации (обработка от глистов)
	DewormingDate string          `json:"dewormingDate,omitempty" example:"2023-01-01T00:00:00Z"`
	CreatedAt     time.Time       `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time       `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	DeletedAt     *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

// Информация о последних анализах питомца
type PetAnalysis struct {
	// ID питомца (соответствует ID в структуре Pet)
	ID string `gorm:"primaryKey;" json:"id" example:"PET-25-000001"`

	// Лейкоз (FeLV)
	LeukemiaDate string       `json:"leukemiaDate,omitempty" example:"2023-01-01T00:00:00Z"`
	LeukemiaType AnalysisType `json:"leukemiaType,omitempty" example:"PCR"`

	// Иммунодефицит (FIV)
	ImmunodeficiencyDate string       `json:"immunodeficiencyDate,omitempty" example:"2023-01-01T00:00:00Z"`
	ImmunodeficiencyType AnalysisType `json:"immunodeficiencyType,omitempty" example:"ELISA"`

	// Гемоплазмоз
	HemoplasmosisDate string       `json:"hemoplasmosisDate,omitempty" example:"2023-01-01T00:00:00Z"`
	HemoplasmosisType AnalysisType `json:"hemoplasmosisType,omitempty" example:"PCR"`

	// Бартонеллез
	BartonellosisDate string       `json:"bartonellosisDate,omitempty" example:"2023-01-01T00:00:00Z"`
	BartonellosisType AnalysisType `json:"bartonellosisType,omitempty" example:"PCR"`

	// Бабезиоз
	BabesiosisDate string       `json:"babesiosisDate,omitempty" example:"2023-01-01T00:00:00Z"`
	BabesiosisType AnalysisType `json:"babesiosisType,omitempty" example:"Microscopy"`

	// Дирофиляриоз
	DirofilariaDate string       `json:"dirofilariaDate,omitempty" example:"2023-01-01T00:00:00Z"`
	DirofilariaType AnalysisType `json:"dirofilariaType,omitempty" example:"PCR"`

	// Эрлихиоз
	EhrlichiosisDate string       `json:"ehrlichiosisDate,omitempty" example:"2023-01-01T00:00:00Z"`
	EhrlichiosisType AnalysisType `json:"ehrlichiosisType,omitempty" example:"Express"`

	// Анаплазмоз
	AnaplasmosisDate string       `json:"anaplasmosisDate,omitempty" example:"2023-01-01T00:00:00Z"`
	AnaplasmosisType AnalysisType `json:"anaplasmosisType,omitempty" example:"PCR"`

	CreatedAt time.Time       `json:"createdAt" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time       `json:"updatedAt" example:"2023-01-01T00:00:00Z"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

// Определяем тип конкретно для префиксов питомца
type PetPrefix string

const (
	// Основной префикс для питомцев
	PetIDPrefix PetPrefix = "PET"
)

// Создание префикса для сущности питомца
func (e PetPrefix) Generate(sequenceNum int) string {
	year := time.Now().Year() % 100
	return fmt.Sprintf("%s-%02d-%06d", e, year, sequenceNum)
}

// Проверка валидности префикса сущности питомца
func (e PetPrefix) IsValid() bool {
	return e == PetIDPrefix
}

// PetType представляет тип животного
type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

// Gender представляет пол животного
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// LivingCondition представляет условия проживания животного
type LivingCondition string

const (
	LivingConditionIndoor  LivingCondition = "indoor"
	LivingConditionLeash   LivingCondition = "leash_walking"
	LivingConditionOutdoor LivingCondition = "self_outdoor"
)

// HealthStatus представляет состояние здоровья животного
type HealthStatus string

const (
	HealthStatusHealthy HealthStatus = "healthy"
	HealthStatusIll     HealthStatus = "ill"
	HealthStatusUnknown HealthStatus = "unknown"
)

// ReproductiveStatus представляет физиологическое состояние животного
type ReproductiveStatus string

const (
	ReproductiveStatusPregnancy ReproductiveStatus = "pregnancy"
	ReproductiveStatusLactation ReproductiveStatus = "lactation"
	ReproductiveStatusEstrus    ReproductiveStatus = "estrus"
	ReproductiveStatusNone      ReproductiveStatus = "none"
)

// PetRole представляет роль питомца в системе донорства крови
type PetRole string

const (
	PetRoleDonor     PetRole = "donor"
	PetRoleRecipient PetRole = "recipient"
)

// AnalysisType представляет метод проведения анализа
type AnalysisType string

const (
	AnalysisTypePCR        AnalysisType = "PCR"        // ПЦР
	AnalysisTypeELISA      AnalysisType = "ELISA"      // ИФА
	AnalysisTypeICA        AnalysisType = "ICA"        // ИХА
	AnalysisTypeMicroscopy AnalysisType = "Microscopy" // Микроскопия мазка
	AnalysisTypeExpress    AnalysisType = "Express"    // Экспресс-тест
)

// В BeforeCreate хуках
func (v *Pet) BeforeCreate(tx *gorm.DB) error {
	var nextVal int
	// Получаем следующий номер из последовательности
	if err := tx.Raw("SELECT nextval('pet_id_seq')").Scan(&nextVal).Error; err != nil {
		return err
	}

	// Используем наш новый тип и его метод
	v.ID = PetIDPrefix.Generate(nextVal)
	return nil
}
