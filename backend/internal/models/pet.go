package models

import (
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
	PetStatus  string  `json:"petStatus" example:"donor"`
	WeightKg   float64 `gorm:"type:numeric" json:"weightKg,omitempty" example:"25.5"`
	BloodGroup string  `json:"bloodGroup,omitempty" example:"DEA 1.1"`
	//==============Дополнительная регистрация для донорства включает в себя=================
	Gender              Gender          `json:"gender,omitempty" example:"male"`
	AgeYears            int             `json:"ageYears,omitempty" example:"3"`
	AgeMonths           int             `json:"ageMonths,omitempty" example:"6"`
	HasChip             bool            `json:"hasChip" example:"false"`
	ChipNumber          string          `gorm:"size:50" json:"chipNumber,omitempty" example:"123456789"`
	PhotoURL            string          `gorm:"size:255" json:"photoUrl,omitempty" example:"https://example.com/photo.jpg"`
	KnowsBloodGroup     bool            `json:"knowsBloodGroup" example:"false"`
	IsGuideDog          bool            `json:"isGuideDog" example:"false"`
	IsTherapist         bool            `json:"isTherapist" example:"false"`
	Breed               string          `gorm:"size:100" json:"breed,omitempty" example:"Лабрадор"`
	Sterilized          bool            `json:"sterilized" example:"false"`
	VaccinationDate     *time.Time      `json:"vaccinationDate,omitempty" example:"2023-01-01T12:00:00Z"`
	DewormingDate       *time.Time      `json:"dewormingDate,omitempty" example:"2023-01-01T12:00:00Z"`
	EctoparasiteDate    *time.Time      `json:"ectoparasiteDate,omitempty" example:"2023-01-01T12:00:00Z"`
	LastTransfusionDate *time.Time      `json:"lastTransfusionDate,omitempty" example:"2023-01-01T12:00:00Z"`
	Latitude            float64         `gorm:"type:numeric" json:"latitude,omitempty" example:"55.7558"`
	Longitude           float64         `gorm:"type:numeric" json:"longitude,omitempty" example:"37.6173"`
	LivingCondition     LivingCondition `json:"livingCondition,omitempty" example:"apartment"`
	DeletedAt           *gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

// В BeforeCreate хуках
func (v *Pet) BeforeCreate(tx *gorm.DB) error {
	var nextVal int
	tx.Raw("SELECT nextval('user_id_seq')").Scan(&nextVal)
	v.ID = PrefixPET.Generate(nextVal)
	return nil
}
