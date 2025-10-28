package models

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
	LivingConditionApartment LivingCondition = "apartment"
	LivingConditionHouse     LivingCondition = "house"
	LivingConditionAviary    LivingCondition = "aviary"
	LivingConditionOther     LivingCondition = "other"
)

// UserRole представляет роль пользователя в системе
type UserRole string

const (
	UserRoleUser   UserRole = "user"
	UserRoleClinic UserRole = "clinic"
	UserRoleAdmin  UserRole = "admin"
	UserRoleDonor  UserRole = "donor"
)

// BloodSearchStatus представляет статус поиска донора
type BloodSearchStatus string

const (
	BloodSearchStatusActive    BloodSearchStatus = "active"
	BloodSearchStatusCompleted BloodSearchStatus = "completed"
	BloodSearchStatusCancelled BloodSearchStatus = "cancelled"
	BloodSearchStatusExpired   BloodSearchStatus = "expired"
)

// BloodStockStatus представляет статус запаса крови
type BloodStockStatus string

const (
	BloodStockStatusActive   BloodStockStatus = "active"
	BloodStockStatusReserved BloodStockStatus = "reserved"
	BloodStockStatusUsed     BloodStockStatus = "used"
	BloodStockStatusExpired  BloodStockStatus = "expired"
)

// DonationStatus представляет статус донорства
type DonationStatus string

const (
	DonationStatusScheduled DonationStatus = "scheduled"
	DonationStatusCompleted DonationStatus = "completed"
	DonationStatusCancelled DonationStatus = "cancelled"
	DonationStatusNoShow    DonationStatus = "no_show"
)

// DonorRequirements представляет требования к донорам
type DonorRequirements struct {
	MinAge           int      `json:"minAge,omitempty"`
	MaxAge           int      `json:"maxAge,omitempty"`
	MinWeight        float64  `json:"minWeight,omitempty"`
	HealthConditions []string `json:"healthConditions,omitempty"`
	Vaccinations     []string `json:"vaccinations,omitempty"`
	BloodTypes       []string `json:"bloodTypes,omitempty"`
}
