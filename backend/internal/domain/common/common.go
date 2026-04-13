package common

// CompensationType represents donor's compensation preference
type CompensationType string

const (
	CompensationFree CompensationType = "free" // Готов помочь безвозмездно
	CompensationPaid CompensationType = "paid" // Не готов помочь бесплатно
	CompensationFood CompensationType = "food" // Готов помочь за корм
)

// PetType представляет тип животного
type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

// Blood groups
var (
	DogBloodGroups = []string{"DEA 1+", "DEA 1-", "unknown"}
	CatBloodGroups = []string{"A", "B", "AB", "unknown"}
)

// GetBloodGroupsByPetType возвращает список групп крови для указанного типа питомца
func GetBloodGroupsByPetType(petType PetType) []string {
	switch petType {
	case PetTypeDog:
		return DogBloodGroups
	case PetTypeCat:
		return CatBloodGroups
	default:
		return []string{}
	}
}
