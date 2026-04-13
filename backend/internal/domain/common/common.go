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
// BloodComponent represents a blood component
type BloodComponent struct {
	ID   string
	Name string
}

// BloodComponents contains predefined blood components
var BloodComponents = []BloodComponent{
	{ID: "BLC-1", Name: "Цельная кровь"},
	{ID: "BLC-2", Name: "Эритроцитарная масса"},
	{ID: "BLC-3", Name: "Свежезамороженная плазма"},
	{ID: "BLC-4", Name: "Замороженная плазма"},
	{ID: "BLC-5", Name: "Тромбоконцентрат"},
	{ID: "BLC-6", Name: "Обогащенная тромбоцитами плазма"},
	{ID: "BLC-7", Name: "Криопреципитат"},
	{ID: "BLC-8", Name: "Криосупернатант"},
}

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
