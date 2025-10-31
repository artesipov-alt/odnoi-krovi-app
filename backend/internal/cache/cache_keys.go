package cache

import "time"

// CacheKey представляет шаблоны ключей для кэша
const (
	// User keys
	UserByIDKey       = "user:id:%d"
	UserByTelegramKey = "user:telegram:%d"
	UserProfileKey    = "user:profile:%d"

	// Pet keys
	PetByIDKey    = "pet:id:%d"
	PetsByUserKey = "pets:user:%d"
	PetProfileKey = "pet:profile:%d"

	// Vet clinic keys
	ClinicByIDKey   = "clinic:id:%d"
	ClinicByUserKey = "clinic:user:%d"
	ClinicsListKey  = "clinics:list"

	// Blood stock keys
	BloodStockByIDKey     = "blood_stock:id:%d"
	BloodStockByClinicKey = "blood_stock:clinic:%d"
	BloodStockListKey     = "blood_stock:list"

	// Blood type keys
	BloodTypeByIDKey  = "blood_type:id:%d"
	BloodTypesListKey = "blood_types:list"

	// Blood component keys
	BloodComponentByIDKey = "blood_component:id:%d"

	// Blood group keys
	BloodGroupsByPetTypeKey = "blood_groups:pet_type:%s"

	// Breed keys
	BreedByIDKey    = "breed:id:%d"
	BreedsListKey   = "breeds:list"
	BreedsByTypeKey = "breeds:type:%s"

	// Location keys
	LocationByIDKey  = "location:id:%d"
	LocationsListKey = "locations:list"

	// Rate limiting keys
	RateLimitKey = "rate_limit:%s:%s"

	// Session keys
	SessionKey = "session:%s"
)

// TTL константы для времени жизни кэша
const (
	// Short TTL для часто меняющихся данных
	ShortTTL = 5 * time.Minute

	// Medium TTL для данных, которые меняются реже
	MediumTTL = 30 * time.Minute

	// Long TTL для справочных данных
	LongTTL = 24 * time.Hour

	// VeryLong TTL для редко меняющихся данных
	VeryLongTTL = 7 * 24 * time.Hour
)
