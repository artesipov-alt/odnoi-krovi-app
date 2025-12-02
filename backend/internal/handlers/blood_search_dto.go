package handlers

// BloodSearchPetRequest представляет данные питомца для добавления в пул поиска крови
// @Description Данные питомца-реципиента для пула поиска крови
type BloodSearchPetRequest struct {
	// ID питомца в системе
	PetID string `json:"petId" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Тип животного (dog, cat и т.д.)
	PetType string `json:"petType" example:"dog"`
	// Группа крови животного
	BloodGroup string `json:"bloodGroup" example:"DEA 1.1+"`
	// Необходимые компоненты крови
	BloodComponents []string `json:"bloodComponents" example:"whole_blood,plasma"`
	// Необходимый объем крови в мл
	BloodVolumeNeeded int32 `json:"bloodVolumeNeeded" example:"500"`
	// Зарезервированный объем крови в мл
	BloodVolumeReserved int32 `json:"bloodVolumeReserved" example:"100"`
	// ID регионов для поиска доноров
	Regions []int32 `json:"regions" example:"1,2,3"`
	// Разрешить уведомления для маленьких питомцев
	SmallPetsNotifyAllowed bool `json:"smallPetsNotifyAllowed" example:"true"`
	// Статус поиска
	Status string `json:"status" example:"active"`
}

// BloodSearchPetResponse представляет статус добавления питомца в пул
// @Description Статус операции с питомцем в пуле поиска крови
type BloodSearchPetResponse struct {
	// ID питомца в системе
	PetID string `json:"petId" example:"123e4567-e89b-12d3-a456-426614174000"`
	// Статус операции
	Status string `json:"status" example:"added"`
}

// BloodSearchFilterRequest представляет фильтры для поиска питомцев в пуле
// @Description Фильтры для поиска питомцев-реципиентов
type BloodSearchFilterRequest struct {
	// Тип животного для поиска
	PetType string `json:"petType" example:"dog"`
	// Группа крови для поиска
	BloodGroup string `json:"bloodGroup" example:"DEA 1.1+"`
	// ID регионов для фильтрации
	Regions []int32 `json:"regions" example:"1,2,3"`
}

// BloodSearchPetsResponse представляет список найденных питомцев
// @Description Список питомцев из пула поиска крови
type BloodSearchPetsResponse struct {
	// Список питомцев
	Pets []BloodSearchPetRequest `json:"pets"`
}
